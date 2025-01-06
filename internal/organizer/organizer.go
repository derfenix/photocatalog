package organizer

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/derfenix/photocatalog/v2/internal/metadata"
)

const (
	defaultDirMode  = 0o777
	defaultFileMode = 0o644
)

type MetaExtractor interface {
	Extract(_ string, data io.Reader) (metadata.Metadata, error)
}

type Mode interface {
	PlaceIt(sourcePath, targetPath string, mode os.FileMode) error
}

func NewOrganizer(mode Mode, source, target string) *Organizer {
	return &Organizer{
		mode:      mode,
		sourceDir: source,
		targetDir: target,

		dirMode:  defaultDirMode,
		fileMode: defaultFileMode,

		metaExtractors: map[string]MetaExtractor{
			"":     &metadata.Default{},
			"jpg":  metadata.Exif{},
			"jpeg": metadata.Exif{},
			"tiff": metadata.Exif{},
		},
	}
}

type Organizer struct {
	mode Mode

	sourceDir string
	targetDir string

	overwrite bool
	dirMode   os.FileMode
	fileMode  os.FileMode
	errLogger func(error)

	metaExtractors map[string]MetaExtractor
}

func (o *Organizer) WithOverwrite() *Organizer {
	o.overwrite = true

	return o
}

func (o *Organizer) WithDirMode(mode os.FileMode) *Organizer {
	o.dirMode = mode

	return o
}

func (o *Organizer) WithFileMode(mode os.FileMode) *Organizer {
	o.fileMode = mode

	return o
}

func (o *Organizer) WithErrLogger(f func(error)) *Organizer {
	o.errLogger = f

	return o
}

func (o *Organizer) logErr(err error) {
	if o.errLogger != nil {
		o.errLogger(err)
	}
}

func (o *Organizer) Watch(ctx context.Context, wg *sync.WaitGroup) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("new watcher: %w", err)
	}

	if err := watcher.Add(o.sourceDir); err != nil {
		return fmt.Errorf("add source dir to watcher: %w", err)
	}

	// Add all subfolders to the watcher.
	err = filepath.WalkDir(o.sourceDir, func(path string, d fs.DirEntry, err error) error {
		if path == o.sourceDir {
			return nil
		}

		if d.IsDir() {
			if err := watcher.Add(path); err != nil {
				o.logErr(fmt.Errorf("add the directory %s to watcher: %w", path, err))
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("add subdirs to watcher: %w", err)
	}

	wg.Add(2)

	go func() {
		defer wg.Done()

		<-ctx.Done()

		if err := watcher.Close(); err != nil {
			o.logErr(fmt.Errorf("close watcher: %w", err))
		}
	}()

	go func() {
		defer wg.Done()

		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Write {
					stat, err := os.Stat(event.Name)
					if err != nil {
						o.logErr(fmt.Errorf("stat %s: %w", event.Name, err))

						continue
					}

					// Add new directories to the watcher.
					if stat.IsDir() {
						if err := watcher.Add(event.Name); err != nil {
							o.logErr(fmt.Errorf("add the directory %s to watcher: %w", event.Name, err))
						}

						continue
					}

					if err := o.processFile(event.Name); err != nil {
						o.logErr(fmt.Errorf("process file %s: %w", event.Name, err))
					}
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (o *Organizer) FullSync(ctx context.Context) error {
	err := filepath.WalkDir(o.sourceDir, func(path string, info fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if info.IsDir() {
			return nil
		}

		if err := o.processFile(path); err != nil {
			log.Printf("Process file `%s` failed: %s", path, err.Error())

			return nil
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walking source dir: %w", err)
	}

	return nil
}

func (o *Organizer) getMetaForPath(fp string) (metadata.Metadata, error) {
	file, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return metadata.Metadata{}, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	meta, err := o.getMetadata(fp, file)
	if err != nil {
		return metadata.Metadata{}, fmt.Errorf("get metadatas: %w", err)
	}

	return meta, nil
}

func (o *Organizer) getMetadata(fp string, data io.Reader) (metadata.Metadata, error) {
	ext := strings.ToLower(filepath.Ext(fp))

	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}

	extractor, ok := o.metaExtractors[ext]
	if !ok {
		extractor = o.metaExtractors[""]
	}

	meta, err := extractor.Extract(fp, data)
	if err != nil || meta.Created.IsZero() {
		// Fallback to default extractor.
		extractor = o.metaExtractors[""]

		meta, err = extractor.Extract(fp, data)
		if err != nil {
			return metadata.Metadata{}, fmt.Errorf("extract metadata: %w", err)
		}
	}

	return meta, nil
}

func (o *Organizer) processFile(sourcePath string) error {
	meta, err := o.getMetaForPath(sourcePath)
	if err != nil {
		return err
	}

	targetPath, err := o.BuildTargetPath(sourcePath, meta)
	if err != nil {
		return fmt.Errorf("build target path: %w", err)
	}

	if o.pathExists(targetPath) && !o.overwrite {
		return nil
	}

	if err := o.ensureTargetPath(filepath.Dir(targetPath)); err != nil {
		return fmt.Errorf("ensure target path: %w", err)
	}

	if err := o.mode.PlaceIt(sourcePath, targetPath, o.fileMode); err != nil {
		return fmt.Errorf("place file to new path: %w", err)
	}

	log.Printf("File %s placed at %s", sourcePath, targetPath)

	return nil
}

func (o *Organizer) BuildTargetPath(sourcePath string, meta metadata.Metadata) (string, error) {
	sourcePath, err := filepath.Rel(o.sourceDir, sourcePath)
	if err != nil {
		return "", fmt.Errorf("get a relative path: %w", err)
	}

	target := filepath.Join(
		o.targetDir,
		strconv.Itoa(meta.Created.Year()),
		strconv.Itoa(int(meta.Created.Month())),
		strconv.Itoa(meta.Created.Day()),
		sourcePath,
	)

	return target, nil
}

func (o *Organizer) ensureTargetPath(targetPath string) error {
	if o.pathExists(targetPath) {
		return nil
	}

	relPath, err := filepath.Rel(o.targetDir, targetPath)
	if err != nil {
		return fmt.Errorf("get a relative path: %w", err)
	}

	dir := o.targetDir

	for _, part := range strings.Split(relPath, string(filepath.Separator)) {
		dir = filepath.Join(dir, part)

		if err := os.Mkdir(dir, o.dirMode); err != nil && !os.IsExist(err) {
			return fmt.Errorf("create target directory path at %s: %w", dir, err)
		}
	}

	return nil
}

func (o *Organizer) pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		o.logErr(fmt.Errorf("pathExists stat %s: %w", path, err))

		return true
	}

	return true
}

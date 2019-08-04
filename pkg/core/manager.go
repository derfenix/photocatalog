package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"

	"photocatalog/pkg/metadata"
)

type Manager struct {
	TargetPath string
	Mode       ManageMode

	processor       func(fp, targetDir string) (string, error)
	extractorsCache map[string]metadata.Extractor
}

func NewManager(target string, mode ManageMode) (*Manager, error) {
	manager := Manager{
		TargetPath: target,
		Mode:       mode,
		processor:  nil,
	}
	if err := manager.initProcessor(); err != nil {
		return nil, err
	}
	manager.extractorsCache = map[string]metadata.Extractor{}
	return &manager, nil
}

func (m *Manager) buildTarget(meta *metadata.Metadata) (string, error) {
	dir, err := m.dirPathFromMeta(meta)
	if err != nil {
		return "", err
	}
	dirPath := path.Join(m.TargetPath, dir)
	err = os.MkdirAll(dirPath, os.FileMode(0770))
	if err != nil {
		return "", err
	}
	return dirPath, nil
}

func (m *Manager) dirPathFromMeta(meta *metadata.Metadata) (string, error) {
	t := meta.Time
	year := t.Format("2006")
	month := t.Format("01")
	day := t.Format("02")
	return path.Join(year, month, day), nil
}

func (m *Manager) getMetadataExtractor(fp string) metadata.Extractor {
	switch strings.ToLower(path.Ext(fp)) {
	case ".jpeg", ".jpg":
		if _, ok := m.extractorsCache["jpeg"]; !ok {
			m.extractorsCache["jpeg"] = metadata.NewJpegExtractor()
		}
		return m.extractorsCache["jpeg"]
	default:
		if _, ok := m.extractorsCache["default"]; !ok {
			m.extractorsCache["default"] = metadata.NewDefaultExtractor()
		}
		return m.extractorsCache["default"]
	}
}

func (m *Manager) initProcessor() error {
	switch m.Mode {
	case Copy:
		m.processor = func(fp, targetDir string) (string, error) {
			_, fn := path.Split(fp)
			target := path.Join(targetDir, fn)
			cmd := exec.Command("cp", "-f", "--reflink=auto", fp, target)
			return target, cmd.Run()
		}
	case Hardlink:
		m.processor = func(fp, targetDir string) (string, error) {
			_, fn := path.Split(fp)
			target := path.Join(targetDir, fn)
			cmd := exec.Command("ln", "-f", fp, target)
			return target, cmd.Run()
		}
	default:
		return fmt.Errorf("failed to init processor: invalid Mode value")
	}
	return nil
}

func (m *Manager) Manage(fp string) error {
	if m.processor == nil {
		return fmt.Errorf("no processor initialized")
	}
	// Skip hidden files
	if strings.HasPrefix(path.Base(fp), ".") {
		return nil
	}

	log.Println("processing", fp)

	extractor := m.getMetadataExtractor(fp)
	if extractor == nil {
		return fmt.Errorf("failed to get md extractor for %s", fp)
	}

	md, err := extractor.Extract(fp)
	if err != nil {
		return errors.WithMessagef(err, "failed to extract md from %s", fp)
	}

	targetDir, err := m.buildTarget(&md)
	if err != nil {
		return errors.WithMessagef(err, "failed to create dir for %s", fp)
	}

	target, err := m.processor(fp, targetDir)
	if err != nil {
		return errors.WithMessagef(err, "failed to process %s to %s", fp, targetDir)
	}

	if m.Mode == Hardlink {
		log.Println(fp, "linked to", target)
	} else if m.Mode == Copy {
		log.Println(fp, "copied to", target)
	}

	log.Println("success")
	return nil
}

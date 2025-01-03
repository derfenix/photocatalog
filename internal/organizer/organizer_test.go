package organizer

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/derfenix/photocatalog/internal/metadata"
	"github.com/derfenix/photocatalog/internal/organizer/modes"
)

func TestOrganizer_FullSync(t *testing.T) {
	t.Parallel()

	source := "./testdata/fullsync/source"
	target := "./testdata/fullsync/target"

	t.Cleanup(func() {
		if err := os.RemoveAll(target); err != nil {
			t.Fatal(err)
		}
	})

	if err := os.Mkdir(target, 0777); err != nil {
		t.Fatalf("create target dir %s failed: %v", target, err)
	}

	ctx := context.Background()

	org := NewOrganizer(modes.HardLink{}, source, target)
	if err := org.FullSync(ctx); err != nil {
		t.Fatalf("full sync failed: %v", err)
	}

	err := filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		sourceFile, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read source file %s: %v", path, err)
		}

		meta, err := (&metadata.Default{}).Extract(path, nil)
		if err != nil {
			return fmt.Errorf("extract metadata from %s: %v", path, err)
		}

		targetPath, err := org.BuildTargetPath(path, meta)
		if err != nil {
			return fmt.Errorf("build target file %s: %v", path, err)
		}

		targetFile, err := os.ReadFile(targetPath)
		if err != nil {
			return fmt.Errorf("read target file %s: %v", targetPath, err)
		}

		if !bytes.Equal(sourceFile, targetFile) {
			return fmt.Errorf("target file content missmatch")
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk dir failed: %v", err)
	}
}

func TestOrganizer_Watch(t *testing.T) {
	t.Parallel()

	source := "./testdata/watcher/source"
	target := "./testdata/watcher/target"

	t.Cleanup(func() {
		if err := os.RemoveAll(target); err != nil {
			t.Fatal(err)
		}
	})

	if err := os.Mkdir(target, 0777); err != nil {
		t.Fatalf("create target dir %s failed: %v", target, err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(cancel)

	wg := sync.WaitGroup{}

	org := NewOrganizer(&modes.HardLink{}, source, target)
	if err := org.Watch(ctx, &wg); err != nil {
		t.Fatalf("watch failed: %v", err)
	}

	nonEmpty, err := checkEmpty(t, target)
	if err != nil {
		t.Fatalf("check empty failed: %v", err)
	}

	if nonEmpty {
		t.Fatal("target dir should not be empty")
	}

	err = os.WriteFile(filepath.Join(source, "20241108_160834.txt"), []byte("test"), 0777)
	if err != nil {
		t.Fatalf("file write failed: %v", err)
	}

	time.Sleep(time.Millisecond)

	nonEmpty, err = checkEmpty(t, target)
	if err != nil {
		t.Fatalf("check empty failed: %v", err)
	}

	if !nonEmpty {
		t.Fatal("target dir should not be empty")
	}

	cancel()

	wg.Wait()
}

func checkEmpty(t *testing.T, target string) (bool, error) {
	t.Helper()

	var nonEmpty bool
	err := filepath.WalkDir(target, func(path string, d fs.DirEntry, err error) error {
		if path == target {
			return nil
		}

		nonEmpty = true

		return nil
	})
	if err != nil {
		t.Fatalf("walk dir failed: %v", err)
	}

	return nonEmpty, err
}

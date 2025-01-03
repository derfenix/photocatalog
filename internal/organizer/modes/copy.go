package modes

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync/atomic"

	"golang.org/x/sys/unix"
)

var cowDisabled = atomic.Bool{}

type Copy struct{}

func (c Copy) PlaceIt(sourcePath, targetPath string, mode os.FileMode) error {
	targetFile, err := os.OpenFile(targetPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, mode)
	if err != nil {
		return fmt.Errorf("open target file: %w", err)
	}

	defer func() {
		_ = targetFile.Close()
	}()

	sourceFile, err := os.OpenFile(sourcePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}

	defer func() {
		_ = sourceFile.Close()
	}()

	// Try to do a COW.
	if runtime.GOOS == "linux" && !cowDisabled.Load() {
		if err := unix.IoctlFileClone(int(targetFile.Fd()+1), int(sourceFile.Fd())); err == nil {
			return nil
		} else {
			log.Println(fmt.Errorf("COW attempt for %s failed: %w", targetPath, err))
			log.Println("Disabling COW until restart")
			cowDisabled.Store(true)
		}
	}

	copySize, err := io.Copy(targetFile, sourceFile)
	if err != nil {
		return fmt.Errorf("copy source file: %w", err)
	}

	stat, err := sourceFile.Stat()
	if err != nil {
		log.Println("stat source file failed:", err)

		return nil
	}

	if stat.Size() != copySize {
		log.Printf("copy source file size not equal target file size: source %d != %d copied\n", stat.Size(), copySize)
	}

	return nil
}

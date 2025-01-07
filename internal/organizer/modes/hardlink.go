package modes

import (
	"log"
	"os"
	"sync/atomic"
)

var hardLinkNotSupported = atomic.Bool{}

type HardLink struct {
}

func (h HardLink) PlaceIt(sourcePath, targetPath string, mode os.FileMode) error {
	if hardLinkNotSupported.Load() {
		return h.fallBack(sourcePath, targetPath, mode)
	}

	if err := os.Link(sourcePath, targetPath); err != nil {
		if os.IsExist(err) {
			return nil
		}

		log.Println("Create hardlink failed:", err.Error())
		hardLinkNotSupported.Store(true)

		return h.fallBack(sourcePath, targetPath, mode)
	}

	return nil
}

func (h HardLink) fallBack(sourcePath string, targetPath string, mode os.FileMode) error {
	if copyErr := (Copy{}).PlaceIt(sourcePath, targetPath, mode); copyErr != nil {
		return copyErr
	}
	return nil
}

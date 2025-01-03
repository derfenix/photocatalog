package modes

import (
	"fmt"
	"os"
)

type Move struct {
}

func (m Move) PlaceIt(sourcePath, targetPath string, mode os.FileMode) error {
	if err := os.Rename(sourcePath, targetPath); err != nil {
		return fmt.Errorf("rename %s to %s: %w", sourcePath, targetPath, err)
	}

	if err := os.Chmod(targetPath, mode); err != nil {
		return fmt.Errorf("chmod hard link: %w", err)
	}

	return nil
}

package modes

import (
	"fmt"
	"os"
)

type SymLink struct {
}

func (s SymLink) PlaceIt(sourcePath, targetPath string, _ os.FileMode) error {
	if err := os.Symlink(sourcePath, targetPath); err != nil {
		if os.IsExist(err) {
			return nil
		}

		return fmt.Errorf("create symlink: %w", err)
	}

	return nil
}

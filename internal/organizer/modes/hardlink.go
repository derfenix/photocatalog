package modes

import (
	"fmt"
	"os"
)

type HardLink struct {
}

func (h HardLink) PlaceIt(sourcePath, targetPath string, mode os.FileMode) error {
	if err := os.Link(sourcePath, targetPath); err != nil {
		if os.IsExist(err) {
			return nil
		}

		return fmt.Errorf("create hard link: %w", err)
	}

	return nil
}

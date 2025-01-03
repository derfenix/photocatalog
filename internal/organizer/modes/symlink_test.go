package modes_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/derfenix/photocatalog/internal/organizer/modes"
)

func TestSymLink_PlaceIt(t *testing.T) {
	t.Parallel()

	const testDataDir = "symlink"

	t.Cleanup(func() {
		if err := os.RemoveAll(fmt.Sprintf("./testdata/%s/target/", testDataDir)); err != nil {
			t.Errorf("error removing ./testdata/%s/target/: %v", testDataDir, err)
		}
	})

	if err := os.Mkdir(fmt.Sprintf("./testdata/%s/target/", testDataDir), 0777); err != nil {
		t.Errorf("error creating ./testdata/%s/target/: %v", testDataDir, err)
	}

	source := fmt.Sprintf("./testdata/%s/source/foo.txt", testDataDir)
	target := fmt.Sprintf("./testdata/%s/target/foo.txt", testDataDir)

	err := SymLink{}.PlaceIt(source, target, 0644)
	if err != nil {
		t.Errorf("place file: %v", err)
	}

	linkedFilePath, err := os.Readlink(target)
	if err != nil {
		t.Errorf("read target file: %v", err)
	}

	if linkedFilePath != source {
		t.Errorf("linked file path is %s, want %s", linkedFilePath, source)
	}
}

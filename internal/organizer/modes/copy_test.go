package modes_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/derfenix/photocatalog/v2/internal/organizer/modes"
)

func TestCopy_PlaceIt(t *testing.T) {
	t.Parallel()

	const testDataDir = "copy"

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

	err := Copy{}.PlaceIt(source, target, 0644)
	if err != nil {
		t.Errorf("place file: %v", err)
	}

	targetData, err := os.ReadFile(target)
	if err != nil {
		t.Errorf("read target file: %v", err)
	}

	sourceData, err := os.ReadFile(source)
	if err != nil {
		t.Errorf("read source file: %v", err)
	}

	if string(targetData) != string(sourceData) {
		t.Error("copy file contents missmatch")
	}
}

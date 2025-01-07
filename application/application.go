package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/derfenix/photocatalog/v2/internal/organizer"
	"github.com/derfenix/photocatalog/v2/internal/organizer/modes"
)

type Application struct {
	config Config
}

func NewApplication(config Config) (Application, error) {
	if err := config.Validate(); err != nil {
		return Application{}, fmt.Errorf("invalid config: %w", err)
	}

	return Application{config: config}, nil
}

func (a *Application) Start(ctx context.Context, wg *sync.WaitGroup) error {
	var mode organizer.Mode

	switch a.config.Mode {
	case ModeCopy:
		mode = modes.Copy{}
	case ModeHardlink:
		mode = modes.HardLink{}
	case ModeMove:
		mode = modes.Move{}
	case ModeSymlink:
		mode = modes.SymLink{}
	default:
		mode = modes.HardLink{}
	}

	org := organizer.NewOrganizer(mode, a.config.SourceDir, a.config.TargetDir).
		WithDirMode(os.FileMode(a.config.DirMode)).
		WithFileMode(os.FileMode(a.config.FileMode)).
		WithErrLogger(func(err error) {
			log.Println("ERROR:", err.Error())
		})

	if a.config.Overwrite {
		org = org.WithOverwrite()
	}

	if a.config.Watch {
		if err := org.Watch(ctx, wg); err != nil {
			return fmt.Errorf("initialize watch: %w", err)
		}
	}

	if !a.config.SkipFullSync {
		if err := org.FullSync(ctx); err != nil {
			return fmt.Errorf("full sync: %w", err)
		}
	}

	return nil
}

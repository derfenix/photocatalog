package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"slices"
	"strconv"
	"sync"
	"syscall"

	"github.com/derfenix/photocatalog/v2/application"
)

func main() {
	cfg := loadCfg()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app, err := application.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}

	if err := app.Start(ctx, wg); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}

func loadCfg() application.Config {
	cfg := application.Config{
		DirMode:  0774,
		FileMode: 0644,
	}

	flag.StringVar(&cfg.SourceDir, "source", "", "Source directory")
	flag.StringVar(&cfg.TargetDir, "target", "", "Target directory")
	flag.BoolVar(&cfg.Overwrite, "overwrite", false, "Overwrite existing files")
	flag.BoolVar(&cfg.Watch, "watch", false, "Watch for changes in the source directory")
	flag.BoolVar(&cfg.Watch, "monitor", false, "Watch for changes in the source directory") // Legacy option
	flag.BoolVar(&cfg.SkipFullSync, "skip-full-sync", false, "Skip full sync at startup")

	flag.Func("dir-mode", "Mode bits for directories can be created while syncing", func(s string) error {
		var err error

		cfg.DirMode, err = strconv.ParseUint(s, 8, 32)
		if err != nil {
			return err
		}

		return nil
	})

	flag.Func("file-mode", "Mode bits for files created while syncing (not applicable for hardlink mode)", func(s string) error {
		var err error

		cfg.FileMode, err = strconv.ParseUint(s, 8, 32)
		if err != nil {
			return err
		}

		return nil
	})

	flag.Func("mode", "Organizing mode", func(s string) error {
		if s == "" {
			cfg.Mode = application.ModeHardlink
		}

		cfg.Mode = application.Mode(s)

		if !slices.Contains(application.SupportedModes, cfg.Mode) {
			return fmt.Errorf("invalid mode, supported modes: %s", application.SupportedModes)
		}

		return nil
	})

	flag.Parse()

	// Legacy fallback
	if cfg.SourceDir == "" {
		log.Println("Source directory not specified. May be using old systemd unit file.")

		cfg.SourceDir = flag.Arg(0)
	}

	return cfg
}

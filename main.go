package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/derfenix/photocatalog/application"
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
		SourceDir: "",
		TargetDir: "",
		Mode:      "",
		Overwrite: false,
		DirMode:   0,
		FileMode:  0,
		Watch:     false,
	}

	flag.StringVar(&cfg.SourceDir, "source", "", "Source directory")
	flag.StringVar(&cfg.TargetDir, "target", "", "Target directory")
	flag.BoolVar(&cfg.Overwrite, "overwrite", false, "Overwrite existing files")
	flag.BoolVar(&cfg.Watch, "watch", false, "Watch for changes in the source directory")

	var dirMode string
	var fileMode string
	flag.StringVar(&dirMode, "dirmode", "0777", "Mode bits for directories can be created while syncing")
	flag.StringVar(&fileMode, "filemode", "0644", "Mode bits for files created while syncing (not applicable for hardlink mode)")

	var mode string
	flag.StringVar(&mode, "mode", "hardlink", "Mode")

	flag.Parse()

	cfg.Mode = application.Mode(mode)

	var err error

	cfg.DirMode, err = strconv.ParseUint(dirMode, 8, 32)
	if err != nil {
		cfg.DirMode = 0o777
	}

	cfg.FileMode, err = strconv.ParseUint(fileMode, 8, 32)
	if err != nil {
		cfg.DirMode = 0o644
	}

	return cfg
}

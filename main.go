package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
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
	cfg := application.Config{}

	flag.StringVar(&cfg.SourceDir, "source", "", "Source directory")
	flag.StringVar(&cfg.TargetDir, "target", "", "Target directory")
	flag.BoolVar(&cfg.Overwrite, "overwrite", false, "Overwrite existing files")
	flag.BoolVar(&cfg.Watch, "watch", true, "Watch for changes in the source directory")
	flag.BoolVar(&cfg.SkipFullSync, "skip-full-sync", false, "Skip full sync at startup")

	var (
		dirMode  string
		fileMode string
		mode     string
	)

	flag.StringVar(&dirMode, "dir-mode", "0777", "Mode bits for directories can be created while syncing")
	flag.StringVar(&fileMode, "file-mode", "0644", "Mode bits for files created while syncing (not applicable for hardlink mode)")
	flag.StringVar(&mode, "mode", "hardlink", "Organizing mode")

	flag.Parse()

	cfg.Mode = application.Mode(mode)

	var err error

	cfg.DirMode, err = strconv.ParseUint(dirMode, 8, 32)
	if err != nil {
		log.Println("Parse -dir-mode failed:", err)

		cfg.DirMode = 0o777
	}

	cfg.FileMode, err = strconv.ParseUint(fileMode, 8, 32)
	if err != nil {
		log.Println("Parse -file-mode failed:", err)

		cfg.DirMode = 0o644
	}

	return cfg
}

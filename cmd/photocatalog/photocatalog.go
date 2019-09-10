package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"

	"github.com/derfenix/photocatalog/pkg/manager"
)

func main() {
	mode := flag.String("mode", "hardlink", "Manage mode: copy or hardlink")
	target := flag.String("target", "./", "Root directory to organize files in")
	monitor := flag.String("monitor", "", "Monitor specified folder for new files")
	updateMtime := flag.Bool("update_mtime", false, "Update mtime on organized files to be equal to the files creation date")

	flag.Parse()
	args := flag.Args()
	log.Println("Using", *target, "as target and", *mode, "as mode")

	var manageMode manager.ManageMode
	switch *mode {
	case "copy":
		manageMode = manager.Copy
	case "hardlink":
		manageMode = manager.Hardlink
	default:
		log.Fatalf("Invalid mode %s", *mode)
	}

	mgr, err := manager.NewManager(*target, manageMode, *updateMtime)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if *monitor == "" {
		processFiles(args, mgr)
	} else {
		startMonitoring(*monitor, mgr)
	}
}

func processFiles(args []string, manager *manager.Manager) {
	var manageErr error
	var gotErrors bool

	if len(args) > 0 {
		var err error
		if len(args) == 1 && strings.HasSuffix(args[0], "/") {
			args, err = filepath.Glob(args[0] + "*")
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Println("Processing", len(args), "files")
		for _, f := range args {
			manageErr = manager.Manage(f)
			if manageErr != nil {
				log.Println(manageErr)
				gotErrors = true
			}
		}
	} else {
		log.Println("No input files")
	}

	if gotErrors {
		log.Println("All files processed, got errors")
	} else {
		log.Println("All files processed without errors")
	}
}

func startMonitoring(monitor string, manager *manager.Manager) {
	var manageErr error

	if !path.IsAbs(monitor) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get CWD: %s", err.Error())
		}
		monitor = path.Join(cwd, monitor)
	}

	log.Println("Monitoring", monitor)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		closeErr := watcher.Close()
		if closeErr != nil {
			log.Println(closeErr)
		}
	}()

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Create {
					if strings.HasSuffix(event.Name, "tmp") {
						continue
					}
					manageErr = manager.Manage(event.Name)
					if manageErr != nil {
						log.Println(manageErr)
					}
				}
			case wErr, ok := <-watcher.Errors:
				log.Println("error:", wErr)
				if !ok {
					return
				}
			}
		}
	}()

	err = watcher.Add(monitor)
	if err != nil {
		log.Fatal(err)
	}

	sig := <-done
	log.Println("Monitoring stopped with", sig, "signal")
}

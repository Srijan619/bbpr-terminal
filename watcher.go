package main

import (
	"log"
	"os"
	"path/filepath"
	"simple-git-terminal/state"
	widgets "simple-git-terminal/widgets/table"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rivo/tview"
)

const debounceInterval = 500 * time.Millisecond

var (
	fileEventTimestamps = make(map[string]time.Time)
	mu                  sync.Mutex // protect fileEventTimestamps
)

func watchFiles(app *tview.Application) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watchAllDirs(watcher, ".") // watch current project root
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			file := event.Name

			// Ignore .log and temp/backup files
			if filepath.Ext(file) == ".log" ||
				strings.HasSuffix(file, "~") ||
				strings.HasSuffix(file, ".swp") {
				continue
			}

			// Only care about writes and creates
			if !(event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) {
				continue
			}

			// Throttle duplicate triggers
			mu.Lock()
			last, seen := fileEventTimestamps[file]
			now := time.Now()
			if seen && now.Sub(last) < debounceInterval {
				mu.Unlock()
				continue
			}
			fileEventTimestamps[file] = now
			mu.Unlock()

			log.Println("Triggering reload for:", event)

			app.QueueUpdateDraw(func() {
				for _, view := range state.PipelineUIState.Views {
					if r, ok := view.(widgets.Refreshable); ok {
						r.Refresh()
					}
				}
			})
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("watcher error:", err)
		}
	}
}

func watchAllDirs(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Optionally skip hidden dirs like .git, .air, etc.
			if info.Name() == ".git" || info.Name() == "tmp" {
				return filepath.SkipDir
			}
			log.Println("Watching:", path)
			return watcher.Add(path)
		}

		// Skip .log files
		if filepath.Ext(info.Name()) == ".log" {
			return nil
		}
		return nil
	})
}

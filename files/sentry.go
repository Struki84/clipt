package files

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/struki84/clipt/internal/tools/library"
)

type FileSentry struct {
	chromaClient library.ChromaClient
	dirPath      string
}

func NewFileSentry(dirPath string) *FileSentry {
	return &FileSentry{
		chromaClient: *library.NewChromaClient(),
		dirPath:      dirPath,
	}

}

func (sentry *FileSentry) ScanFiles() error {
	err := filepath.Walk(sentry.dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		if !info.IsDir() {
			log.Println(path)
		}

		return nil

	})

	if err != nil {
		log.Println("Error scanning files:", err)
		return err
	}

	return nil
}

func (sentry *FileSentry) WatchFiles(ctx context.Context) error {
	log.Println("Watching files...")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Error creating watcher:", err)
		return err
	}

	defer watcher.Close()

	err = watcher.Add(sentry.dirPath)
	if err != nil {
		log.Println("Error adding directory to watcher:", err)
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						log.Println("Error getting file info:", err)
						continue
					}

					if !fileInfo.IsDir() {
						log.Println("File changed:", event.Name)

						path := fmt.Sprintf("%s/%s", sentry.dirPath, event.Name)
						err := sentry.chromaClient.SaveFile(ctx, path, fileInfo)
						if err != nil {
							log.Println("Error saving file to chroma DB:", err)
						}
					}

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("Watcher error:", err)
			}
		}
	}()

	return nil
}

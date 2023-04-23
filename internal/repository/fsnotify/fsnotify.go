package fsnotify

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type CopyFile func(src, dest string) error
type Validator func(path string) error

type Repository struct {
	validator    Validator
	copyFile     CopyFile
	backupFolder string
	syncFiles    map[string]bool
}

type Params struct {
	Validator Validator
	CopyFile  CopyFile
}

func New(params Params) *Repository {
	return &Repository{
		copyFile:  params.CopyFile,
		validator: params.Validator,
		syncFiles: make(map[string]bool),
	}
}

func (r *Repository) Sync(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("Error creating new watcher: %s", err)
	}
	defer func() { _ = watcher.Close() }()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("File changed: %s\n", event.Name)

					backupPath := filepath.Join(r.backupFolder, event.Name)
					err = r.copyFile(event.Name, backupPath)
					if err != nil {
						log.Println("Error copying file to backup folder:", err)
						continue
					}
					fmt.Printf("Copied changed file to backup folder: %s\n", backupPath)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	for filePath, _ := range r.syncFiles {
		err = watcher.Add(filePath)
		if err != nil {
			return fmt.Errorf("Error adding watched folder: %s", err)
		}
	}

	<-done

	return nil
}

func (r *Repository) AddSyncFile(path string) error {
	if err := r.validator(path); err != nil {
		return fmt.Errorf("error validating path: %w", err)
	}

	r.syncFiles[path] = true

	return nil
}

func (r *Repository) SetBackupFolder(path string) error {
	if err := r.validator(path); err != nil {
		return fmt.Errorf("error validating path: %w", err)
	}

	r.backupFolder = path

	return nil
}

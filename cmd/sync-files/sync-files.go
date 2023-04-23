package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/mark-marushak/sync-files/internal/filewather"
	"github.com/mark-marushak/sync-files/internal/repository/fsnotify"
	"github.com/mark-marushak/sync-files/internal/utils"
)

var syncFiles = []string{
	"/home/sandbox/testbackup",
}

func main() {
	repository := fsnotify.New(fsnotify.Params{
		Validator: utils.ValidatePath,
		CopyFile:  utils.CopyFile,
	})

	if err := repository.SetBackupFolder("/home/sandbox/backup"); err != nil {
		panic(err)
	}

	watcher := filewather.New(repository)

	err := watcher.AddFiles(syncFiles...)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err = watcher.Sync(ctx); err != nil {
		return
	}
}

package filewather

import (
	"context"
	"fmt"
)

type Repository interface {
	Sync(ctx context.Context) error
	AddSyncFile(path string) error
	SetBackupFolder(path string) error
}

type Watcher struct {
	repository Repository
}

func New(repository Repository) *Watcher {
	return &Watcher{
		repository: repository,
	}
}

func (w *Watcher) Sync(ctx context.Context) error {
	return w.repository.Sync(ctx)
}

func (w *Watcher) AddFiles(path ...string) error {
	for i := 0; i < len(path); i++ {
		if err := w.repository.AddSyncFile(path[i]); err != nil {
			return fmt.Errorf("error adding file to sync: %w", err)
		}
	}

	return nil
}

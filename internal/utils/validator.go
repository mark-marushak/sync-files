package utils

import (
	"fmt"
	"os"
)

func ValidatePath(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func() { _ = f.Close() }()

	return f.Sync()
}

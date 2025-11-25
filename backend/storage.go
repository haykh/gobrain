package backend

import (
	"os"
	"path/filepath"
)

// Storage abstracts filesystem interactions so that future backends (git, cloud
// sync, etc.) can plug in without rewriting business logic.
type Storage interface {
	Ensure(paths Paths) error
	ListMarkdown(path string) ([]string, error)
	Create(filepath string) (*os.File, error)
	Move(src, dst string) error
	Exists(path string) bool
}

// localStorage is the default Storage implementation backed by the local
// filesystem.
type localStorage struct{}

func (localStorage) Ensure(paths Paths) error {
	for _, path := range paths.All() {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}

func (localStorage) ListMarkdown(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
			filenames = append(filenames, file.Name())
		}
	}
	return filenames, nil
}

func (localStorage) Create(filepath string) (*os.File, error) { //nolint:revive
	return os.Create(filepath)
}

func (localStorage) Move(src, dst string) error {
	return os.Rename(src, dst)
}

func (localStorage) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

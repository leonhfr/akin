package files

import (
	"os"
	"path/filepath"

	"github.com/leonhfr/akin/cfg"
)

func Files(config *cfg.Config, files chan<- string, errors chan<- error) {
	walk(config.Path, config.Extensions, files, errors)
}

func walk(root string, extensions map[string]bool, files chan<- string, errors chan<- error) {
	defer close(files)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errors <- err
			return err
		}
		if isFile(path) && isExtension(path, extensions) {
			files <- path
		}
		return nil
	})
}

func isFile(path string) bool {
	fi, _ := os.Stat(path)
	mode := fi.Mode()
	return mode.IsRegular()
}

func isExtension(path string, extensions map[string]bool) bool {
	ext := filepath.Ext(path)
	_, ok := extensions[ext]
	return ok
}

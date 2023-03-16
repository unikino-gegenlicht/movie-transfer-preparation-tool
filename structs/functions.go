package structs

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// NewFileFromPath takes a file path and creates a new File object from it
func NewFileFromPath(path string) (*File, error) {
	// cleanup the given filepath
	path = filepath.Clean(path)
	// check if the given path exists
	fileStats, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create new file object from path")
	}
	// check if the path is a directory
	if fileStats.IsDir() {
		return nil, fmt.Errorf("cannot create file object for directory: %s", path)
	}
	// now create the new file object
	f := new(File)
	f.Path = path
	f.Size = fileStats.Size()
	f.Extension = filepath.Ext(path)
	return f, nil
}

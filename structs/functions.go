package structs

import (
	"fmt"
	"github.com/pkg/errors"
	consts "movie-transfer-preparation-tool/const"
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

func ConvertMovieToMovieConfigurationEntry(m Movie) *MovieConfigurationEntry {
	// create a new movie configuration entry
	e := new(MovieConfigurationEntry)
	// assign the attributes from the movie to the configuration entry
	e.Title = m.Title
	e.ScreeningDateTime = e.ScreeningDateTime
	e.MovieLanguage = m.AudioLanguage.Code
	e.SubtitleLanguage = m.SubtitleLanguage.Code
	// now build the directory path for the configuration entry
	movieDateTime := e.ScreeningDateTime.Format(consts.DateTimeFolderFormat)
	dirPath := fmt.Sprintf("./%s – %s", movieDateTime, e.Title)
	e.Path = dirPath
	return e
}

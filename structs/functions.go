package structs

import (
	"fmt"
	consts "movie-transfer-preparation-tool/const"
)

func ConvertMovieToMovieConfigurationEntry(m Movie) *MovieConfigurationEntry {
	// create a new movie configuration entry
	e := new(MovieConfigurationEntry)
	// assign the attributes from the movie to the configuration entry
	e.Title = m.Title
	e.ScreeningDateTime = m.ScreeningDate
	if m.AudioLanguage != nil {
		e.MovieLanguage = m.AudioLanguage.Code
	} else {
		e.MovieLanguage = ""
	}
	if m.SubtitleLanguage != nil {
		e.SubtitleLanguage = m.SubtitleLanguage.Code
	} else {
		e.SubtitleLanguage = ""
	}
	// now build the directory path for the configuration entry
	movieDateTime := e.ScreeningDateTime.Format(consts.DateTimeFolderFormat)
	dirPath := fmt.Sprintf("./%s – %s", movieDateTime, e.Title)
	e.Path = dirPath
	return e
}

package structs

import "time"
import "github.com/emvi/iso-639-1"

type Movie struct {
	// Title is the title of the movie
	Title string

	// ScreeningDate contains the date at which the movie will be started. This
	// value should also the time to allow the correct sorting of the movie
	ScreeningDate time.Time

	// AudioLanguage contains the main language of the movie as the ISO 639-1
	// two character code
	AudioLanguage *iso6391.Language

	// SubtitleLanguage contains the main language of the subtitles as the ISO
	// 639-1 two character code.
	SubtitleLanguage *iso6391.Language

	// VideoFile contains the file information for the movie file
	VideoFile File

	// SubtitleFile contains the file information for the subtitle file
	SubtitleFile File
}

type File struct {
	// Path takes the full path of the file
	Path string
	// Extension contains the file extension inferred from the Path
	Extension string
	// Size contains the file size in Bytes
	Size int64
}

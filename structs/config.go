package structs

import "time"

// Configuration contains the configuration written later on to a toml file
type Configuration struct {
	Semester       string                    `toml:"semester"`
	FirstScreening time.Time                 `toml:"firstScreeningDate"`
	LastScreening  time.Time                 `toml:"lastScreeningDate"`
	Movies         []MovieConfigurationEntry `toml:"movies"`
}

type MovieConfigurationEntry struct {
	// Path contains the path to the folder containing the movie file and the
	// subtitle file relative to the drive root.
	Path string `toml:"path"`

	// Title contains the title of the movie
	Title string `toml:"title"`

	// ScreeningDateTime contains the date and time at which the movie will be
	// screened
	ScreeningDateTime time.Time `toml:"screeningDateTime"`

	// MovieLanguage contains the two letter iso639-1 code of the spoken
	// language in the movie
	MovieLanguage string `toml:"movieLanguage"`

	// SubtitleLanguage contains the two letter iso639-1 code of the spoken
	// language in the movie
	SubtitleLanguage string `toml:"subtitleLanguage"`
}

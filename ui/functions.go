package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	iso6391 "github.com/emvi/iso-639-1"
	"movie-transfer-preparation-tool/bindings"
	consts "movie-transfer-preparation-tool/const"
	"movie-transfer-preparation-tool/structs"
	"movie-transfer-preparation-tool/validators"
	"time"
)

var movieDataForm *widget.Form
var newMovie *structs.Movie

func AddMovieOnClick() {
	// create a new form for the needed data and create a object storing the movie information
	movieDataForm = new(widget.Form)
	newMovie = new(structs.Movie)

	// create entries for the needed data
	movieNameEntry := widget.NewEntryWithData(binding.BindString(&newMovie.Title))
	movieNameEntry.SetPlaceHolder("Der Herr der Ringe")
	movieNameEntry.Validator = validators.NoEmptyOrWhitespaces
	movieDataForm.Append("Titel", movieNameEntry)

	movieScreeningDateEntry := widget.NewEntry()
	movieScreeningDateEntry.SetPlaceHolder("01.01.2001 20:00")
	movieScreeningDateEntry.Validator = validators.DateTimeSpaced
	movieDataForm.Append("Vorführungsdatum", movieScreeningDateEntry)

	// get the names of the languages and prepend the option NONE to the selection
	languageSelectionOptions := append([]string{"NONE"}, iso6391.Names...)

	// now create the selection fields for the languages
	audioLanguageSelect := widget.NewSelect(languageSelectionOptions, func(selection string) {
		// resolve the selected language to the code
		if selection == "NONE" {
			newMovie.AudioLanguage = nil
		} else {
			language := iso6391.FromName(selection)
			newMovie.AudioLanguage = &language
		}
	})
	movieDataForm.Append("Sprache der Tonspur", audioLanguageSelect)

	subtitleLanguageSelect := widget.NewSelect(languageSelectionOptions, func(selection string) {
		// resolve the selected language to the code
		if selection == "NONE" {
			newMovie.SubtitleLanguage = nil
		} else {
			language := iso6391.FromName(selection)
			newMovie.SubtitleLanguage = &language
		}
	})
	movieDataForm.Append("Sprache der Untertitel", subtitleLanguageSelect)

	movieDataPopup := dialog.NewCustomConfirm(
		"Metadaten erfassen",
		"Hinzufügen",
		"Abbrechen",
		movieDataForm,
		func(confirmed bool) {
			if !confirmed {
				return
			}
			// parse the date to a time.Time object
			newMovie.ScreeningDate, _ = time.Parse(consts.DateTimeFormat, movieScreeningDateEntry.Text)
			bindings.Movies.Append(newMovie)
		},
		MainWindow,
	)
	movieDataPopup.Resize(fyne.NewSize(600, 350))
	movieDataPopup.Show()
}

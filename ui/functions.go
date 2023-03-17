package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	iso6391 "github.com/emvi/iso-639-1"
	"github.com/rs/zerolog/log"
	_ "github.com/rs/zerolog/log"
	"movie-transfer-preparation-tool/bindings"
	consts "movie-transfer-preparation-tool/const"
	"movie-transfer-preparation-tool/structs"
	"movie-transfer-preparation-tool/validators"
	"os"
	"strings"
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

	// create a dialog for selecting the file for the movie
	movieFileSelectionButton := new(widget.Button)
	movieFileSelectionButton.SetText("Auswählen")
	movieFileSelectionButton.Importance = widget.HighImportance
	movieFileSelectionButton.OnTapped = func() {
		// show a dialog to select the file for the movie
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if err != nil {
				log.Error().Err(err).Msg("unable to select file for movie")
				return
			}
			if closer != nil {
				defer closer.Close()
				newMovie.VideoFile.Path = closer.URI().Path()
				newMovie.VideoFile.Extension = closer.URI().Extension()
				// now calculate the size of the file using os.Stat
				fileStat, _ := os.Stat(newMovie.VideoFile.Path)
				newMovie.VideoFile.Size = fileStat.Size()
				movieFileSelectionButton.SetText(newMovie.VideoFile.Path)
				movieFileSelectionButton.Disable()
			} else {
				log.Info().Msg("no file was selected")
			}
		}, MainWindow)
	}
	movieDataForm.Append("Datei für Film", movieFileSelectionButton)
	// get the names of the languages and prepend the option NONE to the selection
	languageSelectionOptions := append([]string{"NONE"}, iso6391.Names...)

	// now create the selection fields for the languages
	audioLanguageSelect := xwidget.NewCompletionEntry(languageSelectionOptions)
	audioLanguageSelect.OnChanged = func(s string) {
		if len(s) < 2 {
			audioLanguageSelect.HideCompletion()
			return
		}

		var possibleLanguages []string

		for _, lang := range languageSelectionOptions {
			if strings.Contains(strings.ToLower(lang), strings.ToLower(s)) {
				possibleLanguages = append(possibleLanguages, lang)
			}
		}

		if len(possibleLanguages) == 0 {
			audioLanguageSelect.HideCompletion()
			return
		}

		audioLanguageSelect.SetOptions(possibleLanguages)
		audioLanguageSelect.ShowCompletion()
	}
	movieDataForm.Append("Sprache der Tonspur", audioLanguageSelect)

	subtitleFileSelectionButton := new(widget.Button)
	subtitleFileSelectionButton.SetText("Auswählen")
	subtitleFileSelectionButton.Importance = widget.HighImportance
	subtitleFileSelectionButton.OnTapped = func() {
		// show a dialog to select the file for the movie
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if err != nil {
				log.Error().Err(err).Msg("unable to select file for subtitle")
				return
			}
			if closer != nil {
				defer closer.Close()
				newMovie.SubtitleFile.Path = closer.URI().Path()
				newMovie.SubtitleFile.Extension = closer.URI().Extension()
				// now calculate the size of the file using os.Stat
				fileStat, _ := os.Stat(newMovie.VideoFile.Path)
				newMovie.SubtitleFile.Size = fileStat.Size()
				subtitleFileSelectionButton.SetText(newMovie.SubtitleFile.Path)
				subtitleFileSelectionButton.Disable()
			} else {
				log.Info().Msg("no file was selected")
			}
		}, MainWindow)
	}
	movieDataForm.Append("Datei für Untertitel", subtitleFileSelectionButton)
	subtitleLanguageSelect := xwidget.NewCompletionEntry(languageSelectionOptions)
	subtitleLanguageSelect.OnChanged = func(s string) {
		if len(s) < 2 {
			subtitleLanguageSelect.HideCompletion()
			return
		}

		var possibleLanguages []string

		for _, lang := range languageSelectionOptions {
			if strings.Contains(strings.ToLower(lang), strings.ToLower(s)) {
				possibleLanguages = append(possibleLanguages, lang)
			}
		}

		if len(possibleLanguages) == 0 {
			subtitleLanguageSelect.HideCompletion()
			return
		}

		subtitleLanguageSelect.SetOptions(possibleLanguages)
		subtitleLanguageSelect.ShowCompletion()
	}
	movieDataForm.Append("Sprache der Untertitel", subtitleLanguageSelect)

	movieDataPopup := dialog.NewForm(
		"Metadaten erfassen",
		"Hinzufügen",
		"Abbrechen",
		movieDataForm.Items,
		func(confirmed bool) {
			if !confirmed {
				return
			}
			// get the name of the language and put it into the movie
			audioLang := audioLanguageSelect.Text
			if audioLang == "NONE" {
				newMovie.AudioLanguage = nil
			} else {
				al := iso6391.FromName(audioLang)
				newMovie.AudioLanguage = &al
			}
			// get the name of the language and put it into the movie
			subLang := subtitleLanguageSelect.Text
			if subLang == "NONE" {
				newMovie.SubtitleLanguage = nil
			} else {
				sl := iso6391.FromName(subLang)
				newMovie.SubtitleLanguage = &sl
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

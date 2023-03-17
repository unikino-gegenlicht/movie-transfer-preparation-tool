package utilities

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"movie-transfer-preparation-tool/bindings"
	"movie-transfer-preparation-tool/structs"
	"movie-transfer-preparation-tool/types"
	"movie-transfer-preparation-tool/ui"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func WriteFilesToDisk() {
	// create a new context which allows stopping the copy progress
	copyContext, stopCopying := context.WithCancel(context.Background())
	copyProgress := binding.NewFloat()
	// now create a new progress dialog
	progress := ui.NewProgressDialog(
		"Kopiervorgang läuft",
		"Die erstellten Filme werden auf das ausgewälte Laufwerk kopiert",
		copyProgress,
		ui.MainWindow,
	)
	progress.SetOnClosed(func() {
		// if the copy progress was not interrupted by another error
		if copyContext.Err() == nil {
			// open a confirmation popup to check if the user really wants to cancel
			dialog.ShowConfirm("Kopiervorgang abbrechen", "Möchten Sie den Kopiervorgang wirklich abbrechen?", func(cancel bool) {
				if cancel {
					stopCopying()

				}
			}, ui.MainWindow)
		}
	})

	// now get the configuration
	untyped, _ := bindings.Configuration.Get()
	config := untyped.(*structs.Configuration)

	log.Info().Interface("config", config).Send()

	// now calculate the needed space on the storage location
	var neededSpace uint64 = 0
	movies, _ := bindings.Movies.Get()
	for _, m := range movies {
		movie := m.(*structs.Movie)
		neededSpace += uint64(movie.VideoFile.Size)
		neededSpace += uint64(movie.SubtitleFile.Size)
	}
	log.Info().Uint64("neededSpace", neededSpace).Msg("calculated file sizes")

	// now get the available space on the selected drive
	drive, _ := bindings.SelectedDrive.Get()
	if strings.TrimSpace(drive) == "" {
		stopCopying()
		log.Error().Msg("no drive selected for output")
		progress.Hide()
		err := errors.New("no drive selected for output")
		errorDialog := ui.NewErrorDialog("Fehler beim Kopieren", "Es wurde kein Ausgangslaufwerk ausgewälht", err)
		errorDialog.Show()
		return
	}
	availableSpace := GetAvailableSpace(drive)
	log.Info().Uint64("availableSpace", availableSpace).Str("drive", drive).Msg("got free space")

	// now check if the needed space is enough
	if availableSpace <= neededSpace {
		stopCopying()
		log.Error().Msg("insufficient space on target device")
		progress.Hide()
		err := fmt.Errorf("selected storage device needs %d bytes more free space", neededSpace-availableSpace)
		errorDialog := ui.NewErrorDialog("Fehler beim Kopieren", "Das Zielgerät hat nicht genügend freien Speicherplatz", err)
		errorDialog.Show()
		return
	}

	log.Info().Interface("config", config).Msg("added movies to configuration")
	// and write the configuration file
	path := fmt.Sprintf("%s/configuration.toml", drive)
	configurationFile, _ := os.Create(path)
	err := toml.NewEncoder(configurationFile).Encode(config)
	if err != nil {
		stopCopying()
		log.Error().Err(err).Msg("unable to write configuration")
		progress.Hide()
		errorDialog := ui.NewErrorDialog("Fehler beim Kopieren", "Die Konfigurationsdatei konnte nicht geschrieben werden", err)
		errorDialog.Show()
		return
	}

	var totalCopiedBytes uint64
	ticker := time.NewTicker(250 * time.Millisecond)

	go func() {
		for {
			select {
			case <-copyContext.Done():
				return
			case <-ticker.C:
				p := float64(totalCopiedBytes) / float64(neededSpace)
				copyProgress.Set(p)
			}
		}
	}()

	progress.Show()

	// now start iterating over the created movies
	for _, m := range movies {
		// create a waitgroup
		wg := sync.WaitGroup{}
		// typecast the movie into the struct
		movie := m.(*structs.Movie)
		// now convert the movie into a configuration entry
		configEntry := structs.ConvertMovieToMovieConfigurationEntry(*movie)
		// since we now have the configEntry we now can write the files to the target directory
		originMovieFile, err := os.Open(movie.VideoFile.Path)
		if err != nil {
			stopCopying()
			log.Error().Err(err).Msg("unable to open origin movie file")
			progress.Hide()
			errorDialog := ui.NewErrorDialog("Fehler beim Kopieren", "Die angegebene Videodatei kann nicht geöffnet werden", err)
			errorDialog.Show()
			break
		}

		// now use the config entry and the target path to create the directories needed on the target
		// drive
		targetPath := fmt.Sprintf("%s/%s", drive, configEntry.Path)
		targetPath = filepath.Clean(targetPath)
		err = os.MkdirAll(targetPath, 0777)
		if err != nil {
			stopCopying()
			log.Error().Err(err).Msg("unable to open create target directory")
			progress.Hide()
			errorDialog := ui.NewErrorDialog("Fehler beim Kopieren", "Das Zielverzeichnis konnte nicht erstellt werden", err)
			errorDialog.Show()
			break
		}

		// since the target path now exists we can create the target file
		movieTargetFilePath := fmt.Sprintf("%s/%s%s", targetPath, movie.Title, movie.VideoFile.Extension)
		movieTargetFile, err := os.Create(movieTargetFilePath)
		if err != nil {
			stopCopying()
			log.Error().Err(err).Msg("unable to open create target file for movie")
			progress.Hide()
			errorDialog := ui.NewErrorDialog("Fehler beim Kopieren", "Die Zieldatei für den Film konnte nicht erstellt werden",
				err)
			errorDialog.Show()
			break
		}

		// now create the custom writer
		movieFileWriter := &types.CustomWriter{
			Writer:  movieTargetFile,
			Context: copyContext,
		}

		// now start polling to total count of written bytes from the movieFileWriter
		go func() {
			var previousCopiedBytes uint64
			for {
				select {
				case <-copyContext.Done():
					return
				case <-ticker.C:
					movieFileWriter.Mu.Lock()
					delta := uint64(movieFileWriter.BytesWritten) - previousCopiedBytes
					totalCopiedBytes += delta
					previousCopiedBytes = uint64(movieFileWriter.BytesWritten)
					movieFileWriter.Mu.Unlock()
				}
			}
		}()
		go func() {
			wg.Add(1)
			_, err = io.Copy(movieFileWriter, originMovieFile)
			if err != nil {
				if err == context.Canceled {
					log.Info().Msg("user cancelled copy progress")
				}
			}
			wg.Done()
		}()
		if movie.SubtitleLanguage != nil {
			// since we now have the configEntry we now can write the files to the target directory
			originSubtitleFile, err := os.Open(movie.SubtitleFile.Path)
			if err != nil {
				stopCopying()
				log.Error().Err(err).Msg("unable to open origin movie file")
				progress.Hide()
				errorDialog := ui.NewErrorDialog("Fehler beim Kopieren", "Die angegebene Videodatei kann nicht geöffnet werden", err)
				errorDialog.Show()
				break
			}
			// since the target path now exists we can create the target file
			subtitleTargetFilePath := fmt.Sprintf("%s/%s%s", targetPath, movie.SubtitleLanguage.Code, movie.SubtitleFile.Extension)
			subtitleTargetFile, err := os.Create(subtitleTargetFilePath)
			if err != nil {
				stopCopying()
				log.Error().Err(err).Msg("unable to open create target file for movie")
				progress.Hide()
				errorDialog := ui.NewErrorDialog("Fehler beim Kopieren", "Die Zieldatei für den Film konnte nicht erstellt werden",
					err)
				errorDialog.Show()
				break
			}

			// now create the custom writer
			subtitleFileWriter := &types.CustomWriter{
				Writer:  subtitleTargetFile,
				Context: copyContext,
			}
			go func() {
				var previousCopiedBytes uint64
				for {
					select {
					case <-copyContext.Done():
						return
					case <-ticker.C:
						subtitleFileWriter.Mu.Lock()
						delta := uint64(subtitleFileWriter.BytesWritten) - previousCopiedBytes
						totalCopiedBytes += delta
						previousCopiedBytes = uint64(subtitleFileWriter.BytesWritten)
						subtitleFileWriter.Mu.Unlock()
					}
				}
			}()
			go func() {
				wg.Add(1)
				_, err = io.Copy(subtitleFileWriter, originSubtitleFile)
				if err != nil {
					if err == context.Canceled {
						log.Info().Msg("user cancelled copy progress")
					}
				}
				wg.Done()
			}()
			log.Info().Msg("waiting until this movie is copied")
			go func() {
				wg.Wait()
				// now close the open files
				subtitleTargetFile.Close()
				originSubtitleFile.Close()
			}()
		}
		go func() {
			wg.Wait()
			// now close the open files
			movieTargetFile.Close()
			originMovieFile.Close()
		}()
	}

}

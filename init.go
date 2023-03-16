package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"movie-transfer-preparation-tool/ui"
	"os"
	"strconv"
)

func init() {
	// set the global logging level to debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Stack().Logger()
	// now overwrite the output to log to a file and to the console
	logFile, _ := os.OpenFile("movie-transfer-preparation-tool.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 666)

	log.Logger = log.Output(zerolog.MultiLevelWriter(
		logFile,
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "02.01.2006 15:04:05"}))
	log.Info().Msg("Starting Movie Transfer Preparation Tool")
}

func init() {
	// create a splash screen
	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		ui.SplashScreen = drv.CreateSplashWindow()
	} else {
		log.Fatal().Msg("unable to create splash screen. unsupported client")
	}

	// now set the main window to be of a fixed size and resize it
	ui.MainWindow.SetFixedSize(true)
	ui.MainWindow.Resize(fyne.NewSize(800, 600))

}

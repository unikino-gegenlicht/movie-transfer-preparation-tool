package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/image/colornames"
	"image/color"
	"movie-transfer-preparation-tool/bindings"
	"movie-transfer-preparation-tool/resources"
	"movie-transfer-preparation-tool/ui"
	"movie-transfer-preparation-tool/validators"
	"movie-transfer-preparation-tool/vars"
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
	vars.Application.Settings().SetTheme(&ui.CustomTheme{})
}

func init() {
	// create a splash screen
	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		ui.SplashScreen = drv.CreateSplashWindow()
		ui.SplashScreen.SetFixedSize(true)
		ui.SplashScreen.CenterOnScreen()
	} else {
		log.Fatal().Msg("unable to create splash screen. unsupported client")
	}

	// now set the main window to be of a fixed size and resize it
	// ui.MainWindow = vars.Application.NewWindow("Movie Transfer Preperation Tool")
	ui.MainWindow.SetFixedSize(true)
	ui.MainWindow.Resize(fyne.NewSize(800, 600))
}

// setup the splash screen and show it
func init() {
	log.Info().Msg("setting up the splash window")
	splashBackground := canvas.NewImageFromResource(resources.SplashLogo)
	splashBackground.FillMode = canvas.ImageFillContain
	splashBackground.SetMinSize(fyne.NewSize(600, 300))
	// now infer the current color theme
	var textColor color.Color
	if vars.Application.Settings().ThemeVariant() == theme.VariantDark {
		textColor = color.NRGBA{
			R: 255,
			G: 221,
			B: 0,
			A: 255,
		}
	} else {
		textColor = colornames.Black
	}
	currentStepLabel := canvas.NewText("", textColor)
	currentStepLabel.TextSize = 16
	currentStepLabel.Alignment = fyne.TextAlignCenter
	currentStepLabel.TextStyle.Monospace = true
	currentStepLabel.TextStyle.Bold = true
	bindings.CurrentStartupStep.AddListener(binding.NewDataListener(func() {
		currentStepLabel.Text, _ = bindings.CurrentStartupStep.Get()
		currentStepLabel.Refresh()
	}))
	splashLayout := container.New(layout.NewVBoxLayout(),
		splashBackground, currentStepLabel, layout.NewSpacer())
	ui.SplashScreen.SetContent(splashLayout)
	ui.SplashScreen.Show()
}

func init() {
	// now create the input fields for the main semester data
	semesterTitleEntry := widget.NewEntry()
	semesterTitleEntry.Validator = validators.NoEmptyOrWhitespaces

	semesterStartDay := widget.NewEntry()
	semesterStartDay.Validator = validators.Date
	semesterStartDay.TextStyle.Monospace = true

	semesterEndDay := widget.NewEntry()
	semesterEndDay.Validator = validators.Date
	semesterEndDay.TextStyle.Monospace = true

	// now add those inputs to the form
	ui.SemesterDataForm.Append("Semester", semesterTitleEntry)
	ui.SemesterDataForm.Append("Start der Vorführungen", semesterStartDay)
	ui.SemesterDataForm.Append("Ende der Vorführungen", semesterEndDay)

	mainContent := container.New(layout.NewMaxLayout(), ui.SemesterDataForm)
	ui.MainWindow.SetContent(mainContent)
	ui.MainWindow.SetMaster()
	ui.MainWindow.SetIcon(resources.AppIcon)
}

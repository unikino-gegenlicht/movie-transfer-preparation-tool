package main

import (
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"movie-transfer-preparation-tool/bindings"
	"movie-transfer-preparation-tool/ui"
	"movie-transfer-preparation-tool/utilities"
	"movie-transfer-preparation-tool/vars"
	"strings"
	"time"
)

func main() {
	// now create a channel in which we set up the main window and load some
	// other variables
	go func() {
		bindings.CurrentStartupStep.Set("Getting connected/mounted drives")
		// use the utilities to get the mounted drives
		bindings.ExternalDrives.AddListener(binding.NewDataListener(func() {
			d, _ := bindings.ExternalDrives.Get()
			var options []string
			if d == nil {
				return
			}
			externalDrives := d.([][]string)
			for _, externalDrive := range externalDrives {
				option := fmt.Sprintf("%s (%s)", externalDrive[0], externalDrive[1])
				options = append(options, option)
			}
			ui.FileDestinationSelector.Options = options
		}))
		ui.FileDestinationSelector.OnChanged = func(selection string) {
			fields := strings.Fields(selection)
			bindings.SelectedDrive.Set(fields[0])
		}
		bindings.ExternalDrives.Set(utilities.GetExternalDrives())
		// todo: remove sleeps and mock steps
		time.Sleep(1 * time.Second)
		bindings.CurrentStartupStep.Set("reading data")
		time.Sleep(1 * time.Second)
		bindings.CurrentStartupStep.Set("done")
		time.Sleep(1 * time.Second)

		ui.MainWindow.Show()
		ui.SplashScreen.Close()
	}()
	vars.Application.Run()
}

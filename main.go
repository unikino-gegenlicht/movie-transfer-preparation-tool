package main

import (
	"movie-transfer-preparation-tool/bindings"
	"movie-transfer-preparation-tool/ui"
	"movie-transfer-preparation-tool/vars"
	"time"
)

func main() {
	// now create a channel in which we set up the main window and load some
	// other variables
	go func() {
		bindings.CurrentStartupStep.Set("mocking init steps")
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

package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewProgressDialog(title, message string, progress binding.Float, parent fyne.Window) dialog.Dialog {
	// create a progress var
	progressBar := widget.NewProgressBarWithData(progress)
	progressBarLayout := container.New(layout.NewMaxLayout(), progressBar)
	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord
	messageScroll := container.NewVScroll(messageLabel)
	messageScroll.Resize(fyne.NewSize(100, 300))
	i := widget.NewIcon(theme.DefaultTheme().Icon(theme.IconNameWarning))
	i.Resize(fyne.NewSize(24, 24))
	scrollContainer := container.New(layout.NewMaxLayout(), messageScroll)
	dialogContent := container.New(layout.NewBorderLayout(i, progressBarLayout, nil, nil), i, scrollContainer, progressBarLayout)
	d := dialog.NewCustom(title, "Abbrechen", dialogContent, parent)
	d.Resize(fyne.NewSize(600, 400))
	return d
}

func NewErrorDialog(title, message string, e error) dialog.Dialog {
	m := fmt.Sprintf("%s\n\n%v", message, e)
	l := widget.NewLabel(m)
	l.Wrapping = fyne.TextWrapWord
	scroll := container.NewVScroll(l)
	scroll.Resize(fyne.NewSize(100, 200))
	i := widget.NewIcon(theme.DefaultTheme().Icon(theme.IconNameError))
	i.Resize(fyne.NewSize(24, 24))
	scrollContainer := container.New(layout.NewMaxLayout(), scroll)
	c := container.New(layout.NewGridLayout(1), i, scrollContainer)
	d := dialog.NewCustom(title, "OK", c, MainWindow)
	d.Resize(fyne.NewSize(350, 350))
	return d
}

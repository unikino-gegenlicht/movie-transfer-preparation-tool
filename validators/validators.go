package validators

import (
	"fyne.io/fyne/v2/data/validation"
	consts "movie-transfer-preparation-tool/const"
)

var NoEmptyOrWhitespaces = validation.NewRegexp(`[^\s]+`, "Eingabe erforderlich")
var Date = validation.NewTime(consts.DateFormat)
var DateTimeSpaced = validation.NewTime(consts.DateTimeFormat)

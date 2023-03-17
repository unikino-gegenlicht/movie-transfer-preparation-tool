package bindings

import (
	"fyne.io/fyne/v2/data/binding"
	"movie-transfer-preparation-tool/structs"
)

var CurrentStartupStep = binding.NewString()
var ExternalDrives = binding.NewUntyped()
var SelectedDrive = binding.NewString()
var Movies = binding.NewUntypedList()
var Configuration = binding.BindUntyped(structs.Configuration{})

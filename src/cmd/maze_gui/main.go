package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"maze/internal/ui"
)

const (
	windowTitle = "Maze"
)

func main() {
	mazeApp := app.NewWithID(windowTitle)
	window := mazeApp.NewWindow(windowTitle)
	window.Resize(fyne.NewSize(820, 640))

	mazeView := ui.NewMazeView()
	mazeView.Resize(fyne.NewSize(ui.FieldSizePx, ui.FieldSizePx))

	mazeField := container.NewWithoutLayout(mazeView)
	mazeField.Resize(fyne.NewSize(ui.FieldSizePx, ui.FieldSizePx))

	buttons := ui.CreateButtons(window, mazeView)

	content := container.NewHBox(
		buttons,
		container.NewVBox(
			widget.NewLabelWithStyle("Maze", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			mazeField,
		),
	)

	window.SetContent(content)
	window.ShowAndRun()
}

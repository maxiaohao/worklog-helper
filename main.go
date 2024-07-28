package main

import (
	"worklog-helper/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

const APP_TITLE string = "Worklog Helper"

func main() {
	w := app.New().NewWindow(APP_TITLE)

	w.Resize(fyne.NewSize(600, 400))
	w.SetFixedSize(true)

	tabs := container.NewAppTabs(
		container.NewTabItem("Track Time", ui.NewTrackTimeContainer(w)),
		container.NewTabItem("Worklog", ui.NewWorklogTable(w)),
		container.NewTabItem("Settings", ui.NewSettingsPanel(w)),
	)

	w.SetContent(tabs)
	w.ShowAndRun()
}

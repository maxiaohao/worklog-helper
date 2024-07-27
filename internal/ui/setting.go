package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

func NewSettingsForm(w fyne.Window) fyne.CanvasObject {
	email := widget.NewEntry()
	email.SetPlaceHolder("foo@bar.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	atlassianUrlBase := widget.NewEntry()
	atlassianUrlBase.SetPlaceHolder("https://<your-domain>.atlassian.net")
	atlassianUrlBase.Validator = validation.NewRegexp(`^https?:\/\/[^\s/$.?#].[^\s]*$`, "not a valid http(s) url")

	atlassianApiToken := widget.NewPasswordEntry()
	atlassianApiToken.SetPlaceHolder("Atlassian API token")
	atlassianApiToken.Validator = validation.NewRegexp(`^.+$`, "url baase cannot be empty")

	tempoUrlBase := widget.NewEntry()
	tempoUrlBase.SetText("https://api.tempo.io")
	tempoUrlBase.SetPlaceHolder("https://api.tempo.io")
	tempoUrlBase.Validator = validation.NewRegexp(`^https?:\/\/[^\s/$.?#].[^\s]*$`, "not a valid http(s) url")

	tempoApiToken := widget.NewPasswordEntry()
	tempoApiToken.SetPlaceHolder("Tempo API token")
	tempoApiToken.Validator = validation.NewRegexp(`^.+$`, "api token cannot be empty")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Email", Widget: email, HintText: "Your Atlassian Account Email"},
			{Text: "Atlassian Url base", Widget: atlassianUrlBase},
			{Text: "Atlassian API token", Widget: atlassianApiToken, HintText: "Atlassian -> Manage Account -> Security -> Create and manage API tokens"},
			{Text: "Tempo Url base", Widget: tempoUrlBase},
			{Text: "Tempo API token", Widget: tempoApiToken, HintText: "Jira -> Apps -> Tempo -> Settings (left bar) -> API Integration -> New Token"},
		},
		SubmitText: "Save",
		OnSubmit: func() {
			fmt.Println("Saved")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Worklog Helper",
				Content: "Settings saved successfully",
			})
		},
	}
	return form
}

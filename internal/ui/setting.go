package ui

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"worklog-helper/internal/data"
	"worklog-helper/pkg/rest"
	"worklog-helper/pkg/yamlconfig"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	CONFIG_SUB_DIR_NAME = "worklog-helper"
)

func NewSettingsPanel(w fyne.Window) fyne.CanvasObject {
	email := widget.NewEntry()
	email.SetPlaceHolder("foo@bar.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	atlassianUrlBase := widget.NewEntry()
	atlassianUrlBase.SetText("https://<your-organisation>.atlassian.net")
	atlassianUrlBase.SetPlaceHolder("https://<your-organisation>.atlassian.net")
	atlassianUrlBase.Validator = validation.NewRegexp(`^https?:\/\/[^\s/$.?#].[^\s]*$`, "not a valid http(s) url")

	atlassianApiToken := widget.NewEntry()
	atlassianApiToken.SetPlaceHolder("Atlassian API token")

	tempoUrlBase := widget.NewEntry()
	tempoUrlBase.SetText("https://api.tempo.io")
	tempoUrlBase.SetPlaceHolder("https://api.tempo.io")
	tempoUrlBase.Validator = validation.NewRegexp(`^https?:\/\/[^\s/$.?#].[^\s]*$`, "not a valid http(s) url")

	tempoApiToken := widget.NewEntry()
	tempoApiToken.SetPlaceHolder("Tempo API token")

	accountIdLabel := widget.NewLabel("Atlassian Account Id: <To be verified>")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Email", Widget: email, HintText: "Your Atlassian Account Email"},
			{Text: "Atlassian Url base", Widget: atlassianUrlBase, HintText: "Replace <your-organisation>"},
			{Text: "Atlassian API token", Widget: atlassianApiToken, HintText: "Atlassian -> Manage Account -> Security -> Create and manage API tokens"},
			{Text: "Tempo Url base", Widget: tempoUrlBase, HintText: "Default is https://api.tempo.io"},
			{Text: "Tempo API token", Widget: tempoApiToken, HintText: "Jira -> Apps -> Tempo -> Settings (left bar) -> API Integration -> New Token"},
		},
		SubmitText: "Verify & Save",
		OnSubmit: func() {
			atlassianAuthorizationValue := toBasicAuthorization(email.Text, atlassianApiToken.Text)
			accountId, err := getAccountIdByEmail(email.Text, atlassianUrlBase.Text, atlassianAuthorizationValue)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			fmt.Println("==== Atlassian accountId is ", accountId)
			accountIdLabel.SetText(fmt.Sprintf("Atlassian Account Id: %v", accountId))

			settings := data.Settings{
				Email:     email.Text,
				AccountId: accountId,
				AtlassianSettings: data.ApiServerSettings{
					UrlBase:            atlassianUrlBase.Text,
					ApiToken:           atlassianApiToken.Text,
					AuthorizationValue: atlassianAuthorizationValue,
				},
				TempoSettings: data.ApiServerSettings{
					UrlBase:            tempoUrlBase.Text,
					ApiToken:           tempoApiToken.Text,
					AuthorizationValue: fmt.Sprintf("Bearer %v", tempoApiToken.Text),
				},
			}
			yamlconfig.Write(CONFIG_SUB_DIR_NAME, &settings)

			successfulMessage := fmt.Sprintf("The settings are saved successfully to config file at %v", yamlconfig.GetConfigFilePath(CONFIG_SUB_DIR_NAME))
			fmt.Println(successfulMessage)

			dialog.ShowInformation("Saved config file", successfulMessage, w)

			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Worklog Helper",
				Content: successfulMessage,
			})
		},
	}

	// load settings if there is an existing config file
	existingSettings := &data.Settings{}
	err := yamlconfig.Read(existingSettings, CONFIG_SUB_DIR_NAME)
	if err == nil && len(existingSettings.AccountId) > 0 {
		email.SetText(existingSettings.Email)
		atlassianUrlBase.SetText(existingSettings.AtlassianSettings.UrlBase)
		atlassianApiToken.SetText(existingSettings.AtlassianSettings.ApiToken)
		tempoUrlBase.SetText(existingSettings.TempoSettings.UrlBase)
		tempoApiToken.SetText(existingSettings.TempoSettings.ApiToken)
		accountIdLabel.SetText(fmt.Sprintf("Atlassian Account Id: %v", existingSettings.AccountId))

		fmt.Printf("Settings have been loaded from config file at %v\n", yamlconfig.GetConfigFilePath(CONFIG_SUB_DIR_NAME))
	}

	return container.NewVBox(form, accountIdLabel)
}

func toBasicAuthorization(username, password string) string {
	auth := fmt.Sprintf("%v:%v", username, password)
	token := base64.StdEncoding.EncodeToString([]byte(auth))
	return fmt.Sprintf("Basic %v", token)
}

func getAccountIdByEmail(email, atlassianUrlBase, authorizationValue string) (string, error) {
	url := fmt.Sprintf("%v/rest/api/3/user/search", atlassianUrlBase)
	jsonResp, err := rest.SimpleExchange(rest.METHOD_GET, url, authorizationValue, map[string]string{"query": email}, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get Atlassian accountId by email %v", err)
	}
	result := []map[string]any{}
	err = json.Unmarshal([]byte(jsonResp), &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal json response %v, json: %v", err, jsonResp)
	}

	if len(result) == 1 {
		return result[0]["accountId"].(string), nil
	} else {
		return "", fmt.Errorf("expecting result json to have exactly one user but found %v", len(result))
	}
}

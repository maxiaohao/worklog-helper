package ui

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
	"worklog-helper/internal/data"
	"worklog-helper/pkg/rest"
	"worklog-helper/pkg/timefmt"
	"worklog-helper/pkg/yamlconfig"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Issue struct {
	key, id, summary string
}

const MAX_ISSUE_SUMMARY_DISPLAY_CHARS = 75

func NewTrackTimeContainer(w fyne.Window) fyne.CanvasObject {
	issueKeyInput := widget.NewEntry()
	issueKeyInput.SetPlaceHolder("Enter Jira (eg., ISSUE-1234)")

	issueLink := widget.NewHyperlink("", nil)

	recordIcon := theme.MediaRecordIcon()
	stopIcon := theme.MediaStopIcon()

	ticker := time.NewTicker(time.Second)
	ticker.Stop()
	started := make(chan bool)
	stopped := make(chan bool)

	startStopButton := widget.NewButtonWithIcon("Start", recordIcon, func() {})
	startStopButton.Disable()

	verifyButton := widget.NewButtonWithIcon("Verify", theme.SearchIcon(), func() {
		issueLink.SetText("")
		issueLink.SetURL(nil)
		startStopButton.Disable()

		settings := data.Settings{}
		err := yamlconfig.Read(&settings, CONFIG_SUB_DIR_NAME)
		if err != nil && len(settings.AccountId) == 0 {
			dialog.ShowError(fmt.Errorf("please initialise settings first"), w)
			return
		}

		issueKey := strings.TrimSpace(issueKeyInput.Text)
		if len(issueKey) == 0 || len(issueKey) > 20 {
			dialog.ShowError(fmt.Errorf("please enter a valid Jira number"), w)
			return
		}
		issue, err := getIssueByKey(issueKey, &settings)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		issueDisplaySummary := fmt.Sprintf("[%v] %v", issue.key, issue.summary)
		if len(issueDisplaySummary) > MAX_ISSUE_SUMMARY_DISPLAY_CHARS {
			issueDisplaySummary = issueDisplaySummary[:MAX_ISSUE_SUMMARY_DISPLAY_CHARS] + "..."
		}

		issueLink.SetText(issueDisplaySummary)
		issueLink.SetURL(parseURL(fmt.Sprintf("%v/browse/%v", settings.AtlassianSettings.UrlBase, issue.key)))

		startStopButton.Enable()
	})

	workNoteInput := widget.NewMultiLineEntry()
	workNoteInput.SetPlaceHolder("Work notes (optional, 100 chars max)")

	startTimeLabel := widget.NewLabel("")
	elspsedTimeLabel := widget.NewLabel("")

	var startTime time.Time

	go func() {
		for {
			select {
			case <-started:
				startTime = time.Now()
				startTimeLabel.SetText(fmt.Sprintf("Started at: %v", startTime.Format(time.RFC3339)))
				elspsedTimeLabel.SetText(timefmt.ElapsedTime(startTime))
				ticker.Reset(time.Second)
			case <-stopped:
				ticker.Stop()
				elspsedTimeLabel.SetText(elspsedTimeLabel.Text + " (Stopped)")
			case <-ticker.C:
				elspsedTimeLabel.SetText(timefmt.ElapsedTime(startTime))
			}
		}
	}()

	startStopButton.OnTapped = func() {
		if startStopButton.Icon == recordIcon {
			started <- true
			startStopButton.SetText("Stop")
			startStopButton.SetIcon(stopIcon)
			issueKeyInput.Disable()
			verifyButton.Disable()
		} else {
			stopped <- true
			startStopButton.SetText("Start")
			startStopButton.SetIcon(recordIcon)
			issueKeyInput.Enable()
			verifyButton.Enable()

			settings := data.Settings{}
			err := yamlconfig.Read(&settings, CONFIG_SUB_DIR_NAME)
			if err != nil && len(settings.AccountId) == 0 {
				dialog.ShowError(fmt.Errorf("please initialise settings first"), w)
				return
			}
			issue, err := getIssueByKey(strings.TrimSpace(issueKeyInput.Text), &settings)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if err = createWorklog(issue.id, workNoteInput.Text, startTime, &settings); err != nil {
				dialog.ShowError(err, w)
			} else {
				dialog.ShowInformation("Done", "A worklog has been sent to Tempo successfully!", w)
			}
		}
	}

	// TODO:
	// compactModeCheckbox := widget.NewCheck("Compact Mode", func(checked bool) {
	// 	var compactModeDesc string
	// 	if checked {
	// 		compactModeDesc = "ON"
	// 	} else {
	// 		compactModeDesc = "OFF"
	// 	}
	// 	fyne.CurrentApp().SendNotification(&fyne.Notification{
	// 		Title:   "Worklog Helper",
	// 		Content: fmt.Sprintf("Compact Mode: %v", compactModeDesc),
	// 	})
	// })

	vbox := container.NewVBox(
		container.NewGridWithColumns(2, issueKeyInput, verifyButton),
		issueLink,
		workNoteInput,
		startStopButton,
		startTimeLabel,
		elspsedTimeLabel,
		// TODO: compactModeCheckbox,
	)

	return vbox
}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func getIssueByKey(issueKey string, settings *data.Settings) (*Issue, error) {
	url := fmt.Sprintf("%v/rest/api/3/issue/%v", settings.AtlassianSettings.UrlBase, issueKey)
	jsonResp, err := rest.SimpleExchange(rest.METHOD_GET, url, settings.AtlassianSettings.AuthorizationValue, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot to get Jira '%v' (%v)", issueKey, err)
	}

	result := map[string]any{}
	err = json.Unmarshal([]byte(jsonResp), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response %v, json: %v", err, jsonResp)
	}

	issue := Issue{
		key:     result["key"].(string),
		id:      result["id"].(string),
		summary: strings.TrimSpace(result["fields"].(map[string]any)["summary"].(string)),
	}

	return &issue, nil
}

func createWorklog(issueId, workNote string, startTime time.Time, settings *data.Settings) error {
	duration := time.Since(startTime)
	if duration.Seconds() < 10 {
		return fmt.Errorf("duration is too short (less than 10 seconds), no worklog has been sent")
	}

	startDateLocal := startTime.Format("1999-12-31")
	startTimeLocal := startTime.Format("23:12:34")

	issueIdInt, _ := strconv.Atoi(issueId)
	reqBody := map[string]any{
		"authorAccountId":  settings.AccountId,
		"issueId":          issueIdInt,
		"startDate":        startDateLocal,
		"startTime":        startTimeLocal,
		"timeSpentSeconds": duration.Seconds(),
		"billableSeconds":  duration.Seconds(),
		"description":      strings.TrimSpace(workNote),
	}

	url := fmt.Sprintf("%v/4/worklogs", settings.TempoSettings.UrlBase)

	_, err := rest.SimpleExchange(rest.METHOD_POST, url, settings.TempoSettings.AuthorizationValue, nil, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create worklog %v", err)
	}
	return nil
}

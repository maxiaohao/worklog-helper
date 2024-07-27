package ui

import (
	"fmt"
	"net/url"
	"time"
	"worklog-helper/pkg/timefmt"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewTrackTimeContainer(w fyne.Window) fyne.CanvasObject {
	issueKeyInput := widget.NewEntry()
	issueKeyInput.SetPlaceHolder("Jira ID (eg., ISSUE-1234)")

	verifyButton := widget.NewButtonWithIcon("Validate", theme.SearchIcon(), func() {
		// TODO: enable/disable widgets...
	})

	issueLink := widget.NewHyperlink("ISSUE-1234 TODO", parseURL("https://foo.bar"))

	workNoteText := widget.NewMultiLineEntry()
	workNoteText.SetPlaceHolder("Work notes (optional)")

	startTimeLabel := widget.NewLabel("")
	elspsedTimeLabel := widget.NewLabel("")

	ticker := time.NewTicker(time.Second)
	ticker.Stop()
	started := make(chan bool)
	stopped := make(chan bool)

	go func() {
		var startTime time.Time
		for {
			select {
			case <-started:
				startTime = time.Now()
				startTimeLabel.SetText(fmt.Sprintf("Started at: %v", startTime.Format(time.RFC3339)))
				elspsedTimeLabel.SetText(timefmt.ElapsedTime(startTime))
				ticker.Reset(time.Second)
			case <-stopped:
				ticker.Stop()
				elspsedTimeLabel.SetText(elspsedTimeLabel.Text + " (Stopped, added to pending worklogs)")
			case <-ticker.C:
				elspsedTimeLabel.SetText(timefmt.ElapsedTime(startTime))
			}
		}
	}()

	compactModeCheckbox := widget.NewCheck("Compact Mode", func(checked bool) {
		var compactModeDesc string
		if checked {
			compactModeDesc = "ON"
		} else {
			compactModeDesc = "OFF"
		}
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Worklog Helper",
			Content: fmt.Sprintf("Compact Mode: %v", compactModeDesc),
		})
	})

	recordIcon := theme.MediaRecordIcon()
	stopIcon := theme.MediaStopIcon()

	startStopButton := widget.NewButtonWithIcon("Start", recordIcon, func() {})
	startStopButton.OnTapped = func() {
		if startStopButton.Icon == recordIcon {
			started <- true
			startStopButton.SetText("Stop")
			startStopButton.SetIcon(stopIcon)
		} else {
			stopped <- true
			startStopButton.SetText("Start")
			startStopButton.SetIcon(recordIcon)
		}
	}

	vbox := container.NewVBox(
		container.NewGridWithColumns(2, issueKeyInput, verifyButton),
		issueLink,
		workNoteText,
		startStopButton,
		startTimeLabel,
		elspsedTimeLabel,
		compactModeCheckbox,
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

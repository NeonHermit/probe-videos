package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/NeonHermit/probe-videos/pkg/analyzer"
	"github.com/NeonHermit/probe-videos/pkg/utils"
)

var selectedDir string

func main() {
	a := app.New()
	w := a.NewWindow("VideoProber")
	w.Resize(fyne.NewSize(600, 600))

	logText := widget.NewMultiLineEntry()
	logWindow := createLogWindow(a, logText)

	selectFolderButton := widget.NewButton("Select Folder", func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err == nil && dir != nil {
				selectedDir = dir.Path()
			}
		}, w)
	})

	startButton := widget.NewButton("Start Analysis", func() {
		if selectedDir == "" {
			utils.LogMessage("No directory selected.", logText)
			return
		}

		dialog.ShowFileSave(func(save fyne.URIWriteCloser, err error) {
			if err != nil || save == nil {
				utils.LogMessage("Failed to select output file for saving.", logText)
				return
			}
			csvFilePath := save.URI().Path()

			analyzer.RunAnalysis(selectedDir, csvFilePath, func(message string) {
				utils.LogMessage(message, logText)
			})
		}, w)
	})

	showLogWindowButton := widget.NewButton("Show Logs", func() {
		logWindow.Show()
	})

	w.SetContent(container.NewVBox(
		selectFolderButton,
		startButton,
		showLogWindowButton,
	))

	w.ShowAndRun()
}

func createLogWindow(app fyne.App, logText *widget.Entry) fyne.Window {
	logWindow := app.NewWindow("Logs")
	logWindow.SetContent(container.NewVScroll(logText))
	logWindow.Resize(fyne.NewSize(400, 300))
	return logWindow
}

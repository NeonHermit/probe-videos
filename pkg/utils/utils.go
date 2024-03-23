package utils

import (
	"fyne.io/fyne/v2/widget"
	"strings"
)

var supportedExtensions = []string{".mp4", ".avi", ".mkv"}

func IsSupportedFileType(fileName string) bool {
	for _, ext := range supportedExtensions {
		if strings.HasSuffix(strings.ToLower(fileName), ext) {
			return true
		}
	}
	return false
}

func LogMessage(message string, logText *widget.Entry) {
	logText.SetText(logText.Text + message + "\n")
}

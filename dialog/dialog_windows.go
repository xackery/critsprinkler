package dialog

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/critsprinkler/tracker"
	"github.com/xackery/wlk/walk"
)

// MsgBox displays a message box with the given title and message.
func MsgBox(title, message string) error {
	fmt.Println(message)
	// This function will display a message box with the given title and message.
	// It will block the main thread until the user closes the message box.
	// This function is only available on Windows.
	ret := walk.MsgBox(nil, title, message, walk.MsgBoxIconInformation)
	if ret != walk.DlgCmdNone {
		return nil
	}

	return nil
}

// FileDialogBox displays a file dialog box for selecting a file.
func FileDialogBox(path string) (string, error) {
	var err error

	dia := new(walk.FileDialog)

	curDir := "."
	if path != "" {
		curDir = filepath.Dir(path)
	} else {
		curDir, err = os.Getwd()
		if err != nil {
			curDir = "."
		}
	}

	dia.FilePath = curDir
	dia.Filter = "Log Files (*.txt)|eqlog_*.txt"
	dia.Title = "Select Log File"

	ok, err := dia.ShowOpen(nil)
	if err != nil {
		return "", fmt.Errorf("showOpen: %w", err)
	}

	if !ok {
		return "", fmt.Errorf("cancelled")
	}

	path = dia.FilePath
	err = tracker.SetNewPath(path)
	if err != nil {
		return "", fmt.Errorf("setNewPath: %w", err)
	}

	fmt.Println("Selected file:", path)
	return path, nil
}

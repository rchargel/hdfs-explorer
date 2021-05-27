package ui

import (
	"fyne.io/fyne/v2"
)

var fileRepoDialog *FileRepoManagerDialog = NewFileRepoManagerDialog()

func MakeMainMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		makeConnectionMenu(),
	)
}

func makeConnectionMenu() *fyne.Menu {
	manageConnections := &fyne.MenuItem{
		Label: "Manage Connections",
		Action: func() {
			fileRepoDialog.Open()
		},
	}

	quit := &fyne.MenuItem{
		Label:  "Quit",
		Action: onClosedFunc,
	}

	return &fyne.Menu{
		Label: "Connections",
		Items: []*fyne.MenuItem{
			manageConnections,
			quit,
		},
	}

}

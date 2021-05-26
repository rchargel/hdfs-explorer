package ui

import (
	"fyne.io/fyne/v2"
)

func MakeMainMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		makeConnectionMenu(),
	)
}

func makeConnectionMenu() *fyne.Menu {
	manageConnections := &fyne.MenuItem{
		Label: "Manage Connections",
		Action: func() {
			OpenFileSystemRepoManager()
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

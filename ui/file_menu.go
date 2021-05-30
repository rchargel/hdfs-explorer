package ui

import (
	"fyne.io/fyne/v2"
	"github.com/rchargel/hdfs-explorer/files"
)

var fileRepoDialog *FileRepoManagerDialog = NewFileRepoManagerDialog()

type OpenClientConnection func(client files.FileSystemClient)

func MakeMainMenu(clientConnectionFunc OpenClientConnection) *fyne.MainMenu {
	return fyne.NewMainMenu(
		makeConnectionMenu(clientConnectionFunc),
	)
}

func makeConnectionMenu(clientConnectionFunc OpenClientConnection) *fyne.Menu {
	manageConnections := &fyne.MenuItem{
		Label: "Manage Connections",
		Action: func() {
			fileRepoDialog.Open(OnNewConnection(clientConnectionFunc))
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

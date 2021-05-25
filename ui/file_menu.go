package ui

import (
	"fyne.io/fyne/v2"
	"github.com/rchargel/hdfs-explorer/log"
)

func MakeMainMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		makeConnectionMenu(),
	)
}

func makeConnectionMenu() *fyne.Menu {
	newConnection := &fyne.MenuItem{
		Label: "New Connection",
		Action: func() {
			log.Info.Println("Creating new connection")
		},
	}

	openConnection := &fyne.MenuItem{
		Label: "Open Connection",
		Action: func() {
			log.Info.Println("Opening connection")
		},
	}

	quit := &fyne.MenuItem{
		Label: "Quit",
		Action: func() {
			ShowConfirm(
				"Quit?",
				"Are you sure you would like to exit?",
				func() {
					log.Info.Print("Exiting application")
					Window.Close()
				},
			)
		},
	}

	return &fyne.Menu{
		Label: "Connections",
		Items: []*fyne.MenuItem{
			newConnection,
			openConnection,
			quit,
		},
	}

}

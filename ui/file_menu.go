package ui

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/rchargel/hdfs-explorer/log"
)

func MakeMainMenu(w fyne.Window) *fyne.MainMenu {
	return fyne.NewMainMenu(
		makeConnectionMenu(w),
	)
}

func makeConnectionMenu(w fyne.Window) *fyne.Menu {
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
			dialog.NewConfirm(
				"Quit",
				"Are you sure you wish to quit?",
				func(b bool) {
					if b {
						log.Info.Println("Exiting application")
						os.Exit(0)
					}
				},
				w,
			).Show()
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

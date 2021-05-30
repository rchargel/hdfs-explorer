package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/rchargel/hdfs-explorer/log"
)

var (
	Application fyne.App
	Window      fyne.Window

	onClosedFunc func() = func() {
		ShowConfirm("Quit?", "Are you certain you wish to exit?", func() {
			log.Info.Print("Exiting application by user request")
			fileRepoDialog.Close()
			closeAllBrowsers()
			Window.Close()
		})
	}
)

func init() {
	Application = app.New()
	Application.Settings().SetTheme(theme.DarkTheme())

	Window = Application.NewWindow("HDFS Explorer")

	fileBrowserTabs := NewFileBrowserTab()
	Window.SetMainMenu(MakeMainMenu(fileBrowserTabs.AddConnection))
	Window.SetIcon(resourceLogoPng)
	Window.SetContent(fileBrowserTabs.Container())
	Window.Resize(fyne.NewSize(700, 500))

	Window.CenterOnScreen()
	Window.SetCloseIntercept(onClosedFunc)
}

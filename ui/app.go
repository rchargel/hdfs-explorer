package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func CreateApp() fyne.Window {
	app := app.New()
	app.Settings().SetTheme(theme.DarkTheme())

	win := app.NewWindow("HDFS Explorer")
	win.SetMainMenu(MakeMainMenu(win))

	win.Resize(fyne.NewSize(700, 400))
	win.CenterOnScreen()

	return win
}

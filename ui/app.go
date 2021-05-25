package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var (
	Application fyne.App
	Window      fyne.Window
)

func init() {
	Application = app.New()
	Application.Settings().SetTheme(theme.DarkTheme())

	Window = Application.NewWindow("HDFS Explorer")
	Window.SetMainMenu(MakeMainMenu())
	Window.Resize(fyne.NewSize(700, 400))
	Window.CenterOnScreen()
}

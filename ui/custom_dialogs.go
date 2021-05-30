package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/rchargel/hdfs-explorer/log"
)

func ShowFatalError(err error) {
	errorDialog := dialog.NewError(err, Window)
	errorDialog.SetOnClosed(func() {
		log.Error.Fatal(err)
	})
	errorDialog.Show()
}

func ShowConfirm(title, text string, callback func()) {
	dialog.NewConfirm(
		title,
		text,
		func(b bool) {
			if b {
				callback()
			}
		},
		Window,
	).Show()
}

func NewCustomDialog(title string, content fyne.CanvasObject, buttons ...*widget.Button) fyne.Window {
	var win fyne.Window

	closeButton := widget.NewButton(
		"Close",
		func() { win.Hide() },
	)
	buttonList := make([]fyne.CanvasObject, len(buttons)+2)
	buttonList[0] = layout.NewSpacer()
	buttonList[1] = closeButton
	idx := 2
	for _, button := range buttons {
		buttonList[idx] = button
		idx++
	}
	buttonPanel := container.NewHBox(buttonList...)
	container := container.NewBorder(
		nil,
		buttonPanel,
		nil,
		nil,
		content,
	)
	win = Application.NewWindow(title)
	win.SetContent(container)
	return win
}

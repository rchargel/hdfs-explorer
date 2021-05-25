package ui

import "fyne.io/fyne/v2/dialog"

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

package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/rchargel/hdfs-explorer/files"
)

type FileBrowser struct {
	client files.FileSystemClient
	files  binding.StringList
	path   string
}

func NewFileBrowser(client files.FileSystemClient) *FileBrowser {
	b := &FileBrowser{client, binding.NewStringList(), "/"}
	b.updateList()
	return b
}

func (b *FileBrowser) Tab() *container.TabItem {
	// TODO implement custom widget for file exploration
	l := widget.NewListWithData(
		b.files,
		func() fyne.CanvasObject {
			t := canvas.NewText("template", color.Black)
			return t
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			t := co.(*canvas.Text)
			str, _ := di.(binding.String).Get()
			t.Text = str
		},
	)

	return container.NewTabItem(
		b.client.Name(),
		container.NewMax(l),
	)
}

func (b *FileBrowser) Close() {
	b.client.Close()
}

func (b *FileBrowser) updateList() {
	go func() {
		orig := make([]string, 0)
		if b.path != "/" {
			orig = append(orig, "..")
		}
		b.files.Set(orig)
		info, _ := b.client.List(b.path)
		for _, file := range info {
			b.files.Append(file.Name())
		}
	}()
}

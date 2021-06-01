package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/rchargel/hdfs-explorer/files"
)

type FileBrowser struct {
	name  string
	files HdfsFileInfoList
	path  string
}

func NewFileBrowser(client files.FileSystemClient) *FileBrowser {
	b := &FileBrowser{client.Name(), NewHdfsFileInfoList(client), "/"}
	b.update()
	return b
}

func (b *FileBrowser) Tab() *container.TabItem {
	l := NewFileDirectoryList(b.files, b.updatePath)
	scroll := container.NewVScroll(l.container)
	c := container.NewHBox(scroll, layout.NewSpacer())
	return container.NewTabItem(b.name, c)
}

func (b *FileBrowser) Close() {
	b.files.Close()
}

func (b *FileBrowser) update() {
	go func() {
		b.files.UpdatePath(b.path)
	}()
}

func (b *FileBrowser) updatePath(name string) {
	if name == ".." {
		b.path = files.Parent(b.path)
	} else {
		b.path = files.Join(b.path, name)
	}
	b.update()
}

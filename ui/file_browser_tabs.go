package ui

import (
	"fyne.io/fyne/v2/container"
	"github.com/rchargel/hdfs-explorer/files"
	"github.com/rchargel/hdfs-explorer/log"
)

var closeAllBrowsers func()

type FileBrowserTabs struct {
	browsers []*FileBrowser
	tabs     *container.AppTabs
}

func NewFileBrowserTab() *FileBrowserTabs {
	fileBrowserTabs := &FileBrowserTabs{make([]*FileBrowser, 0), container.NewAppTabs()}
	fileBrowserTabs.tabs.SetTabLocation(container.TabLocationTop)

	closeAllBrowsers = func() {
		fileBrowserTabs.CloseAll()
	}
	return fileBrowserTabs
}

func (f *FileBrowserTabs) Container() *container.AppTabs {
	if f.tabs == nil {
		f.tabs = container.NewAppTabs()
		f.tabs.SetTabLocation(container.TabLocationTop)
	}
	return f.tabs
}

func (f *FileBrowserTabs) AddConnection(client files.FileSystemClient) {
	log.Info.Printf("Open Connection %v\n", client.Name())
	if tabIndex := f.getConnection(client.Name()); tabIndex < 0 {
		browser := NewFileBrowser(client)
		f.browsers = append(f.browsers, browser)
		f.tabs.Append(browser.Tab())
		f.tabs.SelectTabIndex(len(f.browsers) - 1)
	} else {
		f.tabs.SelectTabIndex(tabIndex)
	}
}

func (f *FileBrowserTabs) CloseAll() {
	for _, b := range f.browsers {
		b.Close()
	}
}

func (f *FileBrowserTabs) getConnection(client string) int {
	idx := 0
	for _, conn := range f.browsers {
		if conn.name == client {
			return idx
		}
		idx++
	}
	return -1
}

package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/rchargel/hdfs-explorer/files"
	"github.com/rchargel/hdfs-explorer/log"
)

var fileSystemRepo files.FileSystemRepository

func init() {
	repository, err := files.GetFileSystemRepository()

	if err != nil {
		log.Error.Fatal(err)
	}
	fileSystemRepo = repository
}

func fileSystemNames() *[]string {
	list, err := fileSystemRepo.List()
	if err != nil {
		ShowFatalError(err)
	}
	names := make([]string, len(list))
	for i, fileSystem := range list {
		names[i] = fileSystem.Name
	}
	return &names
}

func OpenFileSystemRepoManager() {
	fileSystems := binding.BindStringList(fileSystemNames())
	name := binding.BindString(nil)
	description := binding.BindString(nil)
	addresses := binding.BindString(nil)

	fsList := fileSystemList(fileSystems)
	form := createForm(name, description, addresses)
	dialog := container.NewHBox(fsList, form)

	ShowCustomDialog(
		"Connections",
		dialog,
		widget.NewButton("New", func() {
			val := fmt.Sprintf("Next %d", fileSystems.Length()+1)
			fileSystems.Append(val)
		}),
	)
}

func createForm(nameData, descriptionData, addressData binding.String) *fyne.Container {
	nameLabel := widget.NewLabel("Name: ")
	nameEntry := widget.NewEntryWithData(nameData)
	descriptionLabel := widget.NewLabel("Description: ")
	descriptionEntry := widget.NewEntryWithData(descriptionData)
	descriptionEntry.MultiLine = true
	addressLabel := widget.NewLabel("Addresses (new line sperated):")
	addressEntry := widget.NewEntryWithData(addressData)
	addressEntry.MultiLine = true

	return container.NewVBox(
		nameLabel,
		nameEntry,
		descriptionLabel,
		descriptionEntry,
		addressLabel,
		addressEntry,
	)
}

func fileSystemList(data binding.DataList) *fyne.Container {
	list := widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			return widget.NewLabel("Connections")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	return container.New(layout.NewMaxLayout(), list)
}

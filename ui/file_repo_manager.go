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

type FileRepoManagerDialog struct {
	Title          string
	popup          *widget.PopUp
	fileSystemRepo files.FileSystemRepository
	fsList         binding.DataList
	name           binding.String
	description    binding.String
	addresses      binding.String
}

func NewFileRepoManagerDialog() *FileRepoManagerDialog {
	repository, err := files.GetFileSystemRepository()

	if err != nil {
		log.Error.Fatal(err)
	}
	dialog := &FileRepoManagerDialog{}
	dialog.Title = "Connections"
	dialog.fileSystemRepo = repository
	fileSystems := binding.BindStringList(dialog.fileSystemNames())
	dialog.fsList = fileSystems
	dialog.name = binding.BindString(nil)
	dialog.description = binding.BindString(nil)
	dialog.addresses = binding.BindString(nil)

	return dialog
}

func (d *FileRepoManagerDialog) fileSystemNames() *[]string {
	list, err := d.fileSystemRepo.List()
	if err != nil {
		ShowFatalError(err)
	}
	names := make([]string, len(list))
	for i, fileSystem := range list {
		names[i] = fileSystem.Name
	}
	return &names
}

func (f *FileRepoManagerDialog) Open() {
	if f.popup == nil {
		fsList := f.fileSystemList()
		form := f.createForm()
		content := container.NewHBox(fsList, form)
		newButton := widget.NewButton("New", func() {
			val := fmt.Sprintf("Next %d", f.fsList.Length()+1)
			f.fsList.(binding.ExternalStringList).Append(val)
		})
		f.popup = NewCustomDialog(
			f.Title,
			content,
			newButton,
		)
	}
	f.popup.Show()
}

func (f *FileRepoManagerDialog) createForm() *fyne.Container {
	nameLabel := widget.NewLabel("Name: ")
	nameEntry := widget.NewEntryWithData(f.name)
	descriptionLabel := widget.NewLabel("Description: ")
	descriptionEntry := widget.NewEntryWithData(f.description)
	descriptionEntry.MultiLine = true
	addressLabel := widget.NewLabel("Addresses (new line sperated):")
	addressEntry := widget.NewEntryWithData(f.addresses)
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

func (f *FileRepoManagerDialog) fileSystemList() *fyne.Container {
	list := widget.NewListWithData(
		f.fsList,
		func() fyne.CanvasObject {
			return widget.NewLabel("Connections")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	return container.New(layout.NewMaxLayout(), list)
}

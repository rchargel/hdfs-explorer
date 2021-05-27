package ui

import (
	"errors"
	"fmt"
	"regexp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/rchargel/hdfs-explorer/files"
	"github.com/rchargel/hdfs-explorer/log"
)

type FileRepoManagerDialog struct {
	Title          string
	popup          *widget.PopUp
	listUI         *ListWithListener
	fileSystemRepo files.FileSystemRepository
	fsList         binding.StringList
	name           binding.String
	description    binding.String
	addresses      binding.String
}

func NewFileRepoManagerDialog() *FileRepoManagerDialog {
	repository, err := files.GetFileSystemRepository()

	if err != nil {
		log.Error.Fatal(err)
	}
	dialog := &FileRepoManagerDialog{
		Title:          "Connections",
		popup:          nil,
		fileSystemRepo: repository,
		fsList:         binding.NewStringList(),
		name:           binding.NewString(),
		description:    binding.NewString(),
		addresses:      binding.NewString(),
	}

	dialog.fsList.Set(dialog.fileSystemNames())

	return dialog
}

func (d *FileRepoManagerDialog) addNew() {
	d.listUI.ClearSelected()
	d.name.Set("")
	d.description.Set("")
	d.addresses.Set("")
}

func (d *FileRepoManagerDialog) save() error {
	name, _ := d.name.Get()
	desc, _ := d.description.Get()
	addr, _ := d.addresses.Get()

	re := regexp.MustCompile(`\s+`)
	addresses := re.FindAllString(addr, -1)

	if len(name) == 0 || len(desc) == 0 || addresses == nil {
		return errors.New("Missing field data")
	}
	fs := files.FileSystem{
		Name:        name,
		Description: desc,
		Addresses:   addresses,
	}

	err := d.fileSystemRepo.Store(fs)
	if err == nil {
		d.fsList.Set(d.fileSystemNames())
	}
	d.addNew()
	return err
}

func (d *FileRepoManagerDialog) fileSystemNames() []string {
	list, err := d.fileSystemRepo.List()
	if err != nil {
		ShowFatalError(err)
	}
	names := make([]string, len(list))
	for i, fileSystem := range list {
		names[i] = fileSystem.Name
	}
	return names
}

func (d *FileRepoManagerDialog) Open() {
	if d.popup == nil {
		fsList := d.fileSystemList()
		form := d.createForm()
		content := container.NewHBox(fsList, form)
		newButton := widget.NewButton("New", func() {
			d.addNew()
		})
		saveButton := widget.NewButton("Save", func() {
			if err := d.save(); err != nil {
				dialog.ShowError(err, Window)
			}
		})
		d.popup = NewCustomDialog(
			d.Title,
			content,
			newButton,
			saveButton,
		)
	}
	d.popup.Show()
}

func (d *FileRepoManagerDialog) createForm() *fyne.Container {
	nameLabel := widget.NewLabel("Name: ")
	nameEntry := widget.NewEntryWithData(d.name)
	descriptionLabel := widget.NewLabel("Description: ")
	descriptionEntry := widget.NewEntryWithData(d.description)
	descriptionEntry.MultiLine = true
	addressLabel := widget.NewLabel("Addresses (new line sperated):")
	addressEntry := widget.NewEntryWithData(d.addresses)
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

func (d *FileRepoManagerDialog) fileSystemList() *fyne.Container {
	d.listUI = NewSelectableList(d.fsList)

	d.listUI.AddListener(func(event interface{}) {
		label := event.(string)
		fmt.Printf("Selected label %v\n", label)
	})

	return container.New(layout.NewMaxLayout(), d.listUI.CanvasObject())
}

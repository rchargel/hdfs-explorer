package ui

import (
	"errors"
	"regexp"
	"strings"

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
	popup          fyne.Window
	listUI         *ListWithListener
	fileSystemRepo files.FileSystemRepository
	fsList         binding.StringList
	name           binding.String
	description    binding.String
	addresses      binding.String
}

type OnNewConnection func(client files.FileSystemClient)

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

func (d *FileRepoManagerDialog) edit(label string) error {
	fs, err := d.fileSystemRepo.FindByName(label)

	if err != nil {
		return err
	}
	d.name.Set(fs.Name)
	d.description.Set(fs.Description)
	d.addresses.Set(strings.Join(fs.Addresses, "\n"))
	return nil
}

func (d *FileRepoManagerDialog) new() {
	d.listUI.ClearSelected()
	d.name.Set("")
	d.description.Set("")
	d.addresses.Set("")
}

func (d *FileRepoManagerDialog) save() error {
	name, _ := d.name.Get()
	desc, _ := d.description.Get()
	addr, _ := d.addresses.Get()

	if len(name) == 0 || len(desc) == 0 || len(addr) == 0 {
		return errors.New("Missing field data")
	}

	re := regexp.MustCompile(`\s+`)
	split := re.Split(addr, -1)
	addresses := make([]string, 0)
	if split == nil {
		addresses = []string{strings.TrimSpace(addr)}
	} else {
		for _, v := range split {
			addresses = append(addresses, strings.TrimSpace(v))
		}
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
	d.new()
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

func (d *FileRepoManagerDialog) Close() {
	if d.popup != nil {
		d.popup.Close()
	}
}

func (d *FileRepoManagerDialog) Open(onNewConnection OnNewConnection) {
	if d.popup == nil {
		fsList := d.fileSystemList(onNewConnection)
		form := d.createForm()
		content := container.NewBorder(nil, nil, nil, form, fsList, widget.NewSeparator())
		newButton := widget.NewButton("New", func() {
			d.new()
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

func (d *FileRepoManagerDialog) fileSystemList(onNewConnection OnNewConnection) *fyne.Container {
	d.listUI = NewSelectableList(d.fsList)
	d.listUI.list.Resize(fyne.NewSize(400, 200))

	d.listUI.AddListener(func(event Event) {
		label := event.EventSource
		if event.EventType == "click" {
			if err := d.edit(label); err != nil {
				dialog.ShowError(err, Window)
			}
		} else if event.EventType == "dblclick" {
			fs, _ := d.fileSystemRepo.FindByName(event.EventSource)
			client, err := fs.Connect()
			d.popup.Hide()
			if err != nil {
				dialog.ShowError(err, Window)
			} else {
				onNewConnection(client)
			}
		}
	})

	return container.New(layout.NewMaxLayout(), d.listUI.CanvasObject())
}

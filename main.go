package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/rchargel/hdfs-explorer/files"
)

func main() {
	app := app.New()
	win := app.NewWindow("Hello")

	hello := widget.NewLabel("Hello World!")
	win.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	fs := files.FileSystem{
		Name:        "My FS",
		Description: "A simple description",
		Addresses:   []string{"localhost:9000"},
		User:        "rchargel",
	}
	if conn, err := fs.Connect(); err == nil {
		values, _ := conn.List("/rmstate/FSRMStateRoot")
		for _, value := range values {
			fmt.Printf("Is Dir: %v \n", value.IsDir())
			fmt.Printf("Mod Time: %v \n", value.ModTime())
			fmt.Printf("Access Time: %v \n", value.AccessTime())
			fmt.Printf("Mode: %v \n", value.Mode())
			fmt.Printf("Name: %v \n", value.Name())
			fmt.Printf("Size: %v \n", value.Size())

			fmt.Printf("Owner: %v\n", value.Owner())
			fmt.Printf("Group: %v\n", value.OwnerGroup())
			fmt.Printf("Replication: %v\n", value.Replication())
			fmt.Printf("Block Size: %v\n", value.BlockSize())

		}
		conn.Close()
	} else {
		println(err.Error())
	}

	win.ShowAndRun()
}

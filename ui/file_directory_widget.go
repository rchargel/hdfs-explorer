package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rchargel/hdfs-explorer/files"
)

var sizes [7]float32 = [7]float32{}

var _ fyne.Widget = (*fileDirectoryRow)(nil)
var _ fyne.Tappable = (*fileDirectoryRow)(nil)
var _ fyne.DoubleTappable = (*fileDirectoryRow)(nil)
var _ desktop.Hoverable = (*fileDirectoryRow)(nil)

type FileDirectoryList struct {
	container *fyne.Container
}

func NewFileDirectoryList(data HdfsFileInfoList, openItem func(string)) *FileDirectoryList {
	list := &FileDirectoryList{
		container: container.New(layout.NewVBoxLayout()),
	}
	data.AddListener(binding.NewDataListener(func() {
		infos, _ := data.Get()
		objects := make([]fyne.CanvasObject, len(infos))
		for i, info := range infos {
			objects[i] = newFileDirectoryRow(info, list.deselectAll, openItem)
		}
		list.container.Objects = objects
		list.container.Refresh()
	}))
	return list
}

func (f *FileDirectoryList) deselectAll(name string) {
	for _, obj := range f.container.Objects {
		if obj.(*fileDirectoryRow).info.Name() != name {
			obj.(*fileDirectoryRow).selected = false
			obj.Refresh()
		}
	}
}

type fileDirectoryRow struct {
	widget.BaseWidget
	info  files.HdfsFileInfo
	child fyne.CanvasObject

	mode  *canvas.Text
	owner *canvas.Text
	group *canvas.Text
	size  *canvas.Text
	block *canvas.Text
	time  *canvas.Text
	name  *canvas.Text

	selected, hovered bool
	onTapped          func(name string)
	onDoubleTapped    func(name string)
}

type fileDirectoryRowRenderer struct {
	objects []fyne.CanvasObject

	row *fileDirectoryRow
}

func newFileDirectoryRow(info files.HdfsFileInfo, onTapped, onDoubleTapped func(name string)) *fileDirectoryRow {
	mode := newText(info.Mode().String(), fyne.TextAlignLeading)
	owner := newText(info.Owner(), fyne.TextAlignLeading)
	group := newText(info.OwnerGroup(), fyne.TextAlignLeading)
	size := newText(files.FormatBytes(uint64(info.Size())), fyne.TextAlignTrailing)
	block := newText(files.FormatBytes(info.BlockSize()), fyne.TextAlignTrailing)
	time := newText(info.ModTime().Format("Jan 02 15:04"), fyne.TextAlignLeading)
	name := newText(info.Name(), fyne.TextAlignLeading)

	r := &fileDirectoryRow{
		BaseWidget:     widget.BaseWidget{},
		info:           info,
		mode:           mode,
		owner:          owner,
		group:          group,
		size:           size,
		block:          block,
		time:           time,
		name:           name,
		onTapped:       onTapped,
		onDoubleTapped: onDoubleTapped,
		child: container.New(
			layout.NewGridLayoutWithColumns(7),
			mode, owner, group, size, block, time, name,
		),
	}
	r.ExtendBaseWidget(r)
	return r
}

func (r *fileDirectoryRow) setStyle(style fyne.TextStyle) {
	objects := r.child.(*fyne.Container).Objects
	for _, obj := range objects {
		obj.(*canvas.Text).TextStyle = style
	}
}

func (r *fileDirectoryRow) setColor(color color.Color) {
	objects := r.child.(*fyne.Container).Objects
	for _, obj := range objects {
		obj.(*canvas.Text).Color = color
	}
}

func (r *fileDirectoryRow) MinSize() fyne.Size {
	r.ExtendBaseWidget(r)
	return r.BaseWidget.MinSize()
}

func (r *fileDirectoryRow) Tapped(*fyne.PointEvent) {
	r.selected = !r.selected
	r.onTapped(r.info.Name())
	r.Refresh()
}

func (r *fileDirectoryRow) DoubleTapped(*fyne.PointEvent) {
	r.onDoubleTapped(r.info.Name())
}

func (r *fileDirectoryRow) MouseIn(*desktop.MouseEvent) {
	r.hovered = true
	r.Refresh()
}

func (r *fileDirectoryRow) MouseOut() {
	r.hovered = false
	r.Refresh()
}

func (r *fileDirectoryRow) MouseMoved(*desktop.MouseEvent) {}

func newText(text string, alignment fyne.TextAlign) *canvas.Text {
	return &canvas.Text{
		Text:      text,
		Alignment: alignment,
		Color:     theme.ForegroundColor(),
		TextStyle: fyne.TextStyle{Monospace: true},
		TextSize:  fyne.CurrentApp().Settings().Theme().Size("text") * float32(0.9),
	}
}

func (r *fileDirectoryRow) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)
	objects := []fyne.CanvasObject{r.child}

	return &fileDirectoryRowRenderer{objects, r}
}

func (r *fileDirectoryRowRenderer) Destroy() {
	// nothing done here
}

func (r *fileDirectoryRowRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *fileDirectoryRowRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}

func (r *fileDirectoryRowRenderer) MinSize() fyne.Size {
	objects := r.row.child.(*fyne.Container).Objects
	height := r.row.child.MinSize().Height

	for i, o := range objects {
		sizes[i] = max(sizes[i], o.MinSize().Width)
	}

	return fyne.NewSize(minWidth(), height+theme.Padding()*2)
}

func (r *fileDirectoryRowRenderer) Layout(size fyne.Size) {
	r.row.child.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	r.row.child.Resize(fyne.NewSize(size.Width+theme.Padding()*2, size.Height+theme.Padding()*2))

	x := theme.Padding()
	for i, obj := range r.row.child.(*fyne.Container).Objects {
		width := sizes[i]
		obj.Resize(fyne.NewSize(width, obj.MinSize().Height))
		obj.Move(fyne.NewPos(x, theme.Padding()))
		x += width + theme.Padding()
	}
}

func (r *fileDirectoryRowRenderer) Refresh() {
	if r.row.hovered {
		r.row.setColor(theme.PrimaryColor())
	} else {
		r.row.setColor(theme.ForegroundColor())
	}
	canvas.Refresh(r.row)
}

func minWidth() float32 {
	width := float32(0)
	for _, w := range sizes {
		width += w + theme.Padding()*2
	}
	return width
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

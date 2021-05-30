package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Event struct {
	EventType   string
	EventSource string
}

type EventListener func(event Event)

type ListWithListener struct {
	list     *widget.List
	data     binding.StringList
	selected int

	listeners []EventListener
}

type ClickableLabel struct {
	widget.Label

	listeners []EventListener
}

func NewClickableLabel(text string) *ClickableLabel {
	label := &ClickableLabel{}
	label.ExtendBaseWidget(label)
	label.SetText(text)
	label.listeners = make([]EventListener, 0)
	return label
}

func (l *ClickableLabel) AddListener(listener EventListener) {
	l.listeners = append(l.listeners, listener)
}

func (l *ClickableLabel) Tapped(_ *fyne.PointEvent) {
	for _, listener := range l.listeners {
		listener(Event{"click", l.Text})
	}
}

func (l *ClickableLabel) DoubleTapped(_ *fyne.PointEvent) {
	for _, listener := range l.listeners {
		listener(Event{"dblclick", l.Text})
	}
}

func NewSelectableList(data binding.StringList) *ListWithListener {
	var list *ListWithListener
	list = &ListWithListener{widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			clickableLabel := NewClickableLabel("template")
			clickableLabel.AddListener(list.labelClick)
			return clickableLabel
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			co.(*ClickableLabel).Bind(di.(binding.String))
		},
	), data, -1, make([]EventListener, 0)}
	return list
}

func (l *ListWithListener) labelClick(event Event) {
	if event.EventType == "dblclick" {
		for _, listener := range l.listeners {
			listener(event)
		}
	} else if event.EventType == "click" {
		label := event.EventSource

		idx := l.getIndex(label)
		if idx > -1 {
			l.list.Select(idx)
			if l.selected != idx {
				l.selected = idx
				for _, listener := range l.listeners {
					listener(event)
				}
			}
		} else {
			l.list.Unselect(l.selected)
			l.selected = -1
		}
	}
}

func (l *ListWithListener) getIndex(text string) int {
	values, _ := l.data.Get()

	for i, value := range values {
		if value == text {
			return i
		}
	}
	return -1
}

func (l *ListWithListener) CanvasObject() fyne.CanvasObject {
	return l.list
}

func (l *ListWithListener) AddListener(listener EventListener) {
	l.listeners = append(l.listeners, listener)
}

func (l *ListWithListener) ClearSelected() {
	if l.selected > -1 {
		l.list.Unselect(l.selected)
		l.selected = -1
	}
}

package views

import (
	"simple-git-terminal/events"

	"github.com/rivo/tview"
)

// BaseView defines the interface every reactive view should implement
type BaseView interface {
	Render() tview.Primitive
	Refresh()
	Subscribe(bus *events.Bus)
}

// BaseViewImpl provides a default implementation for BaseView
type BaseViewImpl struct {
	primitive tview.Primitive // The actual tview widget (Table, Box, etc.)
	bus       *events.Bus     // Event bus
	OnRefresh func()          // Optional refresh callback
}

// NewBaseView creates a new BaseViewImpl wrapping a tview primitive
func NewBaseView(primitive tview.Primitive) *BaseViewImpl {
	return &BaseViewImpl{
		primitive: primitive,
	}
}

// Render returns the wrapped primitive for use in the layout
func (b *BaseViewImpl) Render() tview.Primitive {
	return b.primitive
}

// Refresh triggers the refresh callback if set
func (b *BaseViewImpl) Refresh() {
	if b.OnRefresh != nil {
		b.OnRefresh()
	}
}

// Subscribe attaches the view to the event bus
func (b *BaseViewImpl) Subscribe(bus *events.Bus) {
	b.bus = bus
}

// Publish allows derived views to send events to the bus
func (b *BaseViewImpl) Publish(event events.Event) {
	if b.bus != nil {
		b.bus.Publish(event)
	}
}

// SetOnRefresh sets the refresh callback
func (b *BaseViewImpl) SetOnRefresh(cb func()) {
	b.OnRefresh = cb
}

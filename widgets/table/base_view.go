package widgets

import (
	"simple-git-terminal/constants"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Refreshable is implemented by views that can be refreshed (e.g., after file changes)
type Refreshable interface {
	Refresh()
}

type BaseTableView struct {
	*tview.Table
	SelectedRow int
	Selectable  bool
	OnRefresh   func()
}

// Constructor
func NewBaseTableView() *BaseTableView {
	table := tview.NewTable()

	table.
		SetFixed(1, 0).
		SetBackgroundColor(tcell.ColorDefault)

	return &BaseTableView{
		Table:       table,
		SelectedRow: -1,
		Selectable:  true,
	}
}

// Toggle row selection (UI feedback)
func (b *BaseTableView) SetSelectedRow(row int) {
	b.SelectedRow = row
}

func (b *BaseTableView) GetSelectedRow() int {
	return b.SelectedRow
}

func (b *BaseTableView) SetSelectableState(selectable bool) {
	b.Selectable = selectable
	b.SetSelectable(selectable, false)
}

func (b *BaseTableView) UpdateSelectedRow(row int) {
	if row < 0 {
		return
	}

	b.SelectedRow = row
	b.SetCell(row, 0, util.CellFormat(constants.ICON_SELECTED, tcell.ColorOrange))
}

func (b *BaseTableView) UpdateUnSelectedRow(row int) {
	if row < 0 {
		return
	}

	b.SelectedRow = row
	b.SetCell(row, 0, util.CellFormat("", tcell.ColorDefault))
}

func (b *BaseTableView) ClearTable() {
	b.Clear()
}

func (b *BaseTableView) Refresh() {
	b.ClearTable()
	b.SelectedRow = -1
	if b.OnRefresh != nil {
		b.OnRefresh()
	}
}

func (s *BaseTableView) SetOnRefresh(cb func()) {
	s.OnRefresh = cb
}

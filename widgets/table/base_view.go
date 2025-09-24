package widgets

import (
	"log"
	"simple-git-terminal/constants"
	"simple-git-terminal/util"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Selectable interface for views that support SetSelectable
type Selectable interface {
	SetSelectable(selectable bool, allowsMultiple bool)
}

type TableView interface {
	SetCell(row, col int, cell *tview.TableCell)
	SetLoadingCell(cell *tview.TableCell)
	GetRowCount() int
	GetColumnCount() int
}

// Refreshable is implemented by views that can be refreshed (e.g., after file changes)
type Refreshable interface {
	Refresh()
}

type BaseTableView struct {
	*tview.Table
	SelectedRow int
	OnRefresh   func()

	loading bool
}

// Constructor
func NewBaseTableView() *BaseTableView {
	table := tview.NewTable()

	table.
		SetFixed(1, 0).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitleAlign(tview.AlignLeft)

	return &BaseTableView{
		Table:       table,
		SelectedRow: -1,
		loading:     false,
	}
}

func (b *BaseTableView) SetCell(row, col int, cell *tview.TableCell) {
	// Block all updates except to the loading cell
	if b.loading {
		return
	}
	b.Table.SetCell(row, col, cell)
}

func (b *BaseTableView) SetLoadingCell(cell *tview.TableCell) {
	b.loading = true

	text := cell.Text
	if idx := strings.Index(text, "Loading..."); idx != -1 {
		text = text[:idx] // keep everything before "Loading..." // FIXME: Super hacky for now
	}
	cell.SetText(text).
		SetTextColor(tcell.ColorOrange)

	b.Table.SetCell(0, 0, cell)
}

func (b *BaseTableView) ClearLoading() {
	b.loading = false
}

func (b *BaseTableView) SetSelectable(selectable bool, allowsMultiple bool) {
	b.Table.SetSelectable(selectable, allowsMultiple)
}

// Toggle row selection (UI feedback)
func (b *BaseTableView) SetSelectedRow(row int) {
	b.SelectedRow = row
}

func (b *BaseTableView) GetSelectedRow() int {
	return b.SelectedRow
}

func (b *BaseTableView) SetSelectableState(selectable bool) {
	// Optional helper method to set selectable state easily
	b.SetSelectable(selectable, false)
}

func (b *BaseTableView) UpdateSelectedRow(row int) {
	if row < 0 {
		log.Println("[UpdateSelectedRow] Invalid row:", row)
		return
	}

	log.Println("[UpdateSelectedRow] Selecting row:", row)
	b.SelectedRow = row
	b.SetCell(row, 0, util.CellFormat(constants.ICON_SELECTED, tcell.ColorOrange))
}

func (b *BaseTableView) UpdateUnSelectedRow(row int) {
	if row < 0 {
		log.Println("[UpdateUnSelectedRow] Invalid row:", row)
		return
	}

	log.Println("[UpdateUnSelectedRow] Unselecting row:", row)
	// Clear selection icon from the row
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

func (b *BaseTableView) SetOnRefresh(cb func()) {
	b.OnRefresh = cb
}

package support

import (
	"log"
	widgets "simple-git-terminal/widgets/table"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Helper method to update borders of views
func UpdateFocusBorders(focusOrder []tview.Primitive, currentFocusIndex int, activeBorderColor tcell.Color) {
	for i, view := range focusOrder {
		// Check if the view has border-related methods
		if bordered, ok := view.(interface {
			SetBorder(bool) *tview.Box
			SetBorderColor(tcell.Color) *tview.Box
		}); ok {
			if i == currentFocusIndex {
				bordered.SetBorder(true).
					SetBorderColor(activeBorderColor)
			} else {
				bordered.SetBorder(true).
					SetBorderColor(tcell.ColorGrey)
			}
		}
	}
}

func UpdateView(targetView interface{}, content interface{}) {
	if targetView != nil {
		// Check the type of the target view (either Flex or TextView)
		switch v := targetView.(type) {
		case *tview.Flex:
			// If it's a Flex view, clear it and update it
			v.Clear()

			// Handle content based on its type
			switch c := content.(type) {
			case string:
				// If the content is a string, display it in a TextView
				textView := CreateTextviewComponent("", false).SetText(c)
				v.AddItem(textView, 0, 1, true)
			case tview.Primitive:
				// If the content is a tview.Primitive, add it directly
				v.AddItem(c, 0, 1, true)
			default:
				// Handle unsupported content types
				errorView := CreateTextviewComponent("", false).SetText("[red]Unsupported content type[-]")
				v.AddItem(errorView, 0, 1, true)
			}

		case *tview.TextView:
			// If it's a TextView, update the text directly
			switch c := content.(type) {
			case string:
				// If content is a string, update the TextView
				v.SetText(c)
			case tview.Primitive:
				// Handle case if content is another Primitive (optional)
				log.Println("Unsupported Primitive content for TextView")
			default:
				// Handle unsupported content types
				v.SetText("[red]Unsupported content type[-]")
			}

		case *tview.Table:
			switch c := content.(type) {
			case string:
				tcell := CreateTableCell(c, tcell.ColorDefault)
				v.SetCell(0, 0, tcell)
			case tview.Primitive:
				// Handle case if content is another Primitive (optional)
				log.Println("Unsupported Primitive content for TextView")
			default:
				tcell := CreateTableCell("[red]Unsupported content type[-]", tcell.ColorDefault)
				v.SetCell(0, 0, tcell)
			}
		default:
			// If it's neither Flex nor TextView, print an error
			log.Println("[red]Unsupported target view type[-]")
		}
	}
}

func SetTableSelectableIfFocused(focus tview.Primitive, table tview.Primitive, focusOrder []tview.Primitive, currentFocusIndex int) {
	if focusOrder[currentFocusIndex] == focus {
		if t, ok := table.(*tview.Table); ok {
			t.SetSelectable(true, false)
		}
	}
}

func SetTableSelectability(focusOrder []tview.Primitive, currentFocusIndex int, tableMap map[tview.Primitive]tview.Primitive) {
	for view, table := range tableMap {
		if t, ok := table.(widgets.Selectable); ok {
			t.SetSelectable(focusOrder[currentFocusIndex] == view, false)
		}
	}
}

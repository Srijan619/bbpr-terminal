package components

import (
	"github.com/rivo/tview"
)

// Create standard flex view for cohesion and less repeated code
func CreateFlexComponent(title string) *tview.Box {
	flex := tview.NewFlex().
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(title)

	return flex
}

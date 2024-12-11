package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"simple-git-terminal/state"
)

func CreateCheckBoxComponent(label string) *tview.Checkbox {
	checkedStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorGreen)

	uncheckedStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault)

	activatedStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorGreen)

	checkbox := tview.NewCheckbox().
		SetLabel(label).
		SetLabelColor(tcell.ColorDefault).
		SetCheckedString("✓").
		SetCheckedStyle(checkedStyle).
		SetUncheckedStyle(uncheckedStyle).
		SetActivatedStyle(activatedStyle)

	checkbox.SetChangedFunc(func(checked bool) {
		if checked {
			checkbox.SetLabelColor(tcell.ColorGreen)
		} else {
			checkbox.SetLabelColor(tcell.ColorDefault)
		}
		state.SetPRStatusFilter(label, checked)
	})
	return checkbox
}

// Create standard flex view for cohesion and less repeated code
func CreateFlexComponent(title string) *tview.Flex {
	flex := tview.NewFlex()

	flex.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(title)

	return flex
}

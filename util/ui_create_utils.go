package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateCheckBoxComponent(label string, onChange func(bool)) *tview.Checkbox {
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
		// Call the onChange callback to propagate the change
		if onChange != nil {
			onChange(checked)
		}
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

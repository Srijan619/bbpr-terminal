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

	checkbox.SetBackgroundColor(tcell.ColorDefault)
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
		SetTitle(title).
		SetBorderColor(tcell.ColorGrey).
		SetBackgroundColor(tcell.ColorDefault)

	return flex
}

func CreateTextviewComponent(title string, border bool) *tview.TextView {
	textView := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetWrap(true)

	textView.
		SetBorder(border).
		SetBorderColor(tcell.ColorGrey).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDefault).
		SetBorderPadding(1, 1, 1, 1)

	textView.SetTextStyle(tcell.StyleDefault.Background(tcell.ColorDefault))

	return textView
}

func CreateTextAreaComponent(title string, placeholder string) *tview.TextArea {
	textArea := tview.NewTextArea().
		SetPlaceholder(placeholder)

	textArea.SetTitle(title).SetTitleAlign(tview.AlignLeft)

	textArea.SetPlaceholderStyle(tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault))
	textArea.SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault))

	textArea.SetBorder(true).
		SetBorderColor(tcell.ColorGrey).
		SetBackgroundColor(tcell.ColorDefault)

	return textArea
}

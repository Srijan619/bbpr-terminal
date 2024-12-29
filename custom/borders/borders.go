package borders

import "github.com/rivo/tview"

// CustomizeBorders updates both regular and focus borders with rounded corners.
func CustomizeBorders() {
	tview.Borders.TopLeft = '╭'
	tview.Borders.TopRight = '╮'
	tview.Borders.BottomLeft = '╰'
	tview.Borders.BottomRight = '╯'

	tview.Borders.HorizontalFocus = '─'
	tview.Borders.VerticalFocus = '│'
	tview.Borders.TopLeftFocus = '╭'
	tview.Borders.TopRightFocus = '╮'
	tview.Borders.BottomLeftFocus = '╰'
	tview.Borders.BottomRightFocus = '╯'
}

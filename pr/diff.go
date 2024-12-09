package pr

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// Function to generate the diff output view
func GenerateDiffView(diffText string) *tview.TextView {
	// Initialize the TextView to display the diff
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetWrap(true).
		SetScrollable(true)

	// Split the diff text by lines and color them based on the prefix (+ or -)
	var coloredDiff []string
	lines := strings.Split(diffText, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+") {
			// Green for added lines
			coloredDiff = append(coloredDiff, fmt.Sprintf("[green]%s[-]", line))
		} else if strings.HasPrefix(line, "-") {
			// Red for removed lines
			coloredDiff = append(coloredDiff, fmt.Sprintf("[red]%si[-]", line))
		} else {
			// Normal lines without changes
			coloredDiff = append(coloredDiff, line)
		}
	}

	// Join the lines back together and set the text in the TextView
	textView.SetText(strings.Join(coloredDiff, "\n"))
	return textView
}

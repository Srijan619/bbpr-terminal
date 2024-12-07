package pr

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
	"strings"
)

func UpdatePrDetails(prs []types.PR, prDetails *tview.TextView, row int) {
	if row >= 0 && row < len(prs) {
		selectedPR := prs[row]

		// Format the description using glamour for Markdown rendering
		description := formatDescription(selectedPR.Description)

		// Get the color based on the state
		stateColor := util.GetStateColor(selectedPR.State)

		otherColor := tcell.ColorMediumAquamarine

		// Create a formatted string with improved structure and apply the state color
		formattedText := fmt.Sprintf(
			"[::b]State:[-] [%s]%s[-]\n"+
				"[::b]Author:[-] [%s]%s[-]\n"+
				"[::b]Created On:[-] [%s]%s[-]\n"+
				"[::b]Updated On:[-] [%s]%s[-]\n"+
				"[::b]Link:[-] [%s]%s[-]\n"+
				"[::b]Description:[-] \n%s\n",
			stateColor, selectedPR.State,
			otherColor, selectedPR.Author.DisplayName,
			otherColor, selectedPR.CreatedOn,
			otherColor, selectedPR.UpdatedOn,
			otherColor, selectedPR.Links.HTML.Href,
			description, // Rendered Markdown content
		)

		prDetails.SetText(formattedText)
	}
}

// Formats the PR description for display
func formatDescription(description interface{}) string {
	if description == nil {
		return "No description provided."
	}
	if desc, ok := description.(string); ok {
		trimmed := strings.TrimSpace(desc)
		return translateANSI(renderMarkdown(trimmed))
	}
	return "Unsupported description format."
}

// Renders the given Markdown string using the glamour library.
func renderMarkdown(md string) string {
	rendered, err := glamour.Render(md, "dark")
	if err != nil {
		log.Fatalf("Error rendering markdown: %v", err)
	}

	return rendered
}

// Translate ANSI escape sequences into tview-compatible format
func translateANSI(input string) string {
	var buf bytes.Buffer
	w := tview.ANSIWriter(&buf)
	_, err := w.Write([]byte(input))
	if err != nil {
		log.Fatalf("Error translating ANSI: %v", err)
	}
	return buf.String()
}

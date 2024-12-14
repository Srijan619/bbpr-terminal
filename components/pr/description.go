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

func GeneratePRDetail(pr *types.PR) string {
	// Format the description using glamour for Markdown rendering
	description := formatDescription(pr.Description)

	// Get the color based on the state
	stateColor := util.GetStateColor(pr.State)

	otherColor := tcell.ColorMediumAquamarine

	reviewers := GetReviewerNames(pr)

	// Create a formatted string with improved structure and apply the state color
	formattedText := fmt.Sprintf(
		"[::b]Reviewers:[-] [%s]%s[-]\n"+
			"[::b]State:[-] [%s]%s[-]\n"+
			"[::b]Author:[-] [%s]%s[-]\n"+
			"[::b]Created On:[-] [%s]%s[-]\n"+
			"[::b]Updated On:[-] [%s]%s[-]\n"+
			"[::b]Link:[-] [%s]%s[-]\n"+
			"[::b]Description:[-] \n%s\n",
		tcell.ColorOrange, reviewers,
		stateColor, pr.State,
		otherColor, pr.Author.DisplayName,
		otherColor, pr.CreatedOn,
		otherColor, pr.UpdatedOn,
		otherColor, pr.Links.HTML.Href,
		description, // Rendered Markdown content
	)

	return formattedText
}

func GetReviewerNames(pr *types.PR) string {
	var reviewerNames []string

	// Loop through the participants and check if they are REVIEWERs
	for _, participant := range pr.Participants {
		if participant.Role == "REVIEWER" {
			reviewerNames = append(reviewerNames, util.FormatInitials(participant.User.DisplayName))
		}
	}

	if len(reviewerNames) == 0 {
		return "No Reviewers"
	}
	// Join all the reviewer names with a pipe (" | ") separator
	return strings.Join(reviewerNames, " | ")
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

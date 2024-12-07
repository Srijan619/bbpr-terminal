package pr

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
)

func UpdatePrDetails(prs []types.PR, prDetails *tview.TextView, row int) {
	if row >= 0 && row < len(prs) {
		selectedPR := prs[row]
		description := formatDescription(selectedPR.Description)

		// Get the color based on the state
		stateColor := util.GetStateColor(selectedPR.State)

		otherColor := tcell.ColorMediumAquamarine
		// Create a formatted string with improved structure and apply the state color
		prDetails.SetText(fmt.Sprintf(
			"[::b]State:[-] [%s]%s[-]\n"+
				"[::b]Author:[-] [%s]%s[-]\n"+
				"[::b]Created On:[-] [%s]%s[-]\n"+
				"[::b]Updated On:[-] [%s]%s[-]\n"+
				"[::b]Link:[-] [%s]%s[-]\n"+
				"[::b]Description:[-] \n[%s]%s[-]\n",
			stateColor, selectedPR.State,
			otherColor, selectedPR.Author.DisplayName,
			otherColor, selectedPR.CreatedOn,
			otherColor, selectedPR.UpdatedOn,
			otherColor, selectedPR.Links.HTML.Href,
			otherColor, description,
		))
	}
}

// Formats the PR description for display
func formatDescription(description interface{}) string {
	if description == nil {
		return "No description provided."
	}
	if desc, ok := description.(string); ok {
		return strings.TrimSpace(desc)
	}
	return "Unsupported description format."
}

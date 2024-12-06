package pr

import (
	"fmt"
	"github.com/rivo/tview"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
	"strings"
)

func UpdatePrDetails(prs []types.PR, prDetails *tview.TextView, row int) {
	if row >= 0 && row < len(prs) {
		selectedPR := prs[row]
		description := formatDescription(selectedPR.Description)

		// Get the color based on the state
		stateColor := util.GetStateColor(selectedPR.State)

		// Create a formatted string with improved structure and apply the state color
		prDetails.SetText(fmt.Sprintf(
			"[::b]State:[-] [%s]%s[-]\n"+
				"[::b]Author:[-] [green]%s[-]\n"+
				"[::b]Created On:[-] [green]%s[-]\n"+
				"[::b]Updated On:[-] [green]%s[-]\n"+
				"[::b]Link:[-] [green]%s[-]\n"+
				"[::b]Description:[-] [green]%s[-]\n",
			stateColor, selectedPR.State,
			selectedPR.Author.DisplayName,
			selectedPR.CreatedOn,
			selectedPR.UpdatedOn,
			selectedPR.Links.HTML.Href,
			description,
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


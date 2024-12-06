package pr

import (
	"fmt"
	"strings"

	"time"

	"github.com/rivo/tview"

	"simple-git-terminal/types"
	"simple-git-terminal/util"
)

// CreateActivitiesView generates the UI for displaying PR activities in a TextView.
func CreateActivitiesView(activities []types.Activity) *tview.Flex {
	// Create a TextView for displaying activity details
	activityDetails := tview.NewTextView().
		SetDynamicColors(true).
		SetText(GenerateActivityLogs(activities)).
		SetWrap(true)

	// Layout to return
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(activityDetails, 0, 1, true)

	return layout
}

func isEmptyUpdateDetail(update types.UpdateDetail) bool {
	return update.State == "" && update.Title == "" && update.Description == "" && update.Date == "" && len(update.Changes) == 0
}

func GenerateActivityLogs(activities []types.Activity) string {
	var logs []string

	// Separate logs into sections
	updateLogs := []string{"[::b][cyan]Updates:[-]\n"}
	approvalLogs := []string{"[::b][cyan]Approvals:[-]\n"}
	prLogs := []string{"[::b][cyan]Pull Requests:[-]\n"}

	for i, activity := range activities {
		i = i + 1
		switch {
		case !isEmptyUpdateDetail(activity.Update):
			// Handle updates
			for field, change := range activity.Update.Changes {
				log := fmt.Sprintf(
					"[red]{%d}[-] %s edited the [%s]%s[-]: %s → %s [cyan](%s ago)[-]",
					i,
					activity.Update.Author.DisplayName,
					util.GetFieldBasedColor(field),
					field,
					change.Old,
					change.New[:min(len(change.New), 30)], // Truncate if too long
					formatTimeAgo(activity.Update.Date),
				)
				updateLogs = append(updateLogs, log)
			}
		case activity.Approval.User.DisplayName != "":
			// Handle approvals
			log := fmt.Sprintf(
				"[red]{%d}[-] %s [palegreen]APPROVED[-] the pull request [cyan](%s ago)[-]",
				i,
				activity.Approval.User.DisplayName,
				formatTimeAgo(activity.Approval.Date),
			)
			approvalLogs = append(approvalLogs, log)
		case activity.PullRequest.Title != "":
			// Handle pull requests
			log := fmt.Sprintf(
				"[red]{%d}[-] %s OPENED the pull request: %s [cyan](%s ago)[-]",
				i,
				activity.PullRequest.Author.DisplayName,
				activity.PullRequest.Title,
				formatTimeAgo(activity.PullRequest.CreatedOn),
			)
			prLogs = append(prLogs, log)
		}
	}

	// Add the logs and dividers only if there are actual entries in the section
	if len(updateLogs) > 1 { // Check if there are any updates
		logs = append(logs, strings.Join(updateLogs, "\n"))
		logs = append(logs, "[cyan]-------------------------------------------------[-]")
	}
	if len(approvalLogs) > 1 { // Check if there are any approvals
		logs = append(logs, strings.Join(approvalLogs, "\n"))
		logs = append(logs, "[cyan]-------------------------------------------------[-]")
	}
	if len(prLogs) > 1 { // Check if there are any pull requests
		logs = append(logs, strings.Join(prLogs, "\n"))
		logs = append(logs, "[cyan]-------------------------------------------------[-]")
	}

	// Join all the logs together
	return strings.Join(logs, "\n")
}

// Helper function to calculate time ago
func formatTimeAgo(date string) string {
	parsedTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "unknown time"
	}
	duration := time.Since(parsedTime)

	if hours := duration.Hours(); hours > 24 {
		return fmt.Sprintf("%d days", int(hours/24))
	} else if hours > 1 {
		return fmt.Sprintf("%d hours", int(hours))
	} else if minutes := duration.Minutes(); minutes > 1 {
		return fmt.Sprintf("%d minutes", int(minutes))
	}
	return "just now"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

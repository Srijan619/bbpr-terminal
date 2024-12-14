package pr

import (
	"fmt"
	"strings"
	"time"

	"github.com/rivo/tview"

	"simple-git-terminal/types"
	"simple-git-terminal/util"
)

const (
	ICON_UPDATES      = "\ue206 "
	ICON_APPROVAL     = "\u2713 "
	ICON_PULL_REQUEST = "\ue6a6 "
	ICON_EMPTY        = "\uf111 "
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
	updateLogs := []string{ICON_UPDATES + "[::b][darkslateblue]Updates:[-]\n"}
	approvalLogs := []string{ICON_APPROVAL + "[::b][darkslateblue]Approvals:[-]\n"}
	prLogs := []string{ICON_PULL_REQUEST + "[::b][darkslateblue]Pull Requests:[-]\n"}

	itemsCount := 0
	for _, activity := range activities {
		switch {
		case !isEmptyUpdateDetail(activity.Update):
			// Handle updates
			for field, change := range activity.Update.Changes {
				itemsCount++
				log := fmt.Sprintf(
					"[red]{%d}[-] %s edited the [%s]%s[-]: %s → %s [cyan](%s ago)[-]\n",
					itemsCount,
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
			itemsCount++
			log := fmt.Sprintf(
				"[red]{%d}[-] %s [green]APPROVED[-] the pull request [cyan](%s ago)[-]\n",
				itemsCount,
				activity.Approval.User.DisplayName,
				formatTimeAgo(activity.Approval.Date),
			)
			approvalLogs = append(approvalLogs, log)
		case activity.PullRequest.Title != "":
			// Handle pull requests
			itemsCount++
			log := fmt.Sprintf(
				"[red]{%d}[-] %s OPENED the pull request: %s [cyan](%s ago)[-]\n",
				itemsCount,
				activity.PullRequest.Author.DisplayName,
				activity.PullRequest.Title,
				formatTimeAgo(activity.PullRequest.CreatedOn),
			)
			prLogs = append(prLogs, log)
		}
	}

	// Check if there are no activities
	if itemsCount == 0 {
		return ICON_EMPTY + "[::b][red]No activities----![-]"
	}

	// Add the logs and dividers only if there are actual entries in the section
	if len(updateLogs) > 1 { // Check if there are any updates
		logs = append(logs, strings.Join(updateLogs, "\n"))
	}
	if len(approvalLogs) > 1 { // Check if there are any approvals
		logs = append(logs, strings.Join(approvalLogs, "\n"))
	}
	if len(prLogs) > 1 { // Check if there are any pull requests
		logs = append(logs, strings.Join(prLogs, "\n"))
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

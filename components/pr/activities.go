package pr

import (
	"fmt"
	"strings"

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
	// Check if state, title, description, date are empty
	if update.State == "" && update.Title == "" && update.Description == "" && update.Date == "" {
		// Check if there are no changes in reviewers, title, or description
		return len(update.Changes.Reviewers.Added) == 0 &&
			update.Changes.Description.New == "" &&
			update.Changes.Title.New == ""
	}
	return false
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
			// Handle updates for changes in reviewers, title, description, etc.
			if len(activity.Update.Changes.Reviewers.Added) > 0 {
				for _, reviewer := range activity.Update.Changes.Reviewers.Added {
					itemsCount++
					log := fmt.Sprintf(
						"[grey]%d %s[-] %s added [blue]reviewer[-]: %s [grey](%s ago)[-]\n",
						itemsCount,
						ICON_SIDE_ARROW,
						activity.Update.Author.DisplayName,
						reviewer.DisplayName,
						util.FormatTimeAgo(activity.Update.Date),
					)
					updateLogs = append(updateLogs, log)
				}
			}

			if activity.Update.Changes.Title.Old != "" && activity.Update.Changes.Title.New != "" {
				itemsCount++
				log := fmt.Sprintf(
					"[grey]%d %s[-] %s edited the [blue]title[-]: %s → %s [grey](%s ago)[-]\n",
					itemsCount,
					ICON_SIDE_ARROW,
					activity.Update.Author.DisplayName,
					activity.Update.Changes.Title.Old,
					activity.Update.Changes.Title.New,
					util.FormatTimeAgo(activity.Update.Date),
				)
				updateLogs = append(updateLogs, log)
			}

			if activity.Update.Changes.Description.Old != "" && activity.Update.Changes.Description.New != "" {
				itemsCount++
				log := fmt.Sprintf(
					"[grey]%d %s[-] %s edited the [blue]description[-]: %s → %s [grey](%s ago)[-]\n",
					itemsCount,
					ICON_SIDE_ARROW,
					activity.Update.Author.DisplayName,
					activity.Update.Changes.Description.Old,
					activity.Update.Changes.Description.New,
					util.FormatTimeAgo(activity.Update.Date),
				)
				updateLogs = append(updateLogs, log)
			}
			if activity.Update.Title != "" &&
				len(activity.Update.Changes.Reviewers.Added) == 0 &&
				activity.Update.Changes.Description.New == "" &&
				activity.Update.Changes.Title.New == "" {
				itemsCount++
				log := fmt.Sprintf(
					"[grey]%d %s[-] %s [blue]OPENED[-] the pull request: %s [grey](%s ago)[-]\n",
					itemsCount,
					ICON_SIDE_ARROW,
					activity.Update.Author.DisplayName,
					activity.Update.Title,
					util.FormatTimeAgo(activity.PullRequest.CreatedOn),
				)
				updateLogs = append(updateLogs, log)
			}
		case activity.Approval.User.DisplayName != "":
			// Handle approvals (if there's an approval activity)
			itemsCount++
			log := fmt.Sprintf(
				"[grey]%d %s[-] %s [green]APPROVED[-] the pull request [grey](%s ago)[-]\n",
				itemsCount,
				ICON_SIDE_ARROW,
				activity.Approval.User.DisplayName,
				util.FormatTimeAgo(activity.Approval.Date),
			)
			approvalLogs = append(approvalLogs, log)
		}
	}

	// Check if there are no activities
	if itemsCount == 0 {
		return ICON_EMPTY + "[::b][grey]No activities----![-]"
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

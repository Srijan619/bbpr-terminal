package pr

import (
	"fmt"
	"log"
	"sort"
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
	ICON_WARNING      = "\u2260 "
	ICON_COMMENT      = "\uf27b "
)

// CreateActivitiesView generates the UI for displaying PR activities in a TextView.
func CreateActivitiesView(activities []types.Activity) *tview.Flex {
	// Create a TextView for displaying activity details
	activityDetails := util.CreateTextviewComponent("").SetText(GenerateActivityLogs(activities))

	// Layout to return
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).AddItem(activityDetails, 0, 1, true)

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
	updateLogs := []string{ICON_UPDATES + "[red]Updates\n"}
	approvalLogs := []string{ICON_APPROVAL + "[red]Approvals\n"}
	commentLogs := []string{ICON_COMMENT + "[red]Comments\n"}

	itemsCount := 0
	log.Printf("Total activities..%d", len(activities))
	var previousCommitHash string // Track the previous commit hash
	var openPRFound bool

	// Sort activities by CreatedOn
	sort.SliceStable(activities, func(i, j int) bool {
		timeI, errI := time.Parse(time.RFC3339, activities[i].Update.Date)
		timeJ, errJ := time.Parse(time.RFC3339, activities[j].Update.Date)
		if errI != nil || errJ != nil {
			log.Println("Error parsing dates:", errI, errJ)
			return false
		}
		return timeI.Before(timeJ)
	})
	for _, activity := range activities {
		switch {
		case !isEmptyUpdateDetail(activity.Update):
			// Handle PR opening when no changes other than title
			if !openPRFound && activity.Update.Title != "" &&
				len(activity.Update.Changes.Reviewers.Added) == 0 &&
				activity.Update.Changes.Description.New == "" &&
				activity.Update.Changes.Title.New == "" &&
				activity.Update.Source.Commit.Hash == previousCommitHash {
				itemsCount++
				openPRFound = true
				previousCommitHash = activity.Update.Source.Commit.Hash
				log := fmt.Sprintf(
					"[mediumaquamarine] [-]%s [mediumaquamarine]opened[-] the pull request: %s [grey](%s)[-]\n",
					activity.Update.Author.DisplayName,
					activity.Update.Title,
					util.FormatTimeAgo(activity.Update.Date),
				)
				updateLogs = append(updateLogs, log)
			}

			// Check if the commit hash has changed
			if activity.Update.Source.Commit.Hash != previousCommitHash {
				itemsCount++
				previousCommitHash = activity.Update.Source.Commit.Hash
				log := fmt.Sprintf(
					"[orange]\ue729 [-]%s [orange]updated[-] the pull request with a new commit: [steelblue]%s[-] [grey](%s)[-]\n",
					activity.Update.Author.DisplayName,
					activity.Update.Source.Commit.Hash,
					util.FormatTimeAgo(activity.Update.Date),
				)
				updateLogs = append(updateLogs, log)
			}

			// Handle updates for changes in reviewers, title, description, etc.
			if len(activity.Update.Changes.Reviewers.Added) > 0 {
				itemsCount++
				for _, reviewer := range activity.Update.Changes.Reviewers.Added {
					log := fmt.Sprintf(
						"[purple]+ [-]%s added [purple]reviewer[-]: %s [grey](%s)[-]\n",
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
					"[blue]\uea73 [-]%s edited the [blue]title[-]: %s → %s [grey](%s)[-]\n",
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
					"[blue]\uea73 [-]%s edited the [blue]description[-]: %s → %s [grey](%s)[-]\n",
					activity.Update.Author.DisplayName,
					activity.Update.Changes.Description.Old,
					activity.Update.Changes.Description.New,
					util.FormatTimeAgo(activity.Update.Date),
				)
				updateLogs = append(updateLogs, log)
			}

		case activity.Approval.User.DisplayName != "":
			// Handle approvals (if there's an approval activity)
			itemsCount++
			log := fmt.Sprintf(
				"[limegreen] [-]%s [limegreen]APPROVED[-] the pull request [grey](%s)[-]\n",
				activity.Approval.User.DisplayName,
				util.FormatTimeAgo(activity.Approval.Date),
			)
			approvalLogs = append(approvalLogs, log)

		case activity.ChangesRequested.Date != "":
			// Handle Changes requested
			itemsCount++
			log := fmt.Sprintf(
				"[yellow]%s[-]%s [yellow]requested changes[-] [grey](%s)[-]\n",
				ICON_WARNING,
				activity.ChangesRequested.User.DisplayName,
				util.FormatTimeAgo(activity.ChangesRequested.Date),
			)
			updateLogs = append(updateLogs, log)
		case activity.Comment.ID > 0:
			// Handle Comments added
			itemsCount++
			log := fmt.Sprintf(
				"[yellowgreen]%s[-]%s [yellowgreen]added a comment: [-] %s [grey](%s)[-]\n",
				ICON_COMMENT,
				activity.Comment.User.DisplayName,
				activity.Comment.Content.Raw,
				util.FormatTimeAgo(activity.Comment.CreatedOn),
			)
			commentLogs = append(commentLogs, log)
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
	if len(commentLogs) > 1 {
		logs = append(logs, strings.Join(commentLogs, "\n"))
	}

	// Join all the logs together
	return strings.Join(logs, "\n")
}

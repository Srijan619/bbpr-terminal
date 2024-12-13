package pr

import (
	"fmt"
	"log"

	"github.com/rivo/tview"

	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
)

const (
	ICON_LOADING = "\uea75 "
)

func PopulatePRList(prList *tview.Table) *tview.Table {
	prs := GetFilteredPRs()
	log.Printf("PRS....%v", prs)
	// Populate PR list
	util.PopulatePRList(prList, prs)

	// Add a selection function that updates PR details when a PR is selected
	prList.SetSelectedFunc(func(row, column int) {
		//	HandleOnPrSelect(prs, row)
	})

	// Set initial PR details if available
	if len(prs) > 0 {
		prList.Select(0, 0)
		//HandleOnPrSelect(prs, 0)
	}
	return prList
}

func HandleOnPrSelect(prs []types.PR, row int) {
	if row >= 0 && row < len(prs) {
		// Update PR details and set the selected PR
		UpdatePrDetails(prs, state.GlobalState.PrDetails, row)
		state.SetSelectedPR(&prs[row])

		// Update right panel and set header
		state.GlobalState.RightPanelHeader.SetText(formatPRHeader(*state.GlobalState.SelectedPR))

		// Fetch details in parallel using goroutines
		go func() {
			// Show loading spinner for activities
			util.ShowLoadingSpinner(state.GlobalState.ActivityView, func() (string, error) {
				// Fetch activities
				prActivities := bitbucket.FetchBitbucketActivities(state.GlobalState.SelectedPR.ID)
				if prActivities == nil {
					return "", fmt.Errorf("Failed to fetch activities")
				}
				return ICON_LOADING + "Activities fetched!", nil
			}, func(result string, err error) {
				if err != nil {
					util.UpdateActivityView(err.Error())
				} else {
					util.UpdateActivityView(CreateActivitiesView(bitbucket.FetchBitbucketActivities(state.GlobalState.SelectedPR.ID)))
				}
			})

			// Show loading spinner for diff stats
			util.ShowLoadingSpinner(state.GlobalState.DiffStatView, func() (string, error) {
				// Fetch diff stats
				diffStatData := bitbucket.FetchBitbucketDiffstat(state.GlobalState.SelectedPR.ID)
				if diffStatData == nil {
					return "", fmt.Errorf("Failed to fetch diff stats")
				}
				return ICON_LOADING + "Diff stats fetched!", nil
			}, func(result string, err error) {
				if err != nil {
					util.UpdateDiffStatView(err.Error())
				} else {
					util.UpdateDiffStatView(GenerateDiffStatTree(bitbucket.FetchBitbucketDiffstat(state.GlobalState.SelectedPR.ID)))
				}
			})

			// Optionally, set a default message in the diff details view while fetching
			util.UpdateDiffDetailsView("Hover over to a file for quick preview OR Select a file to see diff in full screen")
		}()
	}
}

// Function to populate the PR list

// FormatPRHeader takes the PR details and returns a formatted string
func formatPRHeader(pr types.PR) string {
	// Use fmt.Sprintf to format the header and apply tview's dynamic color syntax
	headerText := fmt.Sprintf(
		"%s\n\n"+
			"[yellow]%s[white] -> "+
			"[green]%s[white]",
		pr.Title,
		pr.Source.Branch.Name,
		pr.Destination.Branch.Name,
	)

	return headerText
}

func GetFilteredPRs() []types.PR {
	var filteredPRs []types.PR

	// Fetch or use cached PRs based on active filters
	if state.PRStatusFilter.Open {
		log.Printf("I am fetching open PRs... now..")
		filteredPRs = append(filteredPRs, bitbucket.FetchPRsByState("OPEN")...)
	}
	if state.PRStatusFilter.Merged {
		log.Printf("I am fetching merged PRs... now..")
		filteredPRs = append(filteredPRs, bitbucket.FetchPRsByState("MERGED")...)
	}
	if state.PRStatusFilter.Declined {
		log.Printf("I am fetching declined PRs... now..")
		filteredPRs = append(filteredPRs, bitbucket.FetchPRsByState("DECLINED")...)
	}

	log.Printf("Length of filtered PRS..%d", len(filteredPRs))
	return filteredPRs
}

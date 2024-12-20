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
	ICON_LOADING    = "\uea75 "
	ICON_SIDE_ARROW = "\u21AA "
)

func PopulatePRList(prList *tview.Table) *tview.Table {
	bitbucket.UpdateFilteredPRs()
	prs := *state.GlobalState.FilteredPRs
	if len(prs) > 0 {
		prList.Select(0, 0)
		HandleOnPrSelect(prs, 0)
	}
	// Populate PR list
	util.PopulatePRList(prList, prs)

	// Add a selection function that updates PR details when a PR is selected
	prList.SetSelectedFunc(func(row, column int) {
		go func() {
			prs := *state.GlobalState.FilteredPRs // use updated prs inside routine
			HandleOnPrSelect(prs, row)
			UpdatePRList()
		}()
	})

	prList.SetSelectionChangedFunc(func(row, column int) {
		//HandleOnPrSelect(prs, row) TODO: It will be really laggy to allow selection update PR as user might navigate up and down fast...maybe need some debouncing?
	})

	return prList
}

func HandleOnPrSelect(prs []types.PR, row int) {
	if state.GlobalState != nil {
		fetchMore := row == len(prs)
		if fetchMore {
			log.Println("Fetch more selected")
		} else {
			// Handle normal cell selection
			log.Printf("Selected cell at row %d", row)
			handleNormalPRSelect(prs, row)
		}

	}
}

// Function to populate the PR list

// FormatPRHeader takes the PR details and returns a formatted string
func formatPRHeaderBranch(pr types.PR) string {
	// Use fmt.Sprintf to format the header and apply tview's dynamic color syntax
	headerText := fmt.Sprintf(
		"[yellow]%s[white] "+ICON_SIDE_ARROW+
			"[green]%s",
		pr.Source.Branch.Name,
		pr.Destination.Branch.Name,
	)

	return headerText
}

func handleNormalPRSelect(prs []types.PR, row int) {
	if row >= 0 && row < len(prs) && state.GlobalState != nil {
		// Fetch details in parallel using goroutines
		go func() {
			state.SetSelectedPR(&prs[row])

			//	 Update right panel and set header
			state.GlobalState.RightPanelHeader.SetTitle(formatPRHeaderBranch(*state.GlobalState.SelectedPR))
			state.GlobalState.RightPanelHeader.SetText(state.GlobalState.SelectedPR.Title)

			// Show loading spinner for PR details
			util.ShowLoadingSpinner(state.GlobalState.PrDetails, func() (interface{}, error) {
				singlePR := bitbucket.FetchPR(prs[row].ID)
				if singlePR == nil {
					return nil, fmt.Errorf("Failed to fetch PR details")
				}
				return singlePR, nil
			}, func(result interface{}, err error) {
				if err != nil {
					util.UpdatePRDetailView(fmt.Sprintf("[red]Error: %v[-]", err))
				} else {
					// Assert result as the correct type: *types.PR
					pr, ok := result.(*types.PR)
					if !ok {
						util.UpdatePRDetailView("[red]Failed to cast PR details[-]")
						return
					}
					util.UpdatePRDetailView(GeneratePRDetail(pr))
				}
			})

			// Show loading spinner for activities
			util.ShowLoadingSpinner(state.GlobalState.ActivityView, func() (interface{}, error) {
				// Fetch activities
				prActivities := bitbucket.FetchBitbucketActivities(state.GlobalState.SelectedPR.ID)
				if prActivities == nil {
					return nil, fmt.Errorf("Failed to fetch activities")
				}
				return prActivities, nil
			}, func(result interface{}, err error) {
				if err != nil {
					util.UpdateActivityView(err.Error())
				} else {
					// Assert result as a slice of Activity
					activities, ok := result.([]types.Activity)
					if !ok {
						util.UpdateActivityView("[red]Failed to cast activities[-]")
						return
					}
					util.UpdateActivityView(CreateActivitiesView(activities))
				}
			})

			// Show loading spinner for diff stats
			util.ShowLoadingSpinner(state.GlobalState.DiffStatView, func() (interface{}, error) {
				// Fetch diff stats
				diffStatData := bitbucket.FetchBitbucketDiffstat(state.GlobalState.SelectedPR.ID)
				if diffStatData == nil {
					return nil, fmt.Errorf("Failed to fetch diff stats")
				}
				return diffStatData, nil
			}, func(result interface{}, err error) {
				if err != nil {
					util.UpdateDiffStatView(err.Error())
				} else {
					// Assert result as string
					diffStat, ok := result.([]types.DiffstatEntry)
					if !ok {
						util.UpdateDiffStatView("[red]Failed to cast diff stats[-]")
						return
					}
					util.UpdateDiffStatView(GenerateDiffStatTree(diffStat))
				}
			})

			// Optionally, set a default message in the diff details view while fetching
			util.UpdateDiffDetailsView("Hover over to a file for quick preview OR Select a file to see diff in full screen")
		}()
	}
}

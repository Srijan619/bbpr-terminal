package pr

import (
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/components/shared"
	"simple-git-terminal/constants"
	"simple-git-terminal/state"
	"simple-git-terminal/support"
	"simple-git-terminal/types"

	"github.com/rivo/tview"
)

func PopulatePRList(prList *tview.Table) *tview.Table {
	UpdateFilteredPRs()
	prs := *state.GlobalState.FilteredPRs
	if len(prs) > 0 {
		prList.Select(0, 0)
		HandleOnPrSelect(prs, 0)
	}
	// Populate PR list
	//	PopulatePRList(prList, prs)

	// Populate pagination
	//
	// Add pagination below the PR list
	log.Printf("Current pagination state..%+v", state.Pagination)

	pagination := shared.NewPaginationComponent(state.Pagination.Page)
	// Initial render (first-time view)
	support.UpdateView(state.GlobalState.PaginationFlex, pagination)

	// Add a selection function that updates PR details when a PR is selected
	prList.SetSelectedFunc(func(row, column int) {
		go func() {
			prs := *state.GlobalState.FilteredPRs // use updated prs inside routine
			HandleOnPrSelect(prs, row)
			ShowSpinnerFetchPRsByQueryAndUpdatePrList()
		}()
	})

	prList.SetSelectionChangedFunc(func(row, column int) {
		// HandleOnPrSelect(prs, row) TODO: It will be really laggy to allow selection update PR as user might navigate up and down fast...maybe need some debouncing?
	})

	return prList
}

func HandleOnPrSelect(prs []types.PR, row int) {
	if state.GlobalState != nil {
		fetchMore := row == len(prs)
		if fetchMore {
			log.Printf("Fetch more...")
		} else {
			handleNormalPRSelect(prs, row)
		}

	}
}

// Function to populate the PR list

// FormatPRHeader takes the PR details and returns a formatted string
func formatPRHeaderBranch(pr types.PR) string {
	// Use fmt.Sprintf to format the header and apply tview's dynamic color syntax
	headerText := fmt.Sprintf(
		"[yellow]%s[white] "+constants.ICON_SIDE_ARROW+
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
			support.ShowLoadingSpinner(state.GlobalState.PrDetails, func() (interface{}, error) {
				singlePR := bitbucket.FetchPR(prs[row].ID)
				if singlePR == nil {
					return nil, fmt.Errorf("Failed to fetch PR details")
				}
				return singlePR, nil
			}, func(result interface{}, err error) {
				if err != nil {
					UpdatePRDetailView(fmt.Sprintf("[red]Error: %v[-]", err))
				} else {
					// Assert result as the correct type: *types.PR
					pr, ok := result.(*types.PR)
					if !ok {
						UpdatePRDetailView("[red]Failed to cast PR details[-]")
						return
					}
					UpdatePRDetailView(GeneratePRDetail(pr))
				}
			})

			// Show loading spinner for activities
			support.ShowLoadingSpinner(state.GlobalState.ActivityView, func() (interface{}, error) {
				// Fetch activities
				prActivities := bitbucket.FetchBitbucketActivities(state.GlobalState.SelectedPR.ID)
				if prActivities == nil {
					return nil, fmt.Errorf("Failed to fetch activities")
				}
				return prActivities, nil
			}, func(result interface{}, err error) {
				if err != nil {
					UpdateActivityView(err.Error())
				} else {
					// Assert result as a slice of Activity
					activities, ok := result.([]types.Activity)
					if !ok {
						UpdateActivityView("[red]Failed to cast activities[-]")
						return
					}
					UpdateActivityView(CreateActivitiesView(activities))
				}
			})

			// Show loading spinner for diff stats
			support.ShowLoadingSpinner(state.GlobalState.DiffStatView, func() (interface{}, error) {
				// Fetch diff stats
				diffStatData := bitbucket.FetchBitbucketDiffstat(state.GlobalState.SelectedPR.ID)
				if diffStatData == nil {
					return nil, fmt.Errorf("Failed to fetch diff stats")
				}
				return diffStatData, nil
			}, func(result interface{}, err error) {
				if err != nil {
					UpdateDiffStatView(err.Error())
				} else {
					// Assert result as string
					diffStat, ok := result.([]types.DiffstatEntry)
					if !ok {
						UpdateDiffStatView("[red]Failed to cast diff stats[-]")
						return
					}
					UpdateDiffStatView(GenerateDiffStatTree(diffStat))
				}
			})

			// Optionally, set a default message in the diff details view while fetching
			UpdateDiffDetailsView("Hover over to a file for quick preview OR Select a file to see diff in full screen")
		}()
	}
}

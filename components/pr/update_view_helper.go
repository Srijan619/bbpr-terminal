package pr

import (
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/support"
	"simple-git-terminal/types"
)

func UpdateActivityView(activityContent interface{}) {
	support.UpdateView(state.GlobalState.ActivityView, activityContent)
}

func UpdateDiffDetailsView(diffContent interface{}) {
	support.UpdateView(state.GlobalState.DiffDetails, diffContent)
}

func UpdateDiffStatView(statContent interface{}) {
	support.UpdateView(state.GlobalState.DiffStatView, statContent)
}

func UpdatePRListView() {
	if state.GlobalState != nil && state.GlobalState.PrList != nil && state.GlobalState.FilteredPRs != nil {
		state.GlobalState.PrList.Clear()
		PopulatePRList(state.GlobalState.PrList)
		//	state.GlobalState.App.Draw()
	}
}

func UpdatePRListErrorView() {
	if state.GlobalState != nil && state.GlobalState.PrList != nil && state.GlobalState.FilteredPRs != nil {
		state.GlobalState.PrList.Clear()
		support.UpdateView(state.GlobalState.PrList, "[red] Error rendering PR list")
		state.GlobalState.App.Draw()
	}
}

func UpdatePRDetailView(content interface{}) {
	support.UpdateView(state.GlobalState.PrDetails, content)
}

func UpdatePRStatusFilterView(content interface{}) {
	if state.GlobalState != nil && state.GlobalState.PRStatusFilter != nil {
		support.UpdateView(state.GlobalState.PRStatusFilter, content)
	}
}

func UpdatePRListWithFilter(filter string, checked bool) {
	state.SetPRStatusFilter(filter, checked)
	ShowSpinnerFetchPRsByQueryAndUpdatePrList()
}

func UpdateFilteredPRs() {
	prs, pagination := bitbucket.FetchPRsByQuery(bitbucket.BuildQuery(""))
	state.SetFilteredPRs(&prs)
	state.SetPagination(&pagination)
}

func ShowSpinnerFetchPRsByQueryAndUpdatePrList() {
	if state.GlobalState != nil {
		state.GlobalState.PrList.Clear()
		support.ShowLoadingSpinner(state.GlobalState.PrList, func() (interface{}, error) {
			prs, pagination := bitbucket.FetchPRsByQuery(bitbucket.BuildQuery(state.SearchTerm))
			return struct {
				PRs        []types.PR
				Pagination types.Pagination
			}{prs, pagination}, nil
		}, func(result interface{}, err error) {
			if err != nil {
				UpdatePRListErrorView()
			} else {
				data := result.(struct {
					PRs        []types.PR
					Pagination types.Pagination
				})

				state.SetFilteredPRs(&data.PRs)
				state.SetPagination(&data.Pagination)
				UpdatePRListView()
				// shared.UpdatePaginationViewUI(state.Pagination.Page)
			}
		})
	}
}

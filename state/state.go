package state

import (
	"github.com/rivo/tview"
	"log"
	"strings"

	"simple-git-terminal/types"
)

type State struct {
	App              *tview.Application
	MainFlexWrapper  *tview.Flex
	PrListFlex       *tview.Flex // Need both flex and table for styling and data
	PrList           *tview.Table
	PrDetails        *tview.TextView
	ActivityView     *tview.Flex
	DiffDetails      *tview.Flex
	DiffStatView     *tview.Flex
	RightPanelHeader *tview.TextView
	CurrentView      tview.Primitive
	PRStatusFilter   *tview.Flex
	PrListSearchBar  *tview.InputField
	PaginationFlex   *tview.Flex

	SelectedPR  *types.PR
	FilteredPRs *[]types.PR
}

var GlobalState *State
var Workspace, Repo string
var IsSearchMode bool
var SearchTerm string
var CurrentUser *types.User
var Pagination *types.Pagination = &types.Pagination{
	Next:    "",
	Page:    1,
	Size:    25,
	PageLen: 0,
}

// InitializeViews initializes all view components except workspace and repo.
func InitializeViews(app *tview.Application, mainFlexWrapper, prListFlex *tview.Flex, prList *tview.Table, prDetails *tview.TextView, activityView, diffDetails, diffStatView, pRStatusFilter *tview.Flex,
	rightPanelHeader *tview.TextView, prListSearchBar *tview.InputField, paginationFlex *tview.Flex) {
	GlobalState = &State{
		App:              app,
		MainFlexWrapper:  mainFlexWrapper,
		PrListFlex:       prListFlex,
		PrList:           prList,
		PrDetails:        prDetails,
		ActivityView:     activityView,
		DiffDetails:      diffDetails,
		DiffStatView:     diffStatView,
		PRStatusFilter:   pRStatusFilter,
		RightPanelHeader: rightPanelHeader,
		PaginationFlex:   paginationFlex,

		PrListSearchBar: prListSearchBar,
	}
}

// SetWorkspaceRepo sets the workspace and repo separately from the main state
func SetWorkspaceRepo(workspace, repo string) {
	Workspace = workspace
	Repo = repo
}

func SetCurrentView(currentView tview.Primitive) {
	GlobalState.CurrentView = currentView
}

// SetSelectedPR sets the selected PR in the global state.
func SetSelectedPR(pr *types.PR) {
	GlobalState.SelectedPR = pr
}

func SetFilteredPRs(prs *[]types.PR) {
	GlobalState.FilteredPRs = prs
}

func SetCurrentUser(user *types.User) {
	CurrentUser = user
}

func SetIsSearchMode(mode bool) {
	IsSearchMode = mode
}

func SetSearchTerm(term string) {
	SearchTerm = term
}

func SetPagination(pagination *types.Pagination) {
	Pagination = pagination
}

type PRStatusFilterType struct {
	Open        bool
	Merged      bool
	Declined    bool
	IAmAuthor   bool
	IAmReviewer bool
}

var PRStatusFilter *PRStatusFilterType

func InitializePRStatusFilter(filter *PRStatusFilterType) {
	if filter == nil {
		filter = &PRStatusFilterType{Open: true, Merged: false, Declined: false, IAmAuthor: false, IAmReviewer: true}
	}
	PRStatusFilter = filter
}

// Provide key (in any format) and whether that key is checked or not ("open | merged | declined | all", true | false)
func SetPRStatusFilter(filterKey string, isChecked bool) {
	trimmedFilterKey := strings.ToLower(strings.TrimSpace(filterKey))
	switch trimmedFilterKey {
	case "open":
		PRStatusFilter.Open = isChecked
	case "merged":
		PRStatusFilter.Merged = isChecked
	case "declined":
		PRStatusFilter.Declined = isChecked
	case "iamauthor":
		PRStatusFilter.IAmAuthor = isChecked
	case "iamreviewer":
		PRStatusFilter.IAmReviewer = isChecked
	case "all":
		PRStatusFilter.Open = isChecked
		PRStatusFilter.Merged = isChecked
		PRStatusFilter.Declined = isChecked
		PRStatusFilter.IAmAuthor = isChecked
		PRStatusFilter.IAmReviewer = isChecked

	}
	log.Printf("Filter updated: %+v\n", PRStatusFilter)
}

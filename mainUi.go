package main

import (
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/components/pr"
	"simple-git-terminal/custom/borders"
	"simple-git-terminal/state"
	"simple-git-terminal/support"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateMainApp() *tview.Application {
	borders.CustomizeBorders()
	app := tview.NewApplication()
	workspace, repoSlug, _ = util.GetRepoAndWorkspace()
	log.Printf("Loading workspace - %s and repo - %s ....", workspace, repoSlug)
	fmt.Printf("Loading workspace - %s and repo - %s ....", workspace, repoSlug)

	if (workspace == "") || (repoSlug == "") {
		log.Fatalf("Not a bitbucket Workspace")
		fmt.Printf("Not a bitbucket Workspace")
	}
	currentUser := bitbucket.FetchCurrentUser()
	state.SetCurrentUser(currentUser)

	state.SetWorkspaceRepo(workspace, repoSlug)
	util.InitMdRenderer() // Markdown renderer takes time, so init it beforehand

	// LEFT

	// PR Status Filter UI
	prStatusFilterFlex := support.CreateFlexComponent("Filters")
	prStatusFilterFlex.AddItem(pr.CreatePRStatusFilterView(), 0, 1, false)

	// PR LIST UI
	prListFlex := support.CreateFlexComponent("Pull Requests   [green]p|P").
		SetDirection(tview.FlexRow)

	prList := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)

	paginationFlex := support.CreateFlexComponent("")

	prList.SetBackgroundColor(tcell.ColorDefault)

	prListSearchBar := support.CreateInputFieldComponent("  Search PR [green]s", " type something....")

	prListFlex.
		AddItem(prList, 0, 1, true)

	leftFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	leftFullFlex.
		AddItem(prStatusFilterFlex, 0, 2, false).
		AddItem(prListSearchBar, 4, 1, false).
		AddItem(prListFlex, 0, 15, true).
		AddItem(paginationFlex, 0, 1, false)

		// Description and Activity

	activityDetails := support.CreateFlexComponent("Activities [green]a|A")

	// MIDDLE
	rightPanelHeader := support.CreateTextviewComponent("", true)
	prDetails := support.CreateTextviewComponent("Description [green]d|D", true)

	middleFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	middleFullFlex.SetBackgroundColor(tcell.ColorDefault)

	middleFullFlex.AddItem(rightPanelHeader, 0, 2, false).
		AddItem(prDetails, 0, 4, false).
		AddItem(activityDetails, 0, 14, false)

		// RIGHT

	diffStatDetails := support.CreateFlexComponent("Diff Tree [green]t|T")
	diffDetails := support.CreateFlexComponent("Diff Content [green]c|C")

	rightFullFlex := tview.NewFlex()

	rightFullFlex.SetBackgroundColor(tcell.ColorDefault)
	rightFullFlex.AddItem(diffStatDetails, 0, 1, false).
		AddItem(diffDetails, 0, 1, false)

	mainFlexWrapper := tview.NewFlex()
	mainFlexWrapper.SetBackgroundColor(tcell.ColorDefault)
	mainFlexWrapper.AddItem(leftFullFlex, 0, 1, true).
		AddItem(middleFullFlex, 0, 1, false).
		AddItem(rightFullFlex, 0, 2, false)

	state.InitializeViews(app, mainFlexWrapper, prListFlex, prList, prDetails, activityDetails, diffDetails, diffStatDetails, prStatusFilterFlex, rightPanelHeader, prListSearchBar, paginationFlex)
	pr.PopulatePRList(prList)

	// Key Bindings
	pr.SetupKeyBindings(func() {
		updateFilter() // TODO: We can do this better in organizing
	})

	app.SetRoot(mainFlexWrapper, true).EnableMouse(true)

	return app
}

func updateFilter() {
	view := pr.CreatePRStatusFilterView()
	pr.UpdatePRStatusFilterView(view)
}

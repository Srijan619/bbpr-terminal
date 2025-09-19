package main

import (
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/components/pipeline"
	"simple-git-terminal/components/pr"
	"simple-git-terminal/custom/borders"
	"simple-git-terminal/state"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateMainAppForBBPipeline() *tview.Application {
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

	// Pipeline Status Filter UI
	ppStatusFilterFlex := util.CreateFlexComponent("Filters")
	ppStatusFilterFlex.AddItem(pr.CreatePRStatusFilterView(), 0, 1, false)

	// Pipelines LIST UI
	ppListFlex := util.CreateFlexComponent("Pipelines   [green]p|P").
		SetDirection(tview.FlexRow)

	ppList := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)

	ppListFlex.SetBackgroundColor(tcell.ColorDefault)

	ppListFlex.
		AddItem(ppList, 0, 1, true)

	leftFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	leftFullFlex.
		AddItem(ppStatusFilterFlex, 0, 2, false).
		AddItem(ppListFlex, 0, 15, true)

		// MIDDLE
	ppDetails := util.CreateTextviewComponent("Details [green]d|D", true)

	middleFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	middleFullFlex.SetBackgroundColor(tcell.ColorDefault)

	middleFullFlex.AddItem(ppDetails, 0, 4, false)

	mainFlexWrapper := tview.NewFlex()
	mainFlexWrapper.SetBackgroundColor(tcell.ColorDefault)
	mainFlexWrapper.AddItem(leftFullFlex, 0, 1, true).
		AddItem(middleFullFlex, 0, 3, false)

	state.InitializePipelineViews(app, mainFlexWrapper, ppListFlex, ppList, ppDetails, nil, nil, nil, nil)

	pipeline.PopulatePipelineList(ppList)

	app.SetRoot(mainFlexWrapper, true).EnableMouse(true)

	return app
}

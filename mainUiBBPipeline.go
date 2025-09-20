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
	ppListFlex := util.CreateFlexComponent("Pipelines [green]p|P").
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
	stepsWrapper := tview.NewFlex()
	stepsWrapper.SetDirection(tview.FlexRow).SetBackgroundColor(tcell.ColorDefault)

	debugView := util.CreateFlexComponent("Debug Info")
	steps := tview.NewFlex()
	steps.SetBackgroundColor(tcell.ColorDefault)

	stepsWrapper.
		AddItem(debugView, 0, 1, true).
		AddItem(steps, 0, 3, true)

		// Right
	stepWrapper := tview.NewFlex()
	stepWrapper.SetDirection(tview.FlexRow).SetBackgroundColor(tcell.ColorDefault)

	step := util.CreateFlexComponent("Individual Step")
	stepCommandLogView := util.CreateFlexComponent("Step Command")

	stepWrapper.
		AddItem(step, 0, 1, true).
		AddItem(stepCommandLogView, 0, 2, true)

	middleFullFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	middleFullFlex.SetBackgroundColor(tcell.ColorDefault)

	middleFullFlex.AddItem(stepsWrapper, 0, 2, false)
	middleFullFlex.AddItem(stepWrapper, 0, 4, false)

	mainFlexWrapper := tview.NewFlex()
	mainFlexWrapper.SetBackgroundColor(tcell.ColorDefault)
	mainFlexWrapper.AddItem(leftFullFlex, 0, 1, true).
		AddItem(middleFullFlex, 0, 3, false)

	state.InitializePipelineViews(app, mainFlexWrapper, ppListFlex, ppList, debugView, steps, step, stepCommandLogView, nil, nil, nil)
	pipeline.PopulatePipelineList(ppList)

	pipeline.SetupKeyBindings()
	app.SetRoot(mainFlexWrapper, true).EnableMouse(true)

	return app
}

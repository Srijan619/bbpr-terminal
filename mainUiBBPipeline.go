package main

import (
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/components/pipeline"
	"simple-git-terminal/components/pr"
	"simple-git-terminal/custom/borders"
	"simple-git-terminal/events"
	"simple-git-terminal/state"
	"simple-git-terminal/support"
	"simple-git-terminal/util"
	"simple-git-terminal/views"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateMainAppForBBPipeline() *tview.Application {
	borders.CustomizeBorders()
	app := tview.NewApplication()
	bus := events.NewBus()

	workspace, repoSlug, _ = util.GetRepoAndWorkspace()
	log.Printf("Loading workspace - %s and repo - %s ....", workspace, repoSlug)
	fmt.Printf("Loading workspace - %s and repo - %s ....", workspace, repoSlug)

	if (workspace == "") || (repoSlug == "") {
		log.Fatalf("Not a bitbucket Workspace")
		fmt.Printf("Not a bitbucket Workspace")
	}
	currentUser := bitbucket.FetchCurrentUser()
	state.SetCurrentUser(currentUser)

	state.InitPartialPipelineState(app, workspace, repoSlug)
	util.InitMdRenderer() // Markdown renderer takes time, so init it beforehand

	// LEFT

	// Pipeline Status Filter UI
	ppStatusFilterFlex := support.CreateFlexComponent("Filters")
	ppStatusFilterFlex.AddItem(pr.CreatePRStatusFilterView(), 0, 1, false)

	// Pipelines LIST UI

	ppList := views.NewPipelineView(bus)

	leftFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	leftFullFlex.
		AddItem(ppStatusFilterFlex, 0, 2, false).
		AddItem(ppList.GetView(), 0, 15, true)

		// MIDDLE
	stepsWrapper := tview.NewFlex()
	stepsWrapper.SetDirection(tview.FlexRow).SetBackgroundColor(tcell.ColorDefault)

	debugView := support.CreateFlexComponent("Debug Info")
	steps := views.NewStepsView(bus)

	stepCommandsView := support.CreateFlexComponent("Script Commands")
	stepsWrapper.
		AddItem(debugView, 0, 2, true).
		AddItem(steps.Render(), 0, 2, true).
		AddItem(stepCommandsView, 0, 4, true)

		// Right
	stepWrapper := tview.NewFlex()
	stepWrapper.SetDirection(tview.FlexRow).SetBackgroundColor(tcell.ColorDefault)

	step := support.CreateFlexComponent("Individual Step")
	stepCommandLogView := support.CreateFlexComponent("Step Command")

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

	state.InitializePipelineViews(mainFlexWrapper, ppList.GetView(), debugView, steps.Render(), step, stepCommandsView, stepCommandLogView, nil, nil, nil)
	// pipeline.PopulatePipelineList()

	pipeline.SetupKeyBindings()
	app.SetRoot(mainFlexWrapper, true).EnableMouse(true)

	return app
}

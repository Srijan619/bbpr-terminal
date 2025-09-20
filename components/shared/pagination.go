package shared

import (
	"fmt"
	"simple-git-terminal/state"
	"simple-git-terminal/support"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewPaginationComponent(currentPage int) *tview.Flex {
	if state.Pagination == nil {
		return nil
	}

	totalItems := state.Pagination.Size
	itemsPerPage := state.Pagination.PageLen
	maxButtons := 5

	// Ensure itemsPerPage is not zero to prevent division by zero
	if itemsPerPage <= 0 {
		itemsPerPage = 1 // Default to 1 to ensure safe calculation
	}

	// Calculate total pages automatically
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage // Ceil(totalItems / itemsPerPage)

	// Ensure currentPage is within valid bounds
	if currentPage < 0 {
		currentPage = 0
	} else if currentPage > totalPages {
		currentPage = totalPages - 1
	}

	buttonFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Determine visible button range
	startPage := max(0, currentPage-maxButtons/2)
	endPage := min(startPage+maxButtons, totalPages)

	// Adjust range for exact `maxButtons` visibility
	if endPage-startPage < maxButtons {
		startPage = max(0, endPage-maxButtons)
	}

	// Add "First" button to go to the first page
	if currentPage > 1 {
		firstButton := tview.NewButton("<<")
		firstButton.SetLabelColor(tcell.ColorGrey).
			SetStyle(tcell.StyleDefault.Background(tcell.ColorDefault)).
			SetSelectedFunc(func() {
				// Update the page in the global state
				currentPage = 1
				// Fetch data or refresh view if needed
				updatePaginationView(currentPage)
			})
		buttonFlex.AddItem(firstButton, 0, 1, false)
	}

	// Add "Previous" button to go to the previous page
	if currentPage > 1 {
		prevButton := tview.NewButton("<").SetLabelColor(tcell.ColorGrey).
			SetStyle(tcell.StyleDefault.Background(tcell.ColorDefault)).
			SetSelectedFunc(func() {
				// Update the page in the global state
				currentPage -= 1
				// Fetch data or refresh view if needed
				updatePaginationView(currentPage)
			})
		buttonFlex.AddItem(prevButton, 0, 1, false)
	}

	// Add numbered page buttons
	for page := startPage; page < endPage; page++ {
		page := page            // Capture loop variable
		displayPage := page + 1 // 1-indexed for display
		button := tview.NewButton(fmt.Sprintf("%d", displayPage)).
			SetSelectedFunc(func() {
				updatePaginationView(displayPage)
			})

		// Highlight the current page
		if currentPage == displayPage {
			button.SetLabelColor(tcell.ColorBlack).
				SetStyle(tcell.StyleDefault.Background(tcell.ColorDarkGreen))
		} else {
			button.SetLabelColor(tcell.ColorGrey).
				SetStyle(tcell.StyleDefault.Background(tcell.ColorDefault))
		}

		buttonFlex.AddItem(button, 0, 1, false)
	}

	// Add "Next" button to go to the next page
	if currentPage < totalPages-1 {
		nextButton := tview.NewButton(">").
			SetLabelColor(tcell.ColorGrey).
			SetStyle(tcell.StyleDefault.Background(tcell.ColorDefault)).
			SetSelectedFunc(func() {
				// Update the page in the global state
				currentPage += 1
				// Fetch data or refresh view if needed
				updatePaginationView(currentPage)
			})
		buttonFlex.AddItem(nextButton, 0, 1, false)
	}

	// Add "Last" button to go to the last page
	if currentPage < totalPages-1 {
		lastButton := tview.NewButton(">>").
			SetLabelColor(tcell.ColorGrey).
			SetStyle(tcell.StyleDefault.Background(tcell.ColorDefault)).
			SetSelectedFunc(func() {
				// Update the page in the global state
				currentPage = totalPages
				// Fetch data or refresh view if needed
				updatePaginationView(currentPage)
			})
		buttonFlex.AddItem(lastButton, 0, 1, false)
	}

	// Meta information display at the bottom of the Flex
	metaText := fmt.Sprintf("Total Items: %d | Items Per Page: %d | Page %d/%d", totalItems, itemsPerPage, currentPage, totalPages) // 1-indexed
	metaInfo := tview.NewTextView().
		SetText("[grey]" + metaText + "[-]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetWordWrap(true).
		SetTextStyle(tcell.StyleDefault.Background(tcell.ColorDefault))

	// Add meta information below the pagination buttons
	metaFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	metaFlex.AddItem(buttonFlex, 0, 3, false) // Add buttonFlex above
	metaFlex.AddItem(metaInfo, 0, 1, false)   // Add meta information below

	return metaFlex
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func UpdatePaginationViewUI(currentPage int) {
	if state.GlobalState != nil {
		pagination := NewPaginationComponent(currentPage)
		support.UpdateView(state.GlobalState.PaginationFlex, pagination)
	}
}

func UpdatePaginationState(currentPage int) {
	if state.Pagination != nil {
		state.Pagination.Page = currentPage
	}
}

func updatePaginationView(currentPage int) {
	UpdatePaginationState(currentPage)
	UpdatePaginationViewUI(currentPage)

	// FIXME: This should be done in PR package itself
	// pr.ShowSpinnerFetchPRsByQueryAndUpdatePrList()
}

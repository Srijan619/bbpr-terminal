package pr

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
	"strings"
)

const (
	ROOT_COLOR     = tcell.ColorDefault
	DIR_COLOR      = tcell.ColorBlue
	FILE_COLOR     = tcell.ColorGrey
	ICON_DIRECTORY = "\uf07b " // Folder icon
	ICON_FILE      = "\uf15b " // File icon
)

// NodeReference structure
type NodeReference struct {
	Path  string
	IsDir bool
}

// GenerateDiffStatTree creates the diff stat tree view
func GenerateDiffStatTree(data []types.DiffstatEntry) *tview.TreeView {
	// Create the root node for the tree
	root := tview.NewTreeNode("Root").
		SetColor(ROOT_COLOR)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function to add directories and files
	add := func(target *tview.TreeNode, path string, isDir bool, displayName string) *tview.TreeNode {
		ref := &NodeReference{
			Path:  path,
			IsDir: isDir,
		}
		node := tview.NewTreeNode(displayName). // Display text with icon
							SetReference(ref).
							SetSelectable(true)
		if isDir {
			node.SetColor(DIR_COLOR)
			node.SetExpanded(true)
		} else {
			node.SetColor(FILE_COLOR)
		}
		target.AddChild(node)
		return node
	}

	// Helper function to handle path splitting into directories and files
	createPathTree := func(target *tview.TreeNode, fullPath, diffStatText string) {
		// Split the path into directories (except the last part which is a file)
		parts := strings.Split(fullPath, "/")
		var currentNode = target

		for i, part := range parts {
			// Check if this is the last part (file)
			if i == len(parts)-1 {
				// Add file node with diff stat text
				getCommentSymbolAsync(fullPath, func(commentSymbol string) {
					part = commentSymbol + part // Prepend comment symbol to file name
					displayName := fmt.Sprintf("%s%s | %s", ICON_FILE, part, diffStatText)
					currentNode = add(currentNode, fullPath, false, displayName)

					// Trigger the UI refresh to make sure the node is updated correctly
					// Set the node's expanded state after the comment is fetched
					currentNode.SetExpanded(true) // Ensure that the folder is expanded after comment update
					state.GlobalState.App.Draw()  // Force a redraw of the UI to reflect updates
				})
			} else {
				// Add directory node
				displayName := ICON_DIRECTORY + part
				// Check if directory already exists
				dirExists := false
				for _, child := range currentNode.GetChildren() {
					if child.GetText() == displayName {
						dirExists = true
						currentNode = child
						break
					}
				}
				// If directory does not exist, create it
				if !dirExists {
					currentNode = add(currentNode, strings.Join(parts[:i+1], "/"), true, displayName)
				}
			}
		}
	}

	// Iterate through the diffstat entries and create nodes for files and directories
	for _, entry := range data {
		var fileName string
		if entry.New != nil {
			fileName = entry.New.Path
		} else if entry.Old != nil {
			fileName = entry.Old.Path
		}

		// Prepare diff stat text with + for lines added and - for lines removed
		var diffStatText string
		if entry.LinesAdded > 0 {
			diffStatText = fmt.Sprintf("[%s]+ %d[-]", tcell.ColorGreen, entry.LinesAdded)
		}
		if entry.LinesRemoved > 0 {
			if len(diffStatText) > 0 {
				diffStatText += " | "
			}
			diffStatText += fmt.Sprintf("[%s]- %d[-]", tcell.ColorRed, entry.LinesRemoved)
		}

		// Create the path structure in the tree
		createPathTree(root, fileName, diffStatText)
	}

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if node.IsExpanded() {
			node.SetExpanded(false)
		} else {
			node.SetExpanded(true)
		}
		OpenFileSpecificDiff(node, true)
	})

	tree.SetChangedFunc(func(node *tview.TreeNode) {
		if state.GlobalState.CurrentView != state.GlobalState.DiffStatView { // Avoid flickering when on full screen view
			OpenFileSpecificDiff(node, false)
		}
		node.SetSelectedTextStyle(tcell.StyleDefault.Foreground(tcell.ColorOrange))
	})
	return tree
}

// OpenFileSpecificDiff opens the diff of the selected file
func OpenFileSpecificDiff(node *tview.TreeNode, fullScreen bool) {
	ref := node.GetReference()
	if ref != nil {
		nodeRef, ok := ref.(*NodeReference)
		if ok && !nodeRef.IsDir {
			log.Printf("Fetching content for path: %s", nodeRef.Path)

			// Use the spinner utility for asynchronous fetch
			util.ShowLoadingSpinner(state.GlobalState.DiffDetails, func() (interface{}, error) {
				return bitbucket.FetchBitbucketDiffContent(state.GlobalState.SelectedPR.ID, nodeRef.Path)
			}, func(result interface{}, err error) {
				if err != nil {
					util.UpdateDiffDetailsView(err.Error())
				} else {
					result, ok := result.(string)
					if !ok {
						util.UpdateActivityView("[red]Failed to cast diff details[-]")
						return
					}
					// Retrieve inline comments for the file and add comment markers to lines
					comments := getInlineComments(*state.GlobalState.SelectedPR, nodeRef.Path)

					log.Printf("Geetting inline comment for a file %s  %v", nodeRef.Path, comments)
					util.UpdateDiffDetailsView(util.GenerateColorizedDiffView(result, comments))
				}
			})

			if fullScreen {
				// Set the DiffDetails view as the active root
				state.GlobalState.App.SetRoot(state.GlobalState.DiffDetails, true)
			}
		}
	}
}

// getInlineComments fetches the inline comments for a file
func getInlineComments(pr types.PR, file string) []types.Comment {
	resultCh := make(chan []types.Comment)

	go func() {
		comments := bitbucket.FetchBitbucketComments(pr.ID)
		var inlineComments []types.Comment

		for _, comment := range comments {
			if !comment.Deleted && comment.Inline.Path != "" && comment.Inline.Path == file {
				inlineComments = append(inlineComments, comment)
			}
		}

		resultCh <- inlineComments
	}()

	return <-resultCh
}

// getCommentSymbolAsync fetches the comment symbol asynchronously and invokes the callback once done
func getCommentSymbolAsync(path string, callback func(string)) {
	go func() {
		comments := getInlineComments(*state.GlobalState.SelectedPR, path)
		if len(comments) > 0 {
			callback("[yellow]" + ICON_COMMENT + "[-]") // Show the comment icon
		} else {
			callback("") // No comment icon
		}
	}()
}


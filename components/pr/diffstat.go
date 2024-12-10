package pr

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
)

const (
	ROOT_COLOR     = tcell.ColorYellow
	DIR_COLOR      = tcell.ColorBlue
	FILE_COLOR     = tcell.ColorGrey
	ICON_DIRECTORY = "\uf07b "
	ICON_FILE      = "\uf15b "
)

func GenerateDiffStatTree(data []types.DiffstatEntry) *tview.TreeView {
	// Create the root node for the tree
	root := tview.NewTreeNode("Root").
		SetColor(ROOT_COLOR)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function to add directories and files
	add := func(target *tview.TreeNode, path string, isDir bool, fullPath string) *tview.TreeNode {
		node := tview.NewTreeNode(path).
			SetReference(fullPath).
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
	createPathTree := func(target *tview.TreeNode, fullPath string, fileNameWithDiffStatText string) {
		// Split the path into directories (except the last part which is a file)
		parts := strings.Split(fullPath, "/")
		var currentNode = target

		for i, part := range parts {
			// Check if this is the last part (file)
			if i == len(parts)-1 {
				// This is the file, so add the file node
				currentNode = add(currentNode, ICON_FILE+fileNameWithDiffStatText, false, fullPath) // Add file node
			} else {
				// This is a directory, check if directory already exists
				dirExists := false
				for _, child := range currentNode.GetChildren() {
					if child.GetText() == ICON_DIRECTORY+part {
						dirExists = true
						currentNode = child
						break
					}
				}
				// If directory does not exist, create it
				if !dirExists {
					currentNode = add(currentNode, ICON_DIRECTORY+part, true, fullPath) // Add directory node
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
		fileNameWithDiffStatText := (fmt.Sprintf("%s | %s", fileName, diffStatText)) // filename is file with path
		createPathTree(root, fileName, fileNameWithDiffStatText)

	}

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if node.IsExpanded() {
			node.SetExpanded(false)
		} else {
			node.SetExpanded(true)
		}

		ref := node.GetReference()
		if ref != nil {
			fullPath, ok := ref.(string)
			if ok {
				util.UpdateDiffDetailsView(util.GenerateFileContentDiffView(state.GlobalState.SelectedPR.Source.Branch.Name, state.GlobalState.SelectedPR.Destination.Branch.Name, fullPath))
				state.GlobalState.App.SetRoot(state.GlobalState.DiffDetails, true)
			}
		}
	})

	return tree
}

package pr

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strings"

	"simple-git-terminal/apis/bitbucket"
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

type NodeReference struct {
	Path  string
	IsDir bool
}

func GenerateDiffStatTree(data []types.DiffstatEntry) *tview.TreeView {
	// Create the root node for the tree
	root := tview.NewTreeNode("Root").
		SetColor(ROOT_COLOR)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function to add directories and files
	add := func(target *tview.TreeNode, path string, isDir bool) *tview.TreeNode {
		ref := &NodeReference{
			Path:  path,
			IsDir: isDir,
		}
		node := tview.NewTreeNode(path).
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
	createPathTree := func(target *tview.TreeNode, fullPath string) {
		// Split the path into directories (except the last part which is a file)
		parts := strings.Split(fullPath, "/")
		var currentNode = target

		for i, part := range parts {
			// Check if this is the last part (file)
			if i == len(parts)-1 {
				// This is the file, so add the file node
				currentNode = add(currentNode, ICON_FILE+part, false) // Add file node
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
					currentNode = add(currentNode, ICON_DIRECTORY+part, true) // Add directory node
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
		createPathTree(root, fileNameWithDiffStatText)

	}

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if node.IsExpanded() {
			node.SetExpanded(false)
		} else {
			node.SetExpanded(true)
		}
		OpenFileSpecificDiff(node)
	})

	return tree
}

func OpenFileSpecificDiff(node *tview.TreeNode) {
	ref := node.GetReference()
	if ref != nil {
		nodeRef, ok := ref.(*NodeReference)
		if ok && !nodeRef.IsDir {
			log.Printf("Fetching content for path: %s", nodeRef.Path)

			// Show a loading placeholder immediately
			util.UpdateDiffDetailsView("Loading...")
			state.GlobalState.App.SetRoot(state.GlobalState.DiffDetails, true)

			// Fetch content
			content, err := bitbucket.FetchBitbucketDiffContent(state.GlobalState.SelectedPR.ID, nodeRef.Path)
			if err != nil {
				util.UpdateDiffDetailsView(err)
			} else {
				util.UpdateDiffDetailsView(util.GenerateColorizedDiffView(content))
			}
			// TODO: This is for local diff, maybe does not make sense?
			//	util.UpdateDiffDetailsView(util.GenerateFileContentDiffView(state.GlobalState.SelectedPR.Source.Branch.Name, state.GlobalState.SelectedPR.Destination.Branch.Name, fullPath))
		}
	}
}

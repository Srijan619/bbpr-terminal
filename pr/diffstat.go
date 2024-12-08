package pr

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"simple-git-terminal/types"
)

const (
	ROOT_COLOR     = tcell.ColorYellow
	DIR_COLOR      = tcell.ColorBlue
	FILE_COLOR     = tcell.ColorGrey
	ICON_DIRECTORY = "\uf07b "
	ICON_FILE      = "\uf15b "
)

var STATIC_DATA = []types.DiffstatEntry{
	{
		Type:         "diffstat",
		LinesAdded:   1,
		LinesRemoved: 0,
		Status:       "added",
		Old:          nil,
		New: &types.DiffFile{
			Path:        "newDir/simple.js",
			Type:        "commit_file",
			EscapedPath: "newDir/simple.js",
			Links: struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			}{
				Self: struct {
					Href string `json:"href"`
				}{
					Href: "https://api.bitbucket.org/2.0/repositories/chapssrijan619/test_repo/src/51e9e27667998ef8dd96d8197783bb838f734ec5/newDir/simple.js",
				},
			},
		},
	},
	{
		Type:         "diffstat",
		LinesAdded:   0,
		LinesRemoved: 1,
		Status:       "removed",
		Old: &types.DiffFile{
			Path:        "test.txt",
			Type:        "commit_file",
			EscapedPath: "test.txt",
			Links: struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			}{
				Self: struct {
					Href string `json:"href"`
				}{
					Href: "https://api.bitbucket.org/2.0/repositories/chapssrijan619/test_repo/src/02f3c13c7f931ef03df9b86676ed29095a8b2ed5/test.txt",
				},
			},
		},
		New: nil,
	},
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
		node := tview.NewTreeNode(path).
			SetReference(path).
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
	createPathTree := func(target *tview.TreeNode, path string, fileNameWithDiffStatText string) {
		// Split the path into directories (except the last part which is a file)
		parts := strings.Split(path, "/")
		var currentNode = target

		for i, part := range parts {
			// Check if this is the last part (file)
			if i == len(parts)-1 {
				// This is the file
				currentNode = add(currentNode, ICON_FILE+fileNameWithDiffStatText, false) // Add file node
			} else {
				// This is a directory
				currentNode = add(currentNode, ICON_DIRECTORY+part, true) // Add directory node
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
		fileNameWithDiffStatText := (fmt.Sprintf("%s | %s", fileName, diffStatText))
		createPathTree(root, fileName, fileNameWithDiffStatText)

	}

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if node.IsExpanded() {
			node.SetExpanded(false)
		} else {
			node.SetExpanded(true)
		}
	})

	return tree
}

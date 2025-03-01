package util

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"simple-git-terminal/types"
	"strings"

	"github.com/rivo/tview"
)

var ICON_COMMENT = "\uf27b "

// Get the name of the current Git repository
// Fetches the Bitbucket workspace and repo slug based on the current git repo.
func GetRepoAndWorkspace() (string, string, error) {
	// Run git remote to get the remote URL
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = getCurrentDir()
	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get git remote URL: %v", err)
	}

	// The output will look like https://bitbucket.org/workspace/repo
	repoURL := strings.TrimSpace(string(out))

	re := regexp.MustCompile(`bitbucket\.org[:/](?P<workspace>[\w-]+)/(?P<repo>[\w-]+)`)
	matches := re.FindStringSubmatch(repoURL)

	if len(matches) < 3 {
		return "", "", fmt.Errorf("failed to parse workspace and repo from URL: %s", repoURL)
	}

	workspace := matches[1]
	repoSlug := matches[2]

	return workspace, repoSlug, nil
}

func GenerateColorizedDiffView(diffText string, comments []types.Comment) *tview.Flex {
	// Create the main diff text view
	diffView := CreateTextviewComponent("", false)

	// Split and colorize diff lines
	var coloredDiff []string
	lines := strings.Split(diffText, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+") {
			coloredDiff = append(coloredDiff, fmt.Sprintf("[green]%s[-]", line))
		} else if strings.HasPrefix(line, "-") {
			coloredDiff = append(coloredDiff, fmt.Sprintf("[red]%s[-]", line)) // Fixed typo (extra "i")
		} else {
			coloredDiff = append(coloredDiff, line)
		}
	}

	// Set colorized diff text
	diffView.SetText(strings.Join(coloredDiff, "\n"))

	// Create a flex container for structuring the layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(diffView, 0, 3, false) // Diff view takes majority of space

	// Add reply input field at the bottom
	reply := CreateInputFieldComponent("🗜️Reply", "type something....")
	layout.AddItem(reply, 1, 1, true)

	return layout
}

func getCurrentDir() string {
	// For testing during local development override
	if os.Getenv("BBPR_APP_ENV") == "development" {
		return "/Users/srijanpersonal/personal_workspace/raw/test_repo"
	} else {
		return "."
	}
}

// Helper method to add comments above diff lines based on line numbers
func addCommentsAboveLines(diffText string, comments []types.Comment) string {
	diffText = removeBeforeAndIncludingHunk(diffText)
	// Map to track which lines have comments
	commentMap := make(map[int][]types.Comment)

	// Loop through comments and add them to the map, associating them with lines
	for _, comment := range comments {
		// Use `From` and `To` to determine the affected lines
		startLine := comment.Inline.From
		endLine := comment.Inline.To
		// If either From or To is nil, set a default behavior
		if startLine == 0 || endLine == 0 {
			// If From or To are nil, treat this as a single-line comment (use From or To as the same line)
			if startLine != 0 {
				endLine = startLine
			} else if endLine != 0 {
				startLine = endLine
			}
		}

		// Insert comments for each affected line
		for line := startLine; line <= endLine; line++ {
			commentMap[line] = append(commentMap[line], comment)
		}
	}

	// Split the diff into lines
	lines := strings.Split(diffText, "\n")
	var result []string

	// Loop through diff lines and insert comments above
	for i, line := range lines {
		lineNumber := i + 1
		relativeLineNumber := i
		if commentLines, exists := commentMap[relativeLineNumber]; exists {
			// Add each comment as a line before the diff line
			for _, comment := range commentLines {
				commentLine := ""

				if comment.Parent.ID > 0 {
					commentLine = fmt.Sprintf("[steelblue]  %s %s %s %s[-]", ICON_SIDE_ARROW, ICON_COMMENT, comment.User.DisplayName, RenderMarkdown(comment.Content.Raw))
				} else {
					// Need to check if the comment was resolved
					if comment.Resolution != nil {
						// Comment is resolved, show a "resolved" marker
						commentLine = fmt.Sprintf("[aquamarine]✔ %s %s (Resolved) %s[-]", ICON_COMMENT, comment.User.DisplayName, RenderMarkdown(comment.Content.Raw))
					} else {
						// If not resolved, display it as normal
						commentLine = fmt.Sprintf("[steelblue]%s %s → %s[-]", ICON_COMMENT, comment.User.DisplayName, RenderMarkdown(comment.Content.Raw))
					}
				}
				// Add the comment line to the result
				result = append(result, commentLine)
			}
		}
		// Add the diff line itself alongside line number
		lineWithNumber := fmt.Sprintf("[grey]%d[-] %s", lineNumber, line)
		result = append(result, lineWithNumber)
	}
	return strings.Join(result, "\n")
}

// Remove diff hunks as they are unnecessary
func removeBeforeAndIncludingHunk(diffText string) string {
	index := strings.Index(diffText, "@@")
	if index == -1 {
		return diffText
	}
	newlineIndex := strings.Index(diffText[index:], "\n")
	if newlineIndex == -1 {
		return diffText[index+2:] // Skip the @@ directly
	}
	return diffText[index+newlineIndex+1:]
}

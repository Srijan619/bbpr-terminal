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

var (
	ICON_COMMENT  = "\uf27b "
	ICON_MARKED   = "[yellow]★[-]"
	ICON_UNMARKED = " "
)

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

func getCurrentDir() string {
	// For testing during local development override
	if os.Getenv("BBPR_APP_ENV") == "development" {
		return "/Users/srijanpersonal/personal_workspace/raw/test_repo"
	} else {
		return "."
	}
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

func formatCommentWithBox(comment types.Comment) string {
	const maxCommentWidth = 40 // 50% of 80-char screen
	content := ""

	if comment.Parent.ID > 0 {
		content = fmt.Sprintf("[steelblue] %s %s %s %s[-]", ICON_SIDE_ARROW, ICON_COMMENT, comment.User.DisplayName, RenderMarkdown(comment.Content.Raw))
	} else if comment.Resolution != nil {
		content = fmt.Sprintf("[aquamarine]✔ %s %s (Resolved) %s[-]", ICON_COMMENT, comment.User.DisplayName, RenderMarkdown(comment.Content.Raw))
	} else {
		content = fmt.Sprintf("[steelblue]%s %s → %s[-]", ICON_COMMENT, comment.User.DisplayName, RenderMarkdown(comment.Content.Raw))
	}
	return content
}

// GenerateColorizedDiffView (unchanged from last working structure but using new formatCommentWithBox)
func GenerateColorizedDiffView(diffText string, comments []types.Comment) *tview.Table {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	diffText = removeBeforeAndIncludingHunk(diffText)
	lines := strings.Split(diffText, "\n")
	commentMap := make(map[int][]types.Comment)
	markedLines := make(map[int]bool)

	for _, comment := range comments {
		startLine := comment.Inline.From
		endLine := comment.Inline.To
		if startLine == 0 || endLine == 0 {
			if startLine != 0 {
				endLine = startLine
			} else if endLine != 0 {
				startLine = endLine
			}
		}
		for line := startLine; line <= endLine; line++ {
			commentMap[line] = append(commentMap[line], comment)
		}
	}

	row := 0
	for i, line := range lines {
		lineNumber := i + 1
		relativeLineNumber := i

		// Diff line
		color := ""
		if strings.HasPrefix(line, "+") {
			color = "[green]"
		} else if strings.HasPrefix(line, "-") {
			color = "[red]"
		}
		marked := markedLines[relativeLineNumber]
		markIcon := ICON_UNMARKED
		if marked {
			markIcon = ICON_MARKED
		}
		lineText := fmt.Sprintf("%s[grey]%d[-] %s%s[-]", markIcon, lineNumber, color, line)
		table.SetCell(row, 0, tview.NewTableCell(lineText).
			SetExpansion(1).
			SetReference(relativeLineNumber))
		row++

		// Comments beneath
		if commentLines, exists := commentMap[relativeLineNumber]; exists {
			fullCommentText := ""
			for _, comment := range commentLines {
				fullCommentText += formatCommentWithBox(comment) + "\n"
			}
			fullCommentText = strings.TrimSpace(fullCommentText)
			table.SetCell(row, 0, tview.NewTableCell(fullCommentText).
				SetExpansion(1).
				SetReference(commentLines[0]))
			row++
		}
	}

	return table
}

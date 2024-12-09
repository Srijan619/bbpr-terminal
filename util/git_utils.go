package util

import (
	"fmt"
	"github.com/rivo/tview"
	"os/exec"
	"regexp"
	"strings"
)

// Get the name of the current Git repository
// Fetches the Bitbucket workspace and repo slug based on the current git repo.
func GetRepoAndWorkspace() (string, string, error) {
	// Run git remote to get the remote URL
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = "." // Ensure we run it from the current directory
	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get git remote URL: %v", err)
	}

	// The output will look like https://bitbucket.org/workspace/repo
	repoURL := strings.TrimSpace(string(out))

	// Regular expression to extract workspace and repository slug from the URL
	re := regexp.MustCompile(`bitbucket\.org[:/](?P<workspace>[\w-]+)/(?P<repo>[\w-]+)`)
	matches := re.FindStringSubmatch(repoURL)

	if len(matches) < 3 {
		return "", "", fmt.Errorf("failed to parse workspace and repo from URL: %s", repoURL)
	}

	workspace := matches[1]
	repoSlug := matches[2]

	return workspace, repoSlug, nil
}

func GenerateLocalDiffView(source string, destination string, filePath string) *tview.TextView {
	// Run git remote to get the remote URL
	cmd := exec.Command("git", "diff", destination, source, "--", filePath)

	// TODO: Only for testing remove while gitting
	//cmd.Dir = "." // Ensure we run it from the current directory
	cmd.Dir = "/Users/srijanpersonal/personal_workspace/raw/test_repo"
	// Capture the output of the command
	out, err := cmd.CombinedOutput()
	if err != nil {
		textView := GenerateColorizedDiffView(fmt.Sprintf("Error running command %s %v\nOutput: %s", cmd, err, string(out)))
		return textView
	}

	return GenerateColorizedDiffView(string(out))
}

func GenerateColorizedDiffView(diffText string) *tview.TextView {
	// Initialize the TextView to display the diff
	textView := tview.NewTextView()

	// Set options for better readability
	textView.SetDynamicColors(true).
		SetWrap(true).
		SetScrollable(true).
		SetBorder(true).
		SetTitle("Diff View").
		SetTitleAlign(tview.AlignLeft)

	// Split the diff text by lines and color them based on the prefix (+ or -)
	var coloredDiff []string
	lines := strings.Split(diffText, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+") {
			// Green for added lines
			coloredDiff = append(coloredDiff, fmt.Sprintf("[green]%s[-]", line))
		} else if strings.HasPrefix(line, "-") {
			// Red for removed lines
			coloredDiff = append(coloredDiff, fmt.Sprintf("[red]%si[-]", line))
		} else {
			// Normal lines without changes
			coloredDiff = append(coloredDiff, line)
		}
	}

	// Join the lines back together and set the text in the TextView
	textView.SetText(strings.Join(coloredDiff, "\n"))
	return textView
}

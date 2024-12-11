package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/rivo/tview"
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

func GenerateFileContentDiffView(source string, destination string, filePath string) *tview.TextView {
	// Run git remote to get the remote URL
	cmd := exec.Command("git", "diff", destination, source, "--", filePath)
	cmd.Dir = getCurrentDir()

	out, err := cmd.CombinedOutput()
	if err != nil {
		textView := GenerateColorizedDiffView(fmt.Sprintf("Error running command %s %v\nOutput: %s", cmd, err, string(out)))
		return textView
	}

	log.Printf("Return file content diff view...%s", cmd)
	return GenerateColorizedDiffView(string(out))
}

func GenerateColorizedDiffView(diffText string) *tview.TextView {
	// Initialize the TextView to display the diff
	textView := tview.NewTextView()

	// Set options for better readability
	textView.SetDynamicColors(true).
		SetWrap(true).
		SetScrollable(true)

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

func getCurrentDir() string {
	// For testing during local development override
	if os.Getenv("BBPR_APP_ENV") == "development" {
		return "/Users/srijanpersonal/personal_workspace/raw/test_repo"
	} else {
		return "."
	}
}

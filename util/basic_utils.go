package util

import (
	"fmt"
	"os/exec"
	"regexp"
	"simple-git-terminal/types"
	"strings"
	"unicode"
)

func removeZeroWidth(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	for _, r := range input {
		if !unicode.Is(unicode.Mn, r) && r != '\u200C' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func SanitizePR(pr types.PR) types.PR {
	pr.Title = removeZeroWidth(pr.Title)
	if desc, ok := pr.Description.(string); ok {
		pr.Description = removeZeroWidth(desc)
	}
	return pr
}

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

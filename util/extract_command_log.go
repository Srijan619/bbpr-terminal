package util

import (
	"bufio"
	"fmt"
	"strings"
)

// ExtractCommandLog extracts the log segment of a particular command from full step logs.
// The fullLog is the entire step log text, command is the command string to match (e.g. "pnpm install --frozen-lockfile").
func ExtractCommandLog(fullLog, commandName string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(fullLog))

	targetLine := "+ " + commandName
	var builder strings.Builder
	capturing := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "+ ") {
			if capturing {
				// We hit next command start — stop capturing
				break
			}
			if line == targetLine {
				// Start capturing from this line
				capturing = true
				builder.WriteString(line + "\n")
			}
			continue
		}

		if capturing {
			builder.WriteString(line + "\n")
		}
	}

	if !capturing {
		return "", fmt.Errorf("command log segment not found for command: %s", commandName)
	}

	return builder.String(), nil
}

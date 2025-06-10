package config

import (
	"fmt"
	"os"
	"strings"
)

const PacmanConfPath = "/etc/pacman.conf"

type Repo struct {
	Name    string
	Enabled bool
}

// Save writes the modified repository statuses back to the pacman.conf file.
func Save(repos []Repo) error {
	input, err := os.ReadFile(PacmanConfPath)
	if err != nil {
		return fmt.Errorf("failed to read pacman.conf for saving: %w", err)
	}

	lines := strings.Split(string(input), "\n")

	repoMap := make(map[string]bool)
	for _, repo := range repos {
		repoMap[repo.Name] = repo.Enabled
	}

	// Findthe repo header, and any subsequent 'Include' or 'Server' lines
	// that belong to it, commented out or not.
	for index, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		potentialRepoLine := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "#"))

		if strings.HasPrefix(potentialRepoLine, "[") && strings.HasSuffix(potentialRepoLine, "]") {
			repoName := strings.Trim(potentialRepoLine, "[]")

			if enabled, ok := repoMap[repoName]; ok {
				// This is one of our managed repos. Toggle it and its children.
				toggleSection(lines, index, enabled)
			}
		}
	}

	output := strings.Join(lines, "\n")
	// Use a temporary file so we don't corrupt the original on failure.
	tmpFile, err := os.CreateTemp(os.TempDir(), "pacman.conf-")
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}

	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			// Log the error, but don't return it.  We're already returning an error from Save.
			fmt.Printf("Error removing temporary file: %v\n", err)
		}
	}()

	if _, err := tmpFile.WriteString(output); err != nil {
		if err := tmpFile.Close(); err != nil {
			return fmt.Errorf("failed to close temp file: %w", err)
		}

		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Overwrite the original file with the temp file's contents.
	err = os.WriteFile(PacmanConfPath, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("failed to overwrite %s: %w", PacmanConfPath, err)
	}

	return nil
}

// toggleSection comments or uncomments a repository section starting from a given line index.
func toggleSection(lines []string, startIndex int, enable bool) {
	for index := startIndex; index < len(lines); index++ {
		// Stop when we hit the next section or an empty line
		trimmed := strings.TrimSpace(lines[index])
		if trimmed == "" || (strings.HasPrefix(trimmed, "[") && index != startIndex) {
			break
		}

		// Don't modify pure comment lines (like documentation/examples)
		if strings.HasPrefix(trimmed, "#") &&
			!strings.Contains(trimmed, "[") &&
			!strings.Contains(trimmed, "Include") &&
			!strings.Contains(trimmed, "CacheServer") &&
			!strings.Contains(trimmed, "Server") &&
			!strings.Contains(trimmed, "SigLevel") &&
			!strings.Contains(trimmed, "Usage") {
			continue
		}

		isCommented := strings.HasPrefix(trimmed, "#")
		if enable && isCommented {
			// Uncomment the line. We find the first '#' and remove it and any space after.
			if idx := strings.Index(lines[index], "#"); idx != -1 {
				lines[index] = lines[index][:idx] + strings.TrimPrefix(lines[index][idx+1:], " ")
			}
		} else if !enable && !isCommented && trimmed != "" {
			// Comment the line out, preserving indentation.
			indent := len(lines[index]) - len(strings.TrimLeft(lines[index], " \t"))
			lines[index] = strings.Repeat(" ", indent) + "# " + strings.TrimLeft(lines[index], " \t")
		}
	}
}

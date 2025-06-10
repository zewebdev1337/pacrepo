package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zewebdev1337/pacrepo/internal/config"
)

// Example/placeholder repositories to ignore.
var blacklistedRepos = []string{
	"custom",
	"repo-name",
}

// Parse reads the pacman.conf file and extracts repository information.
func Parse() ([]config.Repo, error) {
	file, err := os.Open(config.PacmanConfPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("pacman config not found at %s: %w", config.PacmanConfPath, os.ErrNotExist)
		}

		return nil, fmt.Errorf("could not open %s: %w", config.PacmanConfPath, err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error, but don't fail the function.  File closing is best-effort.
			fmt.Printf("error closing file: %v\n", closeErr)
		}
	}()

	var repos []config.Repo

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines or standard comments
		if trimmedLine == "" || !strings.Contains(trimmedLine, "[") {
			continue
		}

		isCommented := strings.HasPrefix(trimmedLine, "#")
		// Temporarily remove the comment hash for parsing
		if isCommented {
			trimmedLine = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "#"))
		}

		// Now check if it's a valid repo header format
		if strings.HasPrefix(trimmedLine, "[") && strings.HasSuffix(trimmedLine, "]") {
			repoName := strings.Trim(trimmedLine, "[]")

			// Ignore the [options] section and blacklisted repos
			if repoName == "options" || isBlacklisted(repoName) {
				continue
			}

			repos = append(repos, config.Repo{
				Name:    repoName,
				Enabled: !isCommented,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning pacman.conf: %w", err)
	}

	return repos, nil
}

// isBlacklisted checks if a repository name is in our blacklist.
func isBlacklisted(repoName string) bool {
	for _, blacklisted := range blacklistedRepos {
		if strings.EqualFold(repoName, blacklisted) {
			return true
		}
	}

	return false
}

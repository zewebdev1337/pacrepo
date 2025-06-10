package tui

import (
	"fmt"
	"os/exec"

	"github.com/zewebdev1337/pacrepo/internal/config"
	"github.com/zewebdev1337/pacrepo/internal/parser"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application's state.
type Model struct {
	repos     []config.Repo // List of repositories.
	cursor    int           // Cursor position.
	statusMsg string        // Status message to display.
	isLoading bool          // Indicates if repos are being loaded.
	saving    bool          // Indicates if the config is being saved.
	quitting  bool          // Indicates if the application should quit.
	spinner   spinner.Model // Spinner for loading or saving.
	err       error         // Stores any error that occurred.
}

// InitialModel initializes the TUI model.
func InitialModel() Model {
	_spinner := spinner.New()                                              // Create a new spinner.
	_spinner.Spinner = spinner.Dot                                         // Set spinner style to Dot.
	_spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // Set spinner style.

	return Model{
		isLoading: true,     // Initially, we are loading.
		spinner:   _spinner, // Set the spinner model.
	}
}

// loadReposCmd runs the config.Parse function and returns the result as a message.
func loadReposCmd() tea.Msg {
	repos, err := parser.Parse() // Parse the configuration.
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err) // Return the error as a wrapped error.
	}

	return reposLoadedMsg(repos) // Return the loaded repos in a custom message type.
}

// saveReposCmd runs the config.Save function.
func saveReposCmd(repos []config.Repo) tea.Cmd {
	return func() tea.Msg {
		err := config.Save(repos) // Save the configuration.
		if err != nil {
			return saveFinishedMsg{err} // Return an error message if the save fails.
		}
		// On successful save, return a special message
		return saveFinishedMsg{nil} // Return a message indicating successful save.
	}
}

// pacmanSyncCmd runs `pacman -Sy` and returns an error if it fails.
func pacmanSyncCmd() tea.Msg {
	// Execute pacman -Sy
	syncCmd := exec.Command("pacman", "-Sy") // Create a command to sync the pacman database.

	output, err := syncCmd.CombinedOutput() // Execute the command and get the output.
	if err != nil {
		// Return a more detailed error message
		return syncFinishedMsg{fmt.Errorf("`pacman -Sy` failed: %w\nOutput:\n%s", err, string(output))}
	}

	return syncFinishedMsg{nil} // Return a message indicating successful sync.
}

// Custom message types to make the Update function clearer.
type reposLoadedMsg []config.Repo        // Message for when repositories are loaded.
type saveFinishedMsg struct{ err error } // Message for when saving finishes.
type syncFinishedMsg struct{ err error } // Message for when sync finishes.

package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the TUI.
func (m Model) Init() tea.Cmd {
	// Start the spinner and load the repos
	return tea.Batch(m.spinner.Tick, loadReposCmd)
}

// Update handles incoming messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle quitting state separately to reduce nesting
	if m.quitting {
		return m.handleQuittingMessages(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle key press messages
		return m.handleKeyMessages(msg)
	case reposLoadedMsg:
		// Handle the message when repositories are loaded.
		m.isLoading = false
		m.repos = msg

		if len(m.repos) == 0 {
			m.statusMsg = "No repositories found in /etc/pacman.conf."
		} else {
			m.statusMsg = "Loaded repositories."
		}

		return m, nil

	case saveFinishedMsg:
		// Handle the message when saving repositories finishes.
		return m.handleSaveFinished(msg)
	case syncFinishedMsg:
		// Handle the message when synchronizing repositories finishes.
		return m.handleSyncFinished(msg)
	case error:
		// Handle error messages.
		m.err = msg
		return m, tea.Quit

	case spinner.TickMsg:
		// Update the spinner on each tick message.
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}

	return m, nil
}

// handleQuittingMessages handles messages when the application is in quitting state.
func (m Model) handleQuittingMessages(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case saveFinishedMsg:
		// Handle the save finished message.
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.statusMsg = "Save successful! Quitting."
		}

		return m, tea.Quit

	case spinner.TickMsg:
		// Update the spinner.
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd

	default:
		return m, nil
	}
}

// handleKeyMessages handles key press messages.
func (m Model) handleKeyMessages(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		// Quit the application on Ctrl+C or 'q'.
		m.quitting = true
		return m, tea.Quit

	case "up", "k":
		// Move the cursor up.
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		// Move the cursor down.
		if m.cursor < len(m.repos)-1 {
			m.cursor++
		}

	case "enter", " ":
		// Toggle repository enabled state on Enter or Space.
		if len(m.repos) > 0 && !m.saving {
			m.repos[m.cursor].Enabled = !m.repos[m.cursor].Enabled
			m.statusMsg = fmt.Sprintf("Toggled %s. Press 's' to save.", m.repos[m.cursor].Name)
		}

	case "s":
		// Start saving and quitting on 's'.
		if !m.saving {
			return m.startSaveAndQuit()
		}

	case "w":
		// Start saving and synchronizing on 'w'.
		if !m.saving {
			return m.startSaveAndSync()
		}
	}

	return m, nil
}

// startSaveAndQuit starts the save process and sets the application to quit.
func (m Model) startSaveAndQuit() (tea.Model, tea.Cmd) {
	m.saving = true
	m.quitting = true
	m.statusMsg = "Saving changes to /etc/pacman.conf..."

	return m, saveReposCmd(m.repos)
}

// startSaveAndSync starts the save and sync process.
func (m Model) startSaveAndSync() (tea.Model, tea.Cmd) {
	m.saving = true
	m.statusMsg = "Saving changes to /etc/pacman.conf..."

	return m, saveReposCmd(m.repos)
}

// handleSaveFinished handles the message received after saving.
func (m Model) handleSaveFinished(msg saveFinishedMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		return m, tea.Quit
	}

	if !m.quitting {
		// 'w' was pressed - proceed to sync
		m.statusMsg = "Save successful! Synchronizing..."
		return m, pacmanSyncCmd
	}
	// 's' was pressed - just save and quit
	m.statusMsg = "Save successful! Quitting."

	return m, tea.Quit
}

// handleSyncFinished handles the message received after synchronizing.
func (m Model) handleSyncFinished(msg syncFinishedMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.statusMsg = fmt.Sprintf("Saved, but failed to sync: %v", msg.err)
	} else {
		m.statusMsg = "Saved and synchronized successfully! Quitting."
	}

	m.saving = false // Entire save+sync process is done
	m.quitting = true

	return m, tea.Quit // Quit after sync
}

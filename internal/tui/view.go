package tui

import (
	"fmt"
	"strings"
)

// View renders the TUI's current state to a string.
func (m Model) View() string {
	// If an error occurred, display it and exit.
	if m.err != nil {
		// Return an error message formatted with the error value.
		return fmt.Sprintf("\nAn error occurred: %v\n\n", m.err)
	}

	// If quitting, show one final status message.
	if m.quitting {
		// Return a quitting message formatted with the status message.
		return fmt.Sprintf("\n%s\n\n", m.statusMsg)
	}

	// Show a spinner while loading, saving, or syncing.
	if m.isLoading {
		// Return a loading message with the spinner view.
		return fmt.Sprintf("\n   %s Loading /etc/pacman.conf...\n\n", m.spinner.View())
	}

	if m.saving {
		// Return a saving message with the spinner and status message.
		return fmt.Sprintf("\n   %s %s\n\n", m.spinner.View(), m.statusMsg)
	}

	// Use a string builder to construct the view efficiently.
	var strBuilder strings.Builder

	// Title
	// Write the title string to the string builder, rendered with title styling.
	strBuilder.WriteString(titleStyle.Render("pacrepo"))
	strBuilder.WriteString("\n\n")

	// List of Repositories
	// Iterate through each repository in the model.
	for index, repo := range m.repos {
		// Determine the cursor string based on the current selection.
		cursor := "  " // non-selected
		if m.cursor == index {
			cursor = "❯ " // cursor!
		}

		// Status text
		// Determine the status text based on the repository's enabled state.
		var status string
		if repo.Enabled {
			status = statusEnabledStyle.Render()
		} else {
			status = statusDisabledStyle.Render()
		}

		// Build the row string
		// Format the row string, including the cursor, repository name, and status.
		row := fmt.Sprintf("%s[%s] %s", cursor, repo.Name, status)

		// Apply styling
		// Apply different styling based on whether the repository is selected.
		if m.cursor == index {
			strBuilder.WriteString(selectedItemStyle.Render(row))
		} else {
			strBuilder.WriteString(itemStyle.Render(row))
		}

		strBuilder.WriteString("\n")
	}

	// Footer Help
	// Write the help text to the string builder, rendered with help styling.
	help := helpStyle.Render("\n↑/↓: navigate • space/enter: toggle • q: quit\ns: save & quit • w: save, sync & quit")
	strBuilder.WriteString(help)

	// Status Message
	// Write the status message to the string builder if it exists.
	if m.statusMsg != "" {
		strBuilder.WriteString("\n\n" + m.statusMsg)
	}

	// Return the final rendered string, wrapped in a document style.
	return docStyle.Render(strBuilder.String())
}

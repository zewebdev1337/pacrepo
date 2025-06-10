package tui

import "github.com/charmbracelet/lipgloss"

var (
	// General.
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	// Title.
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#007BFF")).
			Padding(0, 1)

		// Help / Footer.
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	// List Styles
	// _ = lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderBottom(true).
	// 	BorderForeground(lipgloss.Color("240")).
	// 	MarginRight(2).
	// 	Render("Repositories")

	itemStyle = lipgloss.NewStyle().PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#00BFFF"))

		// Status Styles.
	statusEnabledStyle = lipgloss.NewStyle().
				SetString("[✓] Enabled").
				Foreground(lipgloss.Color("#4CAF50"))

	statusDisabledStyle = lipgloss.NewStyle().
				SetString("[✗] Disabled").
				Foreground(lipgloss.Color("#F44336"))
)

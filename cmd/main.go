package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zewebdev1337/pacrepo/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pacrepo",
	Short: "A TUI to toggle pacman repositories.",
	Long:  `pacrepo provides a simple terminal user interface to toggle your pacman repositories on or off.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if running with sudo
		if os.Geteuid() != 0 {
			fmt.Println("Error: This program must be run with sudo to modify /etc/pacman.conf")
			fmt.Println("Please try: sudo pacrepo")
			os.Exit(1)
		}

		p := tea.NewProgram(tea.Model(tui.InitialModel()), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatalf("Alas, there's been an error: %v", err)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/aemery-cb/pequod/internal/api"
	"github.com/aemery-cb/pequod/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	client := api.CreateClient()
	p := tea.NewProgram(ui.NewWindow(&client), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

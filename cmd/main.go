package main

import (
	"fmt"
	"os"

	"github.com/aemery-cb/pequod/internal/api"
	"github.com/aemery-cb/pequod/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	client := api.CreateClient()
	p := tea.NewProgram(ui.NewWindow(&client), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

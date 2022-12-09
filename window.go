package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type selectedPage int

const (
	nodePage selectedPage = iota
	pvcPage
)

type Window struct {
	pageState selectedPage
	pages     []tea.Model
}

func NewWindow(client *Client) Window {

	pages := []tea.Model{
		NewPodModel(client),
		NewModel(),
	}

	return Window{
		pageState: nodePage,
		pages:     pages,
	}
}

func (w Window) Init() tea.Cmd {
	return tea.Batch(w.pages[nodePage].Init())
}

func (m Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			// tell all pages its quitting time.
			for _, x := range m.pages {
				fmt.Println("passing quit")
				x.Update(msg)
			}
			return m, tea.Quit
		case "tab":
			newNum := int(m.pageState+1) % len(m.pages)
			m.pageState = selectedPage(newNum)
			return m, nil
		}
	}

	newBx, cmd := m.pages[m.pageState].Update(msg)
	cmds = append(cmds, cmd)
	m.pages[m.pageState] = newBx
	return m, tea.Batch(cmds...)
}

func (m Window) View() string {
	return m.pages[m.pageState].View()
}

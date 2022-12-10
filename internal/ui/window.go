package ui

import (
	"github.com/aemery-cb/pequod/internal/api"
	"github.com/aemery-cb/pequod/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

type selectedPage int

const (
	nodePage selectedPage = iota
	pvcPage
)

/**
This is the highest level construct and overall wrapper.

ATM only 1 page is viewable at a time and handles focus between pages
**/
type Window struct {
	pageState selectedPage
	pages     []tea.Model
	client    *api.Client
	stop      chan struct{}
	sub       chan tea.Msg
}

func NewWindow(client *api.Client) Window {

	pages := []tea.Model{
		NewPodModel(),
		NewModel(),
	}

	return Window{
		client:    client,
		stop:      make(chan struct{}),
		sub:       make(chan tea.Msg),
		pageState: nodePage,
		pages:     pages,
	}
}

func (w Window) Init() tea.Cmd {
	w.client.WatchPods("", w.sub, w.stop)
	return tea.Batch(w.pages[nodePage].Init(), waitForActivity(w.sub))
}

func (m Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case common.WaitForActivityMsg:
		return m, waitForActivity(m.sub)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			// tell all pages its quitting time so they can gracefully exit.
			for _, x := range m.pages {
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

// non-blocking command to receive new messages from
// kube client.
func waitForActivity(sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		data := <-sub
		return data
	}
}

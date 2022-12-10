package ui

import (
	"context"

	"github.com/aemery-cb/pequod/internal/api"
	"github.com/aemery-cb/pequod/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

type Page interface {
	tea.Model
	Blur() Page
	Focus() Page
}

/**
This is the highest level construct and overall wrapper.

ATM only 1 page is viewable at a time and handles focus between pages
**/
type Window struct {
	pageState common.SelectedPage
	pages     []Page
	client    *api.Client
	stop      chan struct{}
	sub       chan tea.Msg
}

func NewWindow(client *api.Client) Window {

	pages := []Page{
		NewPodModel(),
		NewLogModel(),
	}
	pages[0] = pages[0].Focus()
	return Window{
		client:    client,
		stop:      make(chan struct{}),
		sub:       make(chan tea.Msg),
		pageState: common.NodePage,
		pages:     pages,
	}
}

func (w Window) Init() tea.Cmd {
	w.client.WatchPods("", w.sub, w.stop)
	return tea.Batch(w.pages[common.NodePage].Init(), waitForActivity(w.sub))
}

func (w Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case common.SelectPageMsg:
		w.pageState = msg.PageNumber
	case common.WaitForActivityMsg:
		return w, waitForActivity(w.sub)
	case common.WatchPodLogsMsg:
		w.client.StreamLogs(context.Background(), *msg.Pod, w.sub)
		return w, tea.Batch(common.ClearPodLogs(), common.SelectPage(common.LogPage), common.WaitForActivity())
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			// tell all pages its quitting time so they can gracefully exit.
			w.stop <- struct{}{}
			for _, x := range w.pages {
				x.Update(msg)
			}
			return w, tea.Quit
		case "tab":
			w.pages[w.pageState] = w.pages[w.pageState].Blur()
			newNum := int(w.pageState+1) % len(w.pages)
			w.pages[newNum] = w.pages[newNum].Focus()
			w.pageState = common.SelectedPage(newNum)
			return w, nil
		}
	}

	for index, page := range w.pages {
		newModel, cmd := page.Update(msg)
		w.pages[index], _ = newModel.(Page)
		cmds = append(cmds, cmd)
	}

	return w, tea.Batch(cmds...)
}

func (m Window) View() string {
	return m.pages[m.pageState].View() + "\nPress 'q' to quit"
}

// non-blocking command to receive new messages from
// kube client.
func waitForActivity(sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		data := <-sub
		return data
	}
}

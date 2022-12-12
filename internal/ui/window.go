package ui

import (
	"context"

	"github.com/aemery-cb/pequod/internal/api"
	"github.com/aemery-cb/pequod/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

/**
This is the highest level construct and overall wrapper.

ATM only 1 pane is viewable at a time and handles focus between panes
**/
type Window struct {
	selectedPane common.SelectedPane
	panes        []Pane
	client       *api.Client
	stop         chan struct{}
	sub          chan tea.Msg
}

func NewWindow(client *api.Client) Window {

	pages := []Pane{
		NewPane(NewPodModel()),
		NewPane(NewLogModel()),
	}
	pages[0].Focus()
	return Window{
		client:       client,
		stop:         make(chan struct{}),
		sub:          make(chan tea.Msg),
		selectedPane: common.NodePane,
		panes:        pages,
	}
}

func (w Window) Init() tea.Cmd {
	w.client.WatchPods("", w.sub, w.stop)
	return tea.Batch(w.panes[common.NodePane].Init(), waitForActivity(w.sub))
}

func (w Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case common.SelectPaneMsg:
		w.panes[w.selectedPane].Blur()
		w.selectedPane = msg.PaneNumber
		w.panes[w.selectedPane].Focus()
	case common.WaitForActivityMsg:
		return w, waitForActivity(w.sub)
	case common.WatchPodLogsMsg:
		w.client.StreamLogs(context.Background(), msg.Pod, w.sub)
		cmds = append(cmds, tea.Batch(common.ClearPodLogs(), common.SelectPane(common.LogPane), common.WaitForActivity()))
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			// tell all pages its quitting time so they can gracefully exit.
			w.stop <- struct{}{}
			for _, x := range w.panes {
				x.Update(msg)
			}
			return w, tea.Quit
		case "tab":
			w.panes[w.selectedPane].Blur()
			newNum := int(w.selectedPane+1) % len(w.panes)
			w.panes[newNum].Focus()
			w.selectedPane = common.SelectedPane(newNum)
			return w, nil
		}
	}

	for index, page := range w.panes {
		newModel, cmd := page.Update(msg)
		w.panes[index], _ = newModel.(Pane)
		cmds = append(cmds, cmd)
	}

	return w, tea.Batch(cmds...)
}

func (m Window) View() string {
	return BaseStyle.Render(m.panes[m.selectedPane].View() + "\nPress 'q' to quit")
}

// non-blocking command to receive new messages from
// kube client.
func waitForActivity(sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		data := <-sub
		return data
	}
}

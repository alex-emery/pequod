package main

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
)

type PodStatus struct {
	name   string
	status string
	uptime string
}

func (r PodStatus) String() string {
	return r.name + " : " + r.status + " - " + r.uptime
}

type model struct {
	client   *Client
	pods     []PodStatus
	quitting bool
	stop     chan struct{}
	sub      chan tea.Msg
}

func newModel(client *Client) model {
	return model{
		client: client,
		stop:   make(chan struct{}),
		sub:    make(chan tea.Msg),
	}
}

func waitForActivity(sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		data := <-sub
		return data
	}
}

func (m model) Init() tea.Cmd {
	m.client.WatchPods(m.sub, m.stop)
	return tea.Batch(
		waitForActivity(m.sub),
	)
}

type NewPodMsg struct{ pod *v1.Pod }
type UpdatePodMsg struct {
	old *v1.Pod
	new *v1.Pod
}
type DeletePodMsg struct{ pod *v1.Pod }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NewPodMsg:
		m.pods = append(m.pods, PodStatus{
			name:   msg.pod.Name,
			status: msg.pod.Status.Message,
		})
		sort.Slice(m.pods, func(i, j int) bool {
			return m.pods[i].name < m.pods[j].name
		})
		return m, waitForActivity(m.sub)
	case UpdatePodMsg:
		return m, waitForActivity(m.sub)
	case DeletePodMsg:
		pods := m.pods
		newList := []PodStatus{}
		for _, pod := range pods {
			if pod.name != msg.pod.Name {
				newList = append(newList, pod)
			}
		}
		m.pods = newList
		return m, waitForActivity(m.sub)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.quitting = true
			m.stop <- struct{}{}

			return m, tea.Quit
		case "s":
			m.stop <- struct{}{}
		case "w":
			m.client.WatchPods(m.sub, m.stop)
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	var s string

	if m.quitting {
		s += "Quitting"
	} else {
		s += "Pods"
	}

	s += "\n\n"

	for _, res := range m.pods {
		s += res.String() + "\n"
	}

	if !m.quitting {
		s += "Press 'q' to exit"
	}

	if m.quitting {
		s += "\n"
	}

	return s
}

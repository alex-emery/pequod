package main

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
)

func PrintPod(pod v1.Pod) string {
	return pod.Namespace + " : " + pod.Name
}

type PodModel struct {
	client   *Client
	pods     []v1.Pod
	quitting bool
	stop     chan struct{}
	sub      chan tea.Msg
	cursor   int
}

func NewPodModel(client *Client) PodModel {
	return PodModel{
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

func (m PodModel) Init() tea.Cmd {
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

func (m PodModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NewPodMsg:
		m.pods = append(m.pods, *msg.pod)

		sort.Slice(m.pods, func(i, j int) bool {
			return m.pods[i].Name < m.pods[j].Name
		})
		return m, waitForActivity(m.sub)
	case UpdatePodMsg:
		return m, waitForActivity(m.sub)
	case DeletePodMsg:
		pods := m.pods
		newList := []v1.Pod{}
		for _, pod := range pods {
			if pod.Name != msg.pod.Name {
				newList = append(newList, pod)
			}
		}
		m.pods = newList
		// TODO: check if the cursor position is now invalidated and move if appropriate
		return m, waitForActivity(m.sub)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.quitting = true
			m.stop <- struct{}{}
			return m, nil
		case "s":
			m.stop <- struct{}{}
		case "w":
			m.client.WatchPods(m.sub, m.stop)
			return m, nil
		case "down":
			if m.cursor < len(m.pods)-1 {
				m.cursor += 1
			}
			return m, nil
		case "up":
			if m.cursor > 0 {
				m.cursor -= 1
			}
		}
	}

	return m, nil
}

func (m PodModel) View() string {
	var s string

	if m.quitting {
		s += "Quitting"
	} else {
		s += "Pods"
	}

	s += "\n\n"

	for index, res := range m.pods {
		if index == m.cursor {
			s += ">"
		} else {
			s += " "
		}
		s += PrintPod(res) + "\n"
	}

	if !m.quitting {
		s += "Press 'q' to exit"
	}

	if m.quitting {
		s += "\n"
	}

	return s
}

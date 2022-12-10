package ui

import (
	"sort"

	"github.com/aemery-cb/pequod/internal/common"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
)

func PrintPod(pod v1.Pod) string {
	return pod.Namespace + " : " + pod.Name
}

type PodModel struct {
	pods   []v1.Pod
	cursor int
}

func NewPodModel() PodModel {
	return PodModel{}
}

func (m PodModel) Init() tea.Cmd {
	return nil
}

func (m PodModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.NewPodMsg:
		m.pods = append(m.pods, *msg.Pod)
		sort.Slice(m.pods, func(i, j int) bool {
			return m.pods[i].Name < m.pods[j].Name
		})
		return m, common.WaitForActivity
	case common.UpdatePodMsg:
		return m, common.WaitForActivity
	case common.DeletePodMsg:
		pods := m.pods
		newList := []v1.Pod{}
		for _, pod := range pods {
			if pod.Name != msg.Pod.Name {
				newList = append(newList, pod)
			}
		}
		m.pods = newList
		// TODO: check if the cursor position is now invalidated and move if appropriate
		return m, common.WaitForActivity
	case tea.KeyMsg:
		switch msg.String() {
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

	s += "Pods"

	s += "\n\n"

	for index, res := range m.pods {
		if index == m.cursor {
			s += ">"
		} else {
			s += " "
		}
		s += PrintPod(res) + "\n"
	}

	return s
}

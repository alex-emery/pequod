package common

import (
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
)

func WaitForActivity() tea.Cmd {
	return func() tea.Msg {
		return WaitForActivityMsg{}
	}
}

func WatchPodLogs(pod *v1.Pod) tea.Cmd {
	return func() tea.Msg {
		return WatchPodLogsMsg{Pod: pod}
	}
}

func ClearPodLogs() tea.Cmd {
	return func() tea.Msg {
		return ClearPodLogsMsg{}
	}
}

func SelectPane(number SelectedPane) tea.Cmd {
	return func() tea.Msg {
		return SelectPaneMsg{
			PaneNumber: number,
		}
	}

}

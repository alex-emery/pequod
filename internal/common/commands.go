package common

import tea "github.com/charmbracelet/bubbletea"

func WaitForActivity() tea.Msg {
	return WaitForActivityMsg{}
}

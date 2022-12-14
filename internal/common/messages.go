package common

import v1 "k8s.io/api/core/v1"

type WaitForActivityMsg struct{}

type NewPodMsg struct{ Pod *v1.Pod }

type UpdatePodMsg struct {
	Old *v1.Pod
	New *v1.Pod
}
type DeletePodMsg struct{ Pod *v1.Pod }

// new log has arrived from a pod.
type NewLogMsg struct {
	Pod     *v1.Pod
	Message string
}

// trigger to start streaming logs from a pod
type WatchPodLogsMsg struct {
	Pod *v1.Pod
}

// clear log message screen
type ClearPodLogsMsg struct{}

type SelectPaneMsg struct {
	PaneNumber SelectedPane
}

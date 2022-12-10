package common

import v1 "k8s.io/api/core/v1"

type NewPodMsg struct{ Pod *v1.Pod }

type UpdatePodMsg struct {
	Old *v1.Pod
	New *v1.Pod
}
type DeletePodMsg struct{ Pod *v1.Pod }

type WaitForActivityMsg struct{}

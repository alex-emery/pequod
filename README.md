# Pequod

Quick example of showing Kubernetes pods in bubble tea.

## Quick start

`go run ./cmd/main.go`

use ↑ and ↓ to choose a pod.
press enter on choosen pod to stream logs.
press tab to switch between pod and log view.

## TODO:
- [ ] allow for multiple Pages to be displayed at the same time.
- [x] log views for pods
- [ ] log view needs to be a list so it can be scrolled.
- [ ] namespace filter for pods
## Notes
https://github.com/charmbracelet/bubbletea/blob/master/examples/realtime/main.go

Commands are ran async so


```
func listenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(900)+100))
			sub <- struct{}{}
		}
	}
}
```

wont block.

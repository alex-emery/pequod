# Pequod

Quick example of showing Kubernetes pods in bubble tea.

## Quick start

`go run .`


## TODO:
- [ ] allow for multiple panes to be displayed at the same time.
- [ ] log views for pods
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

# Pequod

Pequod as in the boat, from Moby-Dick. Variable names are hard...

A tool for exploring Kubernetes clusters.

## Quick start

`go run ./cmd/main.go`

use ↑ and ↓ to choose a pod.
press enter on choosen pod to stream logs.
press tab to switch between pod and log view.

## Example
![An example of the pequod program running](./examples/demo.gif)

## TODO:
- [ ] allow for multiple Pages to be displayed at the same time.
- [x] log views for pods
- [x] log view needs to be a list so it can be scrolled.
- [ ] namespace filter for pods


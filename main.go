package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	go func() {
		ListPods(p.Send)
	}()
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func ListPods(send func(tea.Msg)) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		podsStatuses := make([]PodStatus, len(pods.Items))
		for i, pod := range pods.Items {
			status := PodStatus{
				name:   pod.Name,
				status: string(pod.Status.Conditions[0].Type),
				uptime: pod.Status.StartTime.Time.String(),
			}
			podsStatuses[i] = status
		}
		send(PodMsg{
			pods: podsStatuses,
		})
		time.Sleep(5 * time.Second)
	}
}

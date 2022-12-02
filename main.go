package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	client := GetKubeClient()
	p := tea.NewProgram(newModel(&client), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func GetKubeClient() Client {
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

	return Client{kubeClient: clientset}
}

type Client struct {
	kubeClient *kubernetes.Clientset
}

func (c *Client) GetPods() ([]PodStatus, error) {
	pods, err := c.kubeClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	podStatuses := make([]PodStatus, len(pods.Items))
	for i, pod := range pods.Items {
		status := PodStatus{
			name:   pod.Name,
			status: string(pod.Status.Conditions[0].Type),
			uptime: pod.Status.StartTime.Time.String(),
		}
		podStatuses[i] = status
	}

	return podStatuses, nil

}

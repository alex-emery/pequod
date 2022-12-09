package main

import (
	"context"
	"flag"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

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

func (c *Client) WatchPods(sub chan<- tea.Msg, stop <-chan struct{}) {
	watchlist := cache.NewListWatchFromClient(c.kubeClient.CoreV1().RESTClient(), "pods", v1.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{}, 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				sub <- NewPodMsg{pod: pod}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				sub <- DeletePodMsg{pod: pod}
			},
			UpdateFunc: func(old interface{}, new interface{}) {
				oldPod := old.(*v1.Pod)
				newPod := new.(*v1.Pod)
				sub <- UpdatePodMsg{old: oldPod, new: newPod}
			},
		},
	)

	go controller.Run(stop)
}

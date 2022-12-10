package api

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"io"
	"path/filepath"

	"github.com/aemery-cb/pequod/internal/common"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	kubeClient *kubernetes.Clientset
}

// Creates a client which wraps the kube client in helpers.
func CreateClient() Client {
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

func (c *Client) WatchPods(namespace string, sub chan<- tea.Msg, stop <-chan struct{}) {
	watchlist := cache.NewListWatchFromClient(c.kubeClient.CoreV1().RESTClient(), "pods", namespace, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{}, 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				sub <- common.NewPodMsg{Pod: pod}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				sub <- common.DeletePodMsg{Pod: pod}
			},
			UpdateFunc: func(old interface{}, new interface{}) {
				oldPod := old.(*v1.Pod)
				newPod := new.(*v1.Pod)
				sub <- common.UpdatePodMsg{Old: oldPod, New: newPod}
			},
		},
	)

	go controller.Run(stop)
}

func (c *Client) StreamLogs(ctx context.Context, pod v1.Pod, sub chan<- tea.Msg) {
	req := c.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{Follow: true})
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return
	}
	r := bufio.NewReader(podLogs)

	go func() {
		defer podLogs.Close()
		for {
			bytes, err := r.ReadBytes('\n')
			sub <- common.NewLogMsg{Message: string(bytes)}
			if err != nil {
				if !errors.Is(err, io.EOF) {
					return
				}

				break
			}
		}
	}()
}

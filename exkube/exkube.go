package exkube

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientConfig struct {
	QPS        float32
	Burst      int
	Kubeconfig string
}

// New returns a kubernetes client.
// It tries first with in-cluster config, if it fails it will try with out-of-cluster config.
func New(cc *ClientConfig) (client kubernetes.Interface, err error) {
	client, err = NewInCluster(cc)
	if err == nil {
		return
	}
	if len(cc.Kubeconfig) == 0 {
		dir, err := os.UserHomeDir()
		if err != nil {
			return client, err
		}
		cc.Kubeconfig = filepath.Join(dir, ".kube", "config")
	}
	client, err = NewFromConfig(cc)
	if err != nil {
		return
	}

	return
}

// NewFromConfig returns a new out-of-cluster kubernetes client.
func NewFromConfig(cc *ClientConfig) (client kubernetes.Interface, err error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", cc.Kubeconfig)
	if err != nil {
		return
	}

	cc.apply(config)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
	client = clientset
	return
}

// NewInCluster returns a new in-cluster kubernetes client.
func NewInCluster(cc *ClientConfig) (client kubernetes.Interface, err error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return
	}

	cc.apply(config)

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	client = clientset

	return
}

// newRestConfigInCluster returns a new in-cluster kubernetes rest client.
func newRestConfigInCluster(cc *ClientConfig) (config *rest.Config, err error) {
	// creates the in-cluster config
	config, err = rest.InClusterConfig()
	if err != nil {
		return
	}
	cc.apply(config)
	return
}

// newRestConfigFromConfig returns a new out-of-cluster kubernetes client.
func newRestConfigFromConfig(cc *ClientConfig) (config *rest.Config, err error) {
	// use the current context in kubeconfig
	config, err = clientcmd.BuildConfigFromFlags("", cc.Kubeconfig)
	if err != nil {
		return
	}
	cc.apply(config)
	return
}

func NewRestConfig(cc *ClientConfig) (config *rest.Config, err error) {
	config, err = newRestConfigInCluster(cc)
	if err == nil {
		return
	}
	if len(cc.Kubeconfig) == 0 {
		dir, err := os.UserHomeDir()
		if err != nil {
			return config, err
		}
		cc.Kubeconfig = filepath.Join(dir, ".kube", "config")
	}
	config, err = newRestConfigFromConfig(cc)
	if err != nil {
		return
	}
	return
}

func (cc *ClientConfig) apply(config *rest.Config) {
	if cc.QPS > 0.0 {
		config.QPS = cc.QPS // the default is rest.DefaultQPS which is 5.0
	}

	if cc.Burst > 0 {
		config.Burst = cc.Burst // the default is rest.DefaultBurst which is 10
	}
}

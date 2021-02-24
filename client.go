package kubi

import (
	"fmt"
	"os"

	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

// Interface implements the main interface
type Interface interface {
	setRestConfig() error
	setClientset() error
	setDynamicClient() error
	setDiscoveryClient() error
	setRESTMapper() error
}

// Client is the container of Kubernetes
// interfaces that allow the communication
type Client struct {
	RestConfig      *rest.Config
	Clientset       *kubernetes.Clientset
	DynamicClient   interface{}
	DiscoveryClient *discovery.DiscoveryClient
	RESTMapper      *restmapper.DeferredDiscoveryRESTMapper
}

// NewClient returns a new Client
func NewClient() (*Client, error) {
	client := &Client{}

	if err := client.setRestConfig(); err != nil {
		return nil, err
	}

	if err := client.setClientset(); err != nil {
		return nil, err
	}

	if err := client.setDynamicClient(); err != nil {
		return nil, err
	}

	if err := client.setDiscoveryClient(); err != nil {
		return nil, err
	}

	if err := client.setRESTMapper(); err != nil {
		return nil, err
	}

	return client, nil
}

func (client *Client) setRestConfig() error {
	if _, ok := os.LookupEnv("KUBECONFIG"); ok {
		return client.setRestConfigFromKubeconfig()
	}

	// if this application is running
	// inside a Pod with a Service Account
	if _, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST"); ok {
		return client.setRestConfigFromCluster()
	}

	return fmt.Errorf("without kubernetes access")
}

func (client *Client) setRestConfigFromKubeconfig() error {
	restConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	client.RestConfig = restConfig

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) setRestConfigFromCluster() error {
	restConfig, err := rest.InClusterConfig()
	client.RestConfig = restConfig

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) setClientset() error {
	clientSet, err := kubernetes.NewForConfig(client.RestConfig)
	client.Clientset = clientSet

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) setDynamicClient() error {
	_, err := dynamic.NewForConfig(client.RestConfig)

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) setDiscoveryClient() error {
	discoveryClient, err := discovery.
		NewDiscoveryClientForConfig(client.RestConfig)

	client.DiscoveryClient = discoveryClient

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) setRESTMapper() error {
	client.RESTMapper =
		restmapper.
			NewDeferredDiscoveryRESTMapper(
				memory.NewMemCacheClient(client.DiscoveryClient),
			)

	return nil
}

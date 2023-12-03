package kubelib

import (
	"io"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func RESTConfigFromKubeConfig(kubeconfig string) (*rest.Config, error) {
	// Open the kubeconfig file
	kubeconfigFile, err := os.Open(kubeconfig)
	if err != nil {
		return nil, err
	}
	defer kubeconfigFile.Close()

	// Load the kubeconfig content as bytes
	kubeconfigBytes, err := io.ReadAll(kubeconfigFile)
	if err != nil {
		return nil, err
	}

	// Create a rest.Config from kubeconfig bytes.
	return clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
}

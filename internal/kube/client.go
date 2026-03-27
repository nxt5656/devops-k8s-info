package kube

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset kubernetes.Interface
}

func NewClient() (*Client, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		kubeConfigPath := os.Getenv("KUBECONFIG")
		if kubeConfigPath == "" {
			home, homeErr := os.UserHomeDir()
			if homeErr != nil {
				return nil, fmt.Errorf("resolve user home failed: %w", homeErr)
			}
			kubeConfigPath = filepath.Join(home, ".kube", "config")
		}
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("build kube config from %s failed: %w", kubeConfigPath, err)
		}
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client failed: %w", err)
	}

	return &Client{clientset: clientset}, nil
}

func (c *Client) CheckReady(ctx context.Context) error {
	_, err := c.clientset.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("get server version failed: %w", err)
	}
	return nil
}

func (c *Client) ListNamespaces(ctx context.Context) ([]string, error) {
	res, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list namespaces failed: %w", err)
	}

	items := make([]string, 0, len(res.Items))
	for _, ns := range res.Items {
		items = append(items, ns.Name)
	}
	return items, nil
}

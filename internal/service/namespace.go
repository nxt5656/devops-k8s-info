package service

import (
	"context"
	"time"

	"devops-k8s-info/internal/kube"
	"devops-k8s-info/internal/model"
)

type NamespaceService struct {
	kubeClient *kube.Client
	timeout    time.Duration
}

func NewNamespaceService(kubeClient *kube.Client, timeout time.Duration) *NamespaceService {
	return &NamespaceService{
		kubeClient: kubeClient,
		timeout:    timeout,
	}
}

func (s *NamespaceService) List(ctx context.Context) ([]model.Namespace, error) {
	reqCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	names, err := s.kubeClient.ListNamespaces(reqCtx)
	if err != nil {
		return nil, err
	}

	out := make([]model.Namespace, 0, len(names))
	for _, name := range names {
		out = append(out, model.Namespace{
			Name:      name,
			Status:    "Unknown",
			CreatedAt: "",
		})
	}
	return out, nil
}

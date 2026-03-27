package handler

import (
	"net/http"

	"devops-k8s-info/internal/model"
	"devops-k8s-info/internal/service"
	"devops-k8s-info/pkg/response"

	"github.com/gin-gonic/gin"
)

type NamespaceHandler struct {
	svc *service.NamespaceService
}

func NewNamespaceHandler(svc *service.NamespaceService) *NamespaceHandler {
	return &NamespaceHandler{svc: svc}
}

// ListNamespaces godoc
// @Summary List namespaces
// @Description Get all namespaces from the current Kubernetes cluster.
// @Tags namespaces
// @Produce json
// @Success 200 {object} response.Envelope{data=[]model.Namespace}
// @Failure 502 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /namespaces [get]
func (h *NamespaceHandler) ListNamespaces(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusBadGateway, err.Error())
		return
	}
	if items == nil {
		items = make([]model.Namespace, 0)
	}
	response.OK(c, items)
}

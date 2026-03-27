package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"devops-k8s-info/internal/handler"
	"devops-k8s-info/internal/kube"
	"devops-k8s-info/internal/service"
	"devops-k8s-info/pkg/config"

	"devops-k8s-info/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Zhiyi K8s Info API
// @version 1.0
// @description API service for reading Kubernetes resources.
// @BasePath /api/v1/k8s-info
func main() {
	cfg := config.Load()
	basePath := fmt.Sprintf("/api/v1/%s", cfg.APIBaseSegment)
	docs.SwaggerInfo.BasePath = basePath

	kubeClient, err := kube.NewClient()
	if err != nil {
		log.Fatalf("init k8s client failed: %v", err)
	}

	nsService := service.NewNamespaceService(kubeClient, cfg.KubeRequestTimeout)
	nsHandler := handler.NewNamespaceHandler(nsService)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.KubeRequestTimeout)
		defer cancel()

		if err := kubeClient.CheckReady(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group(basePath)
	{
		api.GET("/namespaces", nsHandler.ListNamespaces)
	}

	srv := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server starting on %s", cfg.ServerAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

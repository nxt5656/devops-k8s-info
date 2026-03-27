package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServerAddr         string
	KubeRequestTimeout time.Duration
	KubeConfigPath     string
	APIBaseSegment     string
}

func Load() Config {
	addr := envOrDefault("SERVER_ADDR", ":8080")
	timeoutSeconds := envOrDefaultInt("KUBE_REQUEST_TIMEOUT_SECONDS", 8)
	kubeConfigPath := os.Getenv("KUBECONFIG")
	apiBaseSegment := normalizePathSegment(envOrDefault("API_BASE_SEGMENT", "k8s-info"), "k8s-info")

	return Config{
		ServerAddr:         addr,
		KubeRequestTimeout: time.Duration(timeoutSeconds) * time.Second,
		KubeConfigPath:     kubeConfigPath,
		APIBaseSegment:     apiBaseSegment,
	}
}

func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func envOrDefaultInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return parsed
}

func normalizePathSegment(raw, fallback string) string {
	trimmed := strings.Trim(raw, "/ ")
	if trimmed == "" {
		return fallback
	}
	return trimmed
}

# devops-k8s-info

Kubernetes 资源信息采集 API（Golang）。

项目目标是通过 REST API 提供集群资源只读查询能力，并以容器方式部署到 Kubernetes 集群。

## 功能概览

- 健康检查：`/healthz`、`/readyz`
- 资源查询前缀：`/api/v1/{api_base_segment}`
- Swagger：`/swagger/index.html`

`{api_base_segment}` 由环境变量 `API_BASE_SEGMENT` 控制，默认 `k8s-info`。

## 目录结构

- `cmd/`：程序入口
- `internal/`：业务逻辑
- `pkg/`：公共模块
- `deploy/k8s/`：Kubernetes 清单
- `doc/`：项目说明文档
- `docs/`：Swagger 生成产物

## 本地运行

```bash
go run main.go
```

## 生成 Swagger 文档

```bash
swag init -g cmd/server/main.go -o docs
```

## 镜像构建与发布

默认镜像地址：`ghcr.io/nxt5656/devops-k8s-info`

可通过 GitHub Actions 自动构建并推送镜像（见 `doc/github-actions-image.md`）。

## Kubernetes 部署

当前仓库已提供：

- `deploy/k8s/deployment.yaml`

部署命令：

```bash
kubectl apply -f deploy/k8s/deployment.yaml
```

注意：该 Deployment 使用了 `serviceAccountName: k8s-info-sa`，请先在集群中创建对应的 `ServiceAccount` 与最小权限 RBAC（`get/list/watch`：pods、services、namespaces、deployments），否则 Pod 可能无法正常访问 Kubernetes API。

建议同时补齐：

- Service（ClusterIP）
- ServiceAccount
- Role / RoleBinding

## 文档

- `doc/k8s-resource-api-project-plan.md`：整体方案与里程碑
- `doc/swagger-usage.md`：Swagger 生成与访问说明
- `doc/github-actions-image.md`：镜像自动构建说明

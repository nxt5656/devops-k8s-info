# Kubernetes 资源信息采集 API 项目方案（Golang）

## 1. 项目目标

开发一个基于 Golang 的服务，通过 API 获取 Kubernetes 集群内资源信息，并以统一接口对外提供查询能力，最终以容器方式部署到 Kubernetes 集群中。

核心目标：
- 面向集群资源信息查询（只读场景）。
- 提供稳定、可扩展的 REST API。
- 支持在集群内安全运行（RBAC 最小权限）。
- 便于后续扩展（多集群、缓存、监控、鉴权）。

## 2. MVP 范围（第一阶段）

第一阶段建议先聚焦单集群、核心资源，快速交付可用版本。

功能范围：
- 获取命名空间列表。
- 获取 Pod 列表与详情（按 namespace 过滤）。
- 获取 Deployment 列表与详情。
- 获取 Service 列表与详情。
- 基础健康检查接口（`/healthz`、`/readyz`）。
- Prometheus 指标接口（`/metrics`）。

暂不纳入 MVP：
- 资源变更推送（watch/stream）。
- 多集群聚合查询。
- 复杂权限模型（多租户细粒度 RBAC）。
- Web UI（仅提供 API）。

## 3. 技术选型

- 语言：Go 1.22+
- HTTP 框架：Gin（或标准库 `net/http`，建议 Gin 提升开发效率）
- Kubernetes 客户端：`client-go`
- 配置管理：环境变量 + 配置文件（可选）
- 日志：`zap`（结构化日志）
- 指标：`prometheus/client_golang`
- API 文档：Swagger（`swaggo/swag` + `gin-swagger`）
- 容器化：Docker
- 部署：Kubernetes Deployment + Service

## 4. 系统架构

逻辑分层建议：

1. API 层（`handler`）
- 处理请求参数、返回格式、错误码映射。

2. Service 层（`service`）
- 聚合业务逻辑、参数校验、调用 Kubernetes 访问层。

3. Kubernetes 访问层（`kube`）
- 封装 `client-go` 访问逻辑，屏蔽 API 细节。

4. 模型层（`model`）
- 定义对外响应结构体（避免直接暴露 K8s 原始对象）。

5. 基础设施层（`pkg`）
- 配置、日志、中间件、错误处理、通用工具。

## 5. 建议目录结构

```text
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── kube/
│   ├── model/
│   └── middleware/
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── pkg/
│   ├── config/
│   ├── logger/
│   └── response/
├── deploy/
│   ├── k8s/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── serviceaccount.yaml
│   │   ├── role.yaml
│   │   └── rolebinding.yaml
│   └── Dockerfile
├── doc/
│   └── k8s-resource-api-project-plan.md
├── go.mod
└── go.sum
```

## 6. API 设计（MVP）

建议统一前缀：`/api/v1/{api_base_segment}`（由环境变量 `API_BASE_SEGMENT` 控制，默认 `k8s-info`）

1. 健康检查
- `GET /healthz`：进程存活
- `GET /readyz`：依赖就绪（K8s API 可访问）

2. 命名空间
- `GET /api/v1/{api_base_segment}/namespaces`

3. Pod
- `GET /api/v1/{api_base_segment}/pods?namespace=default`
- `GET /api/v1/{api_base_segment}/pods/{namespace}/{name}`

4. Deployment
- `GET /api/v1/{api_base_segment}/deployments?namespace=default`
- `GET /api/v1/{api_base_segment}/deployments/{namespace}/{name}`

5. Service
- `GET /api/v1/{api_base_segment}/services?namespace=default`
- `GET /api/v1/{api_base_segment}/services/{namespace}/{name}`

6. Swagger 文档
- `GET /swagger/index.html`：Swagger UI
- `GET /swagger/doc.json`：OpenAPI JSON

响应格式建议：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

错误返回建议：
- 参数错误：`400`
- 资源不存在：`404`
- K8s API 访问失败：`502`
- 服务内部错误：`500`

Swagger 规范建议：
- 所有接口必须添加注释（summary、tag、param、success、failure、router）。
- 统一鉴权声明（如后续启用 Bearer Token，在文档中定义 security scheme）。
- 在 CI 中增加文档一致性检查（注释变更后需同步 `swag init` 产物）。

## 7. Kubernetes 访问策略

运行方式：
- 集群内：使用 `InClusterConfig()`
- 本地调试：使用 kubeconfig（`~/.kube/config`）

连接建议：
- 设置合理超时（如 5s~10s）。
- 对 K8s API 失败增加重试（有限次数 + 指数退避）。
- 高频查询资源可加短期缓存（如 5~15 秒）降低 API 压力。

## 8. 安全与权限（RBAC）

原则：
- 最小权限，只授予读取能力（`get/list/watch`）。
- 仅授权必要资源（namespaces/pods/deployments/services）。

示例权限方向：
- `apiGroups: [""]` + `resources: ["pods", "services", "namespaces"]`
- `apiGroups: ["apps"]` + `resources: ["deployments"]`
- `verbs: ["get", "list", "watch"]`

## 9. 可观测性

日志：
- JSON 结构化输出。
- 每个请求带 request_id、耗时、状态码、错误信息。

指标：
- 请求总数、请求耗时、错误率（按接口和状态码维度）。
- K8s API 调用耗时与失败次数。

健康探针：
- livenessProbe -> `/healthz`
- readinessProbe -> `/readyz`

## 10. 部署方案

部署资源：
- Deployment：2 副本起步（高可用）
- Service：ClusterIP（集群内访问）
- ServiceAccount + Role + RoleBinding

资源建议：
- requests: `cpu: 100m`, `memory: 128Mi`
- limits: `cpu: 500m`, `memory: 512Mi`

发布流程建议：
- CI：`go test` + `go vet` + `swag init` + 镜像构建
- CD：推送镜像后更新 Deployment

## 11. 开发里程碑

第 1 周：
- 完成项目骨架、配置管理、日志、健康检查。
- 接入 Kubernetes client，打通命名空间查询接口。

第 2 周：
- 完成 Pod/Deployment/Service 查询接口。
- 补充统一错误处理和响应结构。
- 完成 Swagger 注释与 `/swagger` 路由接入。

第 3 周：
- 补充指标、探针、RBAC、容器化与集群部署清单。
- 完成基础压测与稳定性验证。

第 4 周：
- 回归测试、文档完善、准备上线。

## 12. 风险与应对

1. K8s API 请求压力过高
- 应对：分页、限流、短期缓存、超时与重试。

2. 权限配置不当导致接口失败
- 应对：提供最小可用 RBAC 模板 + 启动时自检。

3. 返回数据结构不稳定
- 应对：内部模型与 K8s 原始对象解耦，版本化 API（`/api/v1/{api_base_segment}`）。

4. 后续需求扩展过快
- 应对：分层架构与模块化设计，预留多集群适配接口。

5. 接口变更后文档不一致
- 应对：将 `swag init` 和文档校验纳入 CI，禁止未更新文档的接口变更合并。

## 13. 下一步建议

在本方案基础上，优先执行以下任务：

1. 初始化项目结构（`cmd/internal/pkg/deploy`）。
2. 先实现 `GET /api/v1/{api_base_segment}/namespaces` 作为端到端样板。
3. 完成 Dockerfile 与最小 RBAC 清单，先在测试集群验证部署。

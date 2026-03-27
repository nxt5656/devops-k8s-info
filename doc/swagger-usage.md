# Swagger 使用说明

## 1. 安装 swag

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

确保 `$GOPATH/bin` 在 `PATH` 中，或直接使用完整路径执行 `swag`。

## 2. 生成 Swagger 文档

在项目根目录执行：

```bash
swag init -g cmd/server/main.go -o docs
```

生成产物位于 `docs/` 目录（如 `docs.go`、`swagger.json`、`swagger.yaml`）。

## 3. 启动后访问

- Swagger UI: `GET /swagger/index.html`
- OpenAPI JSON: `GET /swagger/doc.json`

## 4. 基础路径配置

- 接口基础路径格式：`/api/v1/{segment}`
- `{segment}` 来自环境变量 `API_BASE_SEGMENT`
- 未设置时默认：`k8s-info`

示例：

```bash
export API_BASE_SEGMENT=prod-k8s
```

此时命名空间接口路径为：

```text
/api/v1/prod-k8s/namespaces
```


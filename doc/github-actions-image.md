# GitHub Actions 镜像构建说明

已新增工作流文件：

- `.github/workflows/docker-image.yml`

## 触发条件

- push 到 `main` / `master`
- push tag（如 `v1.0.0`）
- pull request 到 `main` / `master`（仅构建不推送）
- 手动触发（`workflow_dispatch`）

## 镜像地址

- `ghcr.io/nxt5656/devops-k8s-info`

## 标签策略

- 分支标签（branch）
- tag 标签（如 `v1.0.0`）
- commit sha 标签
- 默认分支额外生成 `latest`

## 权限与前置条件

- 工作流权限：
  - `contents: read`
  - `packages: write`
- 使用 `secrets.GITHUB_TOKEN` 推送到 GHCR。

## 推送后验证

1. 打开 GitHub 仓库 `Actions` 页面，确认工作流执行成功。
2. 打开仓库 `Packages`，确认镜像已生成。
3. 在集群中更新 Deployment 镜像 tag 后发布。


# AI Agent Skill

本仓库用于分发面向编码 Agent 的 AI 支持资产。仓库包含可复用的 skills、sub-agents、references、templates，以及用于把这些资产安装到本地 Agent 运行时或项目目录中的 `ferryPilot` 安装工具。

## 仓库内容

- `AISupport/`：可安装的 AI 支持资产包。
- `utils/ferryPilotGo/`：`ferryPilot` CLI 安装器的 Go 实现。
- `specs/`：特性规格、实施计划和任务清单。
- `docs/`：项目结构、维护方式和发布流程的项目级文档。

## 快速开始

可以从 GitHub Release 下载 `ferryPilot` 二进制文件，也可以在本地构建：

```bash
cd utils/ferryPilotGo
go build -o bin/ferryPilot ./cmd/ferryPilot
```

全局安装一个 AISupport 包：

```bash
./bin/ferryPilot -g speckit
```

安装到当前项目：

```bash
./bin/ferryPilot -p speckit
```

为指定 target agent 安装：

```bash
./bin/ferryPilot -g -t codex speckit
./bin/ferryPilot -p -t cursor speckit
```

如果省略 package 名称，并且存在多个可选 package，`ferryPilot` 会提示用户选择。

## 发布构建

当推送匹配 `v*` 的 tag 时，GitHub Actions 会自动构建发布产物。当前发布矩阵包括：

- `windows/amd64`
- `linux/amd64`
- `darwin/amd64`
- `darwin/arm64`

产物会附加到 GitHub Release，例如 `ferryPilot-windows-amd64.tar.gz`。每个归档包都包含可执行文件和 `config/file_map.json`。

## 文档

- [项目概览](docs/project.md)
- [ferryPilot Go 安装器](docs/ferryPilot.md)
- [发布流程](docs/release.md)


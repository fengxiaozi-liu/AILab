# 项目架构

本仓库是 AI 支持资产仓库，核心目标是维护 `AISupport/` 下的 Agent 资产，并通过 `utils/ferryaide/` 提供跨平台安装工具。

## 顶层结构

| 路径 | 职责 |
| --- | --- |
| `AISupport/` | 可分发、可安装的 AI 支持资产包 |
| `utils/ferryaide/` | 安装 `AISupport` 资产的 Go CLI 工具 |
| `docs/` | 项目架构和发布产物说明 |
| `.github/workflows/` | GitHub Release 自动构建流程 |

## AISupport

`AISupport/<package>` 是资产安装单元。每个 package 可以包含 skills、sub-agents、templates、references 等内容。

当前主要 package：

| Package | 内容 |
| --- | --- |
| `speckit` | 结构化规格、澄清、计划、任务、实施、评审、总结工作流 |
| `kratos` | 面向 Kratos 项目的架构、流水线和评审支持资产 |

## ferryaide

`ferryaide` 位于 `utils/ferryaide/`，用于把 `AISupport/<package>` 安装到目标 Agent 的全局目录或当前项目目录。

主要模块：

| 模块 | 职责 |
| --- | --- |
| `cmd/ferryaide` | CLI 入口 |
| `internal/cli` | 解析 `-g`、`-p`、`-t`、`--config` 等参数 |
| `internal/config` | 读取内置默认配置或外部配置 |
| `internal/assets` | 获取 AISupport 数据源，默认从 git clone 到临时目录 |
| `internal/packages` | 发现可安装 package |
| `internal/install` | 解析目标目录并复制安装文件 |
| `internal/transform` | 对需要转换的 sub-agent 文件做格式转换 |

## 主流程

1. 用户运行 `ferryaide -p speckit` 或 `ferryaide -g speckit`。
2. CLI 解析安装模式、target agent 和 package 名称。
3. 配置模块读取编译进二进制的默认 `file_map.json`；如果传入 `--config`，则读取外部配置。
4. 资产模块根据配置拉取 git 数据源到系统临时目录。
5. package 模块发现 `AISupport/` 下可安装的 package；如果用户未指定 package，则在终端中提供上下键选择。
6. install 模块根据 target 映射把 skills、sub-agents 等文件安装到目标目录。
7. 需要转换的 sub-agent markdown 会在安装过程中转换为目标格式。
8. 安装完成后，临时 checkout 被清理。

## 发布流程

推送匹配 `v*` 的 tag 后，GitHub Actions 会为 Windows、Linux、macOS 交叉编译 `ferryaide`，并把单文件可执行产物上传到 GitHub Release。

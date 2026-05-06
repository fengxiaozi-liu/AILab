# 特性规格说明：ferryPilot Go 语言重写

**Feature Branch**: `rewrite-ferrypilot-go`
**Created**: 2026-05-06
**Status**: Ready
**Input**: 用户描述："现在ferrypliot有点问题：1. 使用python语言导致不能交叉编译windows、linux，macos系统；2. 需要重写为go语言，重写为go语言之后就不需要主动发布，建立对应.github/workflow流程就可以。补充澄清：Go 重写版本放在 utils/ferryPilotGo，工具名称改为 ferryPilot，重写过程中不要引用当前 ferryPilot 文件。执行完成之后为这个项目增加 docs 文件夹，描述我们这个项目，增加 AGENTS.md 文件与 README.md 文件。"

---

## 需求描述

### 背景与动机

ferryPilot 当前是一个用于安装仓库内 `AISupport/` 资产的 CLI 工具，现有实现使用 Python 与 PyInstaller 打包。由于当前语言与打包方式不利于面向 Windows、Linux、macOS 的跨平台交叉编译，项目需要将 ferryPilot 重写为 Go 语言，以降低多平台构建与分发成本。

同时，迁移到 Go 后，发布流程应由 GitHub Actions 工作流自动完成，避免维护者在本地手动构建、打包或发布可执行文件。

### 目标用户

- **工具维护者**：负责维护 ferryPilot、配置构建流水线并验证发布产物。
- **AISupport 使用者**：在 Windows、Linux、macOS 上使用 ferryPilot 安装本仓库提供的 AI 支持资产。
- **仓库协作者**：通过 CI/CD 结果确认 ferryPilot 在目标平台上可构建、可发布、可使用。

### 核心需求

- **RQ-001**: ferryPilot 必须以 Go 语言重写，并移除对 Python 运行时、Python 依赖安装和 PyInstaller 打包流程的发布依赖。
- **RQ-002**: Go 版本 ferryPilot 的工具名称必须为 `ferryPilot`，并保留现有 CLI 的核心使用方式，包括全局安装、项目内安装和指定目标 agent 的能力。
- **RQ-003**: Go 版本 ferryPilot 必须继续扫描仓库 `AISupport/` 下的一级包目录，并允许用户选择待安装的包。
- **RQ-004**: Go 版本 ferryPilot 必须继续按照既有文件映射规则安装 package 内的可安装内容。
- **RQ-005**: Go 版本 ferryPilot 必须继续支持 Codex `sub-agents/*.md` 在安装时转换为 `.toml` 的既有行为。
- **RQ-006**: 项目必须提供 GitHub Actions 工作流，用于构建 Windows amd64、Linux amd64、macOS amd64、macOS arm64 平台的 ferryPilot 可执行产物。
- **RQ-007**: GitHub Actions 工作流必须在创建 tag 时自动构建并将产物附加到 GitHub Release，使维护者无需在本地手动发布 ferryPilot 可执行文件。
- **RQ-008**: ferryPilot 的使用文档必须更新为 Go 重写后的构建、安装和自动化发布说明，不再把 Python/PyInstaller 作为主要发布路径。
- **RQ-009**: Go 重写版本必须位于 `utils/ferryPilotGo`，并作为独立实现交付。
- **RQ-012**: Go 版本 ferryPilot 必须通过 `file_map.json` 描述 git 数据源、默认 target 和安装映射；运行时应将 git 数据源下载到临时目录，再从临时目录复制或转换到目标目录。
- **RQ-010**: Go 重写完成后，项目必须新增 `docs/` 目录，用于描述项目定位、结构、核心资产、ferryPilot 工具和后续维护信息。
- **RQ-011**: Go 重写完成后，项目根目录必须新增或更新 `AGENTS.md` 和 `README.md`，分别面向 AI agent 协作规范与项目使用者说明项目。

### 关键实体（若涉及数据）

- **ferryPilot CLI**: 面向用户的命令行安装工具，提供全局安装、项目安装、目标 agent 选择等入口。
- **AISupport Package**: `AISupport/` 下的一级目录，例如 `speckit`、`kratos`，是用户可选择的安装单元。
- **File Map**: 描述 package 内容如何复制到目标位置的映射规则。
- **Git Data Source**: `file_map.json` 中描述的远程 git 仓库，作为 AISupport 资产来源。
- **Install Target**: 安装目标位置，包括用户级全局目录、当前项目目录，以及指定 agent 对应的目录结构。
- **Release Artifact**: GitHub Actions 构建出的跨平台可执行产物。
- **Project Documentation**: 项目级说明文档集合，包括 `docs/`、根目录 `AGENTS.md` 和根目录 `README.md`。

### 约束与边界

- Go 重写后的 ferryPilot 需要覆盖 Windows amd64、Linux amd64、macOS amd64、macOS arm64。
- Go 重写版本的实现目录为 `utils/ferryPilotGo`，产物与文档中的工具名称统一为 `ferryPilot`。
- Go 重写版本不得通过 Go embed 将 AISupport 文件集成进二进制；必须按 `file_map.json` 配置下载 git 数据源到临时目录，再执行安装。
- 迁移后不应要求用户或维护者安装 Python 依赖来完成 ferryPilot 的常规构建、打包或发布。
- 自动化流程应位于仓库 `.github/workflows/` 下，并能在创建 tag 时通过 GitHub Actions 生成 GitHub Release 产物。
- 重写范围聚焦 ferryPilot，不改变 `AISupport/` 下 package 的业务内容与目录语义。
- 重写过程中不得依赖、导入、复制或运行当前 `utils/ferryPilot` 下的 Python 实现文件；现有 README 和用户可见行为只能作为需求参考。
- 除非后续澄清另有要求，Go 版本应以现有 README 描述的用户可见行为作为兼容基线。
- 项目级文档应描述当前仓库整体，而不仅是 ferryPilot 单个工具；但不应改变 `AISupport/` 资产本身的业务语义。

---

## 待澄清问题

### CQ-001: 目标平台架构范围

> **类别**: 约束与权衡

**Q (提问)**:
Windows、Linux、macOS 需要覆盖哪些 CPU 架构？该范围会影响 GitHub Actions 构建矩阵和产物命名。

*参考选项*:
- A. 仅覆盖 amd64/x86_64
- B. 覆盖 amd64 和 arm64
- C. Windows/Linux 覆盖 amd64，macOS 覆盖 amd64 和 arm64

**A (澄清结论)**:
选择 C。Windows/Linux 覆盖 amd64，macOS 覆盖 amd64 和 arm64；产物命名与 Go 的 `GOOS/GOARCH` 保持一致，例如 `windows-amd64`、`linux-amd64`、`darwin-amd64`、`darwin-arm64`。

### CQ-002: 自动发布触发方式

> **类别**: 完成信号

**Q (提问)**:
GitHub Actions 应在什么事件下生成可下载产物？这会影响“无需主动发布”的验收标准。

*参考选项*:
- A. 每次 push 到主分支都构建并上传 workflow artifacts
- B. 创建 tag 时构建，并附加到 GitHub Release
- C. pull request 只做构建验证，tag 时生成发布产物

**A (澄清结论)**:
选择 B。创建 tag 时触发构建，并将跨平台可执行产物附加到 GitHub Release。

### CQ-003: Go 重写目录与工具名称

> **类别**: 约束与权衡

**Q (提问)**:
Go 重写版本应放置在哪个目录，最终工具名称是什么，以及是否允许引用当前 ferryPilot 文件？

**A (澄清结论)**:
Go 重写版本放在 `utils/ferryPilotGo`，最终工具名称统一为 `ferryPilot`。重写过程中不得依赖、导入、复制或运行当前 `utils/ferryPilot` 下的 Python 实现文件；现有 README 和用户可见行为仅作为需求参考。

### CQ-004: 项目级文档交付

> **类别**: 完成信号

**Q (提问)**:
Go 重写完成后，项目是否需要补充项目级文档资产？

**A (澄清结论)**:
需要。完成后必须新增 `docs/` 目录描述项目定位、结构、核心资产、ferryPilot 工具和维护信息；根目录必须新增或更新 `README.md` 面向项目使用者说明项目，新增或更新 `AGENTS.md` 面向 AI agent 说明协作规范。

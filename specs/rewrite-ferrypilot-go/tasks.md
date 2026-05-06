# 任务清单：ferryPilot Go 语言重写

**Spec**: `specs/rewrite-ferrypilot-go/spec.md`
**Plan**: `specs/rewrite-ferrypilot-go/plan.md`
**Created**: 2026-05-06
**项目识别结果**: CLI 工具 + AI 支持资产仓库
**技术链依赖**: Phase 1 使用 Go module 与 CLI 工程骨架知识；Phase 2 使用 Go 标准库 CLI、路径、文件系统与测试知识；Phase 3 使用 `file_map.json` 配置、git 数据源下载、安装映射、Markdown front matter 转 TOML 与临时目录集成测试知识；Phase 4 使用 GitHub Actions、Go 交叉编译与 GitHub Release 发布知识；Phase 5 使用项目文档、agent 协作规范与端到端验收知识。

---

## 格式说明

- `[TaskID]`：如 `T001`，全局递增唯一编号
- `[技术域]`：该任务所属模块或架构层，如 `[DB]` `[Service]` `[API]` `[UI]` `[Config]`
- `[RQ-xxx]`：对应需求编号，无法对应时使用临时功能点编号
- 每条任务必须包含精确文件路径或明确产物

---

<!-- 1:1 对应执行拓扑中的真实节点结构 -->

## Phase 1：Go 独立工程骨架

- [✅️] T001 [Go/Scaffold] [RQ-001,RQ-009] 创建独立 Go module 与基础目录结构，不引用现有 Python 实现 `utils/ferryPilotGo/go.mod`, `utils/ferryPilotGo/cmd/ferryPilot/`, `utils/ferryPilotGo/internal/`
- [✅️] T002 [Go/CLI] [RQ-002] 建立 CLI 入口骨架，输出工具名统一为 `ferryPilot`，预留全局/项目/target 参数解析入口 `utils/ferryPilotGo/cmd/ferryPilot/main.go`
- [✅️] T003 [Go/App] [RQ-002,RQ-003,RQ-004,RQ-005] 定义应用层执行模型、选项对象、安装结果对象与错误出口约定 `utils/ferryPilotGo/internal/app/`
- [✅️] T004 [Test] [RQ-001,RQ-002] 添加基础 Go 测试入口与参数解析测试框架，确保 `go test ./...` 可运行 `utils/ferryPilotGo/**/*_test.go`

## Phase 2：CLI 参数、配置与目标路径模型

- [✅️] T005 [CLI] [RQ-002] 实现 `-g/--global`、`-p/--project`、`-t/--target` 参数解析和互斥校验 `utils/ferryPilotGo/internal/cli/`
- [✅️] T006 [Config] [RQ-004,RQ-012] 通过 Go 版本 `file_map.json` 定义 git 数据源、安装映射规则和 target agent 规则，不读取或复用 `utils/ferryPilot/src/config/file_map.json` `utils/ferryPilotGo/config/file_map.json`, `utils/ferryPilotGo/internal/config/`
- [✅️] T007 [Install] [RQ-002,RQ-004] 实现 global/project 安装目标解析，覆盖用户 home、当前工作目录和 target agent 目录规则 `utils/ferryPilotGo/internal/install/target.go`
- [✅️] T008 [Test] [RQ-002,RQ-004] 为参数冲突、缺失模式、未知 target、路径解析增加表驱动单元测试 `utils/ferryPilotGo/internal/cli/`, `utils/ferryPilotGo/internal/install/`

## Phase 3：AISupport 资产发现、安装与转换

- [✅️] T009 [Assets] [RQ-003,RQ-004,RQ-012] 实现资产源抽象，支持按 `file_map.json` 将 git 数据源下载到 tmp，并从 tmp 中的 `AISupport/` 读取 `utils/ferryPilotGo/internal/assets/`
- [✅️] T010 [Assets] [RQ-003] 实现仅扫描 `AISupport/` 一级目录的 package discovery，过滤非目录和空目录 `utils/ferryPilotGo/internal/packages/`
- [✅️] T011 [Install] [RQ-004] 实现基于 Go 内置映射规则的目录创建、文件复制、覆盖与安装摘要 `utils/ferryPilotGo/internal/install/`
- [✅️] T012 [Transform] [RQ-005] 实现 Codex `sub-agents/*.md` front matter 与正文到 `.toml` 的转换逻辑 `utils/ferryPilotGo/internal/transform/`
- [✅️] T013 [App] [RQ-002,RQ-003,RQ-004,RQ-005] 串联 package 选择、目标解析、复制安装、sub-agent 转换和结果汇总 `utils/ferryPilotGo/internal/app/`
- [✅️] T014 [Test] [RQ-003,RQ-004,RQ-005] 使用测试文件系统和临时目录覆盖 package 发现、文件复制、sub-agent 转换和安装摘要 `utils/ferryPilotGo/internal/assets/`, `utils/ferryPilotGo/internal/packages/`, `utils/ferryPilotGo/internal/install/`, `utils/ferryPilotGo/internal/transform/`

## Phase 4：构建、嵌入资产与 GitHub Release

- [✅️] T015 [Build] [RQ-001,RQ-006] 增加本地构建脚本或 Makefile 目标，生成名为 `ferryPilot` 的可执行文件 `utils/ferryPilotGo/Makefile`
- [✅️] T016 [Config/Build] [RQ-003,RQ-004,RQ-006,RQ-012] 增加发布配置打包机制，将 `config/file_map.json` 与可执行文件一起归档，运行时由配置下载 git 数据源到 tmp `utils/ferryPilotGo/config/file_map.json`, `.github/workflows/release.yml`
- [✅️] T017 [CI] [RQ-006,RQ-007] 新增 tag 触发 GitHub Actions release workflow，覆盖 `windows/amd64`、`linux/amd64`、`darwin/amd64`、`darwin/arm64` `.github/workflows/release.yml`
- [✅️] T018 [CI] [RQ-007] 在 release workflow 中配置 Release 创建/更新、artifact 命名、Windows `.exe` 后缀和 `contents: write` 权限 `.github/workflows/release.yml`
- [✅️] T019 [Test/CI] [RQ-001,RQ-006] 增加 CI 构建前测试步骤，执行 `go test ./...` 并验证多平台 `go build` `.github/workflows/release.yml`

## Phase 5：项目级文档与用户文档

- [✅️] T020 [Docs] [RQ-008] 更新 Go 版本 ferryPilot 使用、构建、安装和发布说明，移除 Python/PyInstaller 作为主要发布路径的描述 `utils/ferryPilotGo/README.md`
- [✅️] T021 [Docs] [RQ-010] 新增项目级 `docs/`，描述仓库定位、目录结构、AISupport 资产、ferryPilot 工具和维护流程 `docs/`
- [✅️] T022 [Docs] [RQ-011] 新增根目录 `README.md`，面向项目使用者说明项目用途、快速开始、目录结构和 release 获取方式 `README.md`
- [✅️] T023 [Docs/Agent] [RQ-011] 新增根目录 `AGENTS.md`，面向 AI agent 说明协作规则、事实源、禁止引用旧 Python 实现和 speckit 工作流入口 `AGENTS.md`
- [✅️] T024 [Docs] [RQ-008,RQ-010,RQ-011] 对齐所有文档中的工具名、目录名、构建矩阵和发布触发方式 `README.md`, `AGENTS.md`, `docs/`, `utils/ferryPilotGo/README.md`

## Phase 6：验收与防回归检查

- [✅️] T025 [Verify] [RQ-001,RQ-009] 检查 Go 实现、测试、构建脚本和 workflow 未依赖、导入、复制或运行 `utils/ferryPilot` 下的 Python 文件 `utils/ferryPilotGo/`, `.github/workflows/release.yml`
- [✅️] T026 [Verify] [RQ-002,RQ-003,RQ-004,RQ-005] 执行本地端到端安装验证，覆盖 `ferryPilot -g`、`ferryPilot -p`、`-t codex` 与 sub-agent `.toml` 产物 `utils/ferryPilotGo/`
- [✅️] T027 [Verify] [RQ-001,RQ-006] 执行 `go test ./...` 与本地 `go build`，确认产物名为 `ferryPilot` `utils/ferryPilotGo/`
- [✅️] T028 [Verify] [RQ-006,RQ-007] 静态检查 GitHub Actions workflow 的 tag 触发、矩阵、权限和 release artifact 命名 `.github/workflows/release.yml`
- [✅️] T029 [Verify] [RQ-001,RQ-002,RQ-003,RQ-004,RQ-005,RQ-006,RQ-007,RQ-008,RQ-009,RQ-010,RQ-011,RQ-012] 建立需求覆盖检查，逐项确认 RQ-001 至 RQ-012 均有实现、测试或文档证据 `specs/rewrite-ferrypilot-go/tasks.md`

## 需求覆盖证据

| RQ | 实现 / 验证证据 |
|----|----------------|
| RQ-001 | `utils/ferryPilotGo/go.mod`、`cmd/ferryPilot`、`go test ./...`、本地 `go build` |
| RQ-002 | `internal/cli`、`internal/app`、端到端 `-g/-p/-t codex speckit` 验证 |
| RQ-003 | `internal/assets`、`internal/packages`、package discovery 单元测试 |
| RQ-004 | `config/file_map.json`、`internal/config`、`internal/install`、安装复制集成测试与端到端验证 |
| RQ-005 | `internal/transform`、sub-agent `.toml` 转换测试与端到端产物验证 |
| RQ-006 | `.github/workflows/release.yml` 构建矩阵、本地 build 验证 |
| RQ-007 | `.github/workflows/release.yml` tag 触发、`contents: write`、Release artifact 上传配置 |
| RQ-008 | `utils/ferryPilotGo/README.md` |
| RQ-009 | `utils/ferryPilotGo/` 独立实现与旧实现精确扫描无依赖命中 |
| RQ-010 | `docs/project.md`、`docs/ferryPilot.md`、`docs/release.md` |
| RQ-011 | `README.md`、`AGENTS.md` |
| RQ-012 | `utils/ferryPilotGo/config/file_map.json`、`internal/assets` git clone 到 tmp、端到端 `--config` 数据源验证 |

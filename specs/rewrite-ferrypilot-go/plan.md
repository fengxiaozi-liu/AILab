# 实施计划：ferryPilot Go 语言重写

**Feature Branch**: `rewrite-ferrypilot-go`
**Created**: 2026-05-06
**Status**: Ready
**Spec**: `specs/rewrite-ferrypilot-go/spec.md`

---

## 摘要

本计划将 ferryPilot 从 Python/PyInstaller 发布路径迁移为独立 Go CLI，并在 `utils/ferryPilotGo` 下交付。新工具最终名称统一为 `ferryPilot`，保留现有用户可见安装行为：扫描 `AISupport/` 一级包、选择 package、按映射规则安装 skills/sub-agents 等资产，并继续将 Codex `sub-agents/*.md` 转换为 `.toml`。

交付重点包括：Go CLI 模块化实现、`file_map.json` 驱动的 git 数据源下载与安装映射、跨平台构建、tag 触发的 GitHub Release 发布流程，以及项目根文档 `README.md`、`AGENTS.md` 与 `docs/`。实施过程中不得依赖、导入、复制或运行当前 `utils/ferryPilot` 下的 Python 实现文件；现有 README 与用户可见行为仅作为需求参考。

## 项目适配输入

- 项目类型: CLI 工具 + AI 支持资产仓库
- 结构约束: Go 重写版本固定在 `utils/ferryPilotGo`；自动化流程位于 `.github/workflows/`；项目级文档位于根目录与 `docs/`
- 交付物类型: 代码 / 配置 / 文档 / 发布产物
- 特殊验证要求: Go 单元测试、跨平台构建矩阵验证、tag release workflow 配置校验、文档存在性与内容检查

---

## Phase 0：调研与关键决策

### 实体与关系分析

| 实体 | 关键属性 | 与其他实体关系 | 来源 RQ |
|------|---------|---------------|---------|
| ferryPilot CLI | 命令名 `ferryPilot`、安装模式、target agent、交互选择 | 调用 Package Discovery、Installer、Transformer 完成安装 | RQ-001, RQ-002 |
| AISupport Package | package 名称、skills、sub-agents、README 等资产 | 从 git 数据源下载到 tmp 后，作为 CLI 的安装源，被复制到 Install Target | RQ-003, RQ-004, RQ-012 |
| File Map | git 数据源、默认 target、源路径、目标路径、是否需要转换 | 指导 Data Source 下载和 Installer 写入目标目录 | RQ-004, RQ-005, RQ-012 |
| Install Target | global/project 模式、agent 类型、目标根目录 | 接收 Installer 复制或转换后的资产 | RQ-002, RQ-004 |
| Release Artifact | GOOS、GOARCH、文件名、校验信息 | 由 GitHub Actions 构建并附加到 Release | RQ-006, RQ-007 |
| Project Documentation | `README.md`、`AGENTS.md`、`docs/` | 说明项目定位、协作方式、工具维护方式 | RQ-008, RQ-010, RQ-011 |

### 状态流转（若适用）

| 实体 | 状态列表 | 流转规则 | 来源 RQ |
|------|---------|---------|---------|
| 安装执行 | 初始化 -> 发现 package -> 用户选择 -> 解析目标 -> 复制/转换 -> 汇总结果 -> 失败退出或成功退出 | 参数校验通过后进入发现；无 package 或目标不可写时失败退出；每个文件操作记录结果 | RQ-002, RQ-003, RQ-004, RQ-005 |
| Release Artifact | 未构建 -> 已构建 -> 已归档 -> 已附加 Release | 创建 tag 触发 workflow；矩阵构建全部完成后上传到 GitHub Release | RQ-006, RQ-007 |
| Project Documentation | 缺失 -> 草案完成 -> 与实现同步 -> 可维护 | Go 重写完成后补齐项目级文档，并与工具行为和目录结构保持一致 | RQ-008, RQ-010, RQ-011 |

### 依赖调研

| 依赖对象 | 用途 | 已有能力 | 需要新增或修改 |
|----------|------|----------|----------------|
| Go toolchain | 实现 CLI、测试、交叉编译 | 仓库当前无 Go 模块 | 在 `utils/ferryPilotGo` 初始化 Go module，并建立标准包结构 |
| Go 标准库 | 参数解析、文件系统、路径处理、JSON 配置、临时目录、测试 | 可覆盖主要需求 | 优先使用 `flag`、`os`、`path/filepath`、`encoding/json`、`testing`，减少外部依赖 |
| Git CLI | 下载远程 AISupport 数据源到 tmp | 系统需具备 git | 运行时根据 `file_map.json` 执行浅 clone，完成后清理临时目录 |
| AISupport 目录 | 安装源资产 | 已包含 `speckit`、`kratos` 等 package | 从 git 数据源下载后读取 `AISupport/` |
| GitHub Actions | 跨平台构建与发布 | 仓库当前无 `.github/` | 新增 tag 触发 release workflow，覆盖 windows/linux/darwin 矩阵 |
| 项目文档 | 项目说明与协作规范 | 根目录当前无 `README.md`、`AGENTS.md`、`docs/` | 新增项目级文档，覆盖仓库定位、目录、安装、维护与 agent 协作规则 |

### 调研结论

- **DR-001**: Go 版本作为独立实现放在 `utils/ferryPilotGo`
  - 决策: 初始化独立 Go module，不引用当前 `utils/ferryPilot` 下的 Python 文件
  - 理由: 满足独立重写与跨平台交叉编译目标，也避免新实现被旧实现结构牵制
  - 排除方案: 在 `utils/ferryPilot` 原目录内替换实现；从 Python 文件迁移逻辑或复用配置

- **DR-002**: CLI 优先使用 Go 标准库实现
  - 决策: 参数解析、路径处理、文件复制、嵌入资产和测试优先使用标准库
  - 理由: CLI 功能边界清晰，标准库足够覆盖核心需求，减少跨平台发布风险
  - 排除方案: 引入大型 CLI 框架作为第一版基础依赖

- **DR-003**: 资产源采用 `file_map.json` 描述的 git 数据源
  - 决策: 运行时读取 `file_map.json`，将配置的 git 仓库浅 clone 到 tmp，再从 `tmp/AISupport/<package>` 安装到目标目录
  - 理由: AISupport 文件不进入二进制，数据来源、映射规则和 target 行为都由配置文件表达
  - 排除方案: 使用 Go embed 将 AISupport 资产打包进二进制；隐式扫描本地仓库根作为默认发布路径

- **DR-004**: GitHub Release 由 tag 触发
  - 决策: 创建 tag 时构建 `windows-amd64`、`linux-amd64`、`darwin-amd64`、`darwin-arm64` 并附加到 GitHub Release
  - 理由: 与 CQ-002 对齐，减少维护者本地手动发布操作
  - 排除方案: 每次 push 都发布产物；只上传 workflow artifacts 不创建 Release

---

## Phase 1：技术设计

### 模块设计

| 模块 | 职责 | 主要输入 | 主要输出 | 来源 RQ |
|------|------|---------|---------|---------|
| `cmd/ferryPilot` | CLI 入口、参数解析、执行编排 | `-g/--global`、`-p/--project`、`-t/--target`、用户交互输入 | 退出码、安装摘要、错误信息 | RQ-002 |
| `internal/app` | 应用服务层，串联发现、选择、安装、汇总 | CLI options、工作目录、home 目录 | 安装结果对象 | RQ-002, RQ-003 |
| `internal/packages` | 发现 AISupport 一级 package，过滤无效目录 | 本地或嵌入式 asset filesystem | package 列表 | RQ-003 |
| `internal/assets` | 按 `file_map.json` 下载 git 数据源到 tmp，并暴露可读资产文件系统 | git repository、ref、临时目录 | 可读资产文件系统与清理动作 | RQ-003, RQ-004, RQ-012 |
| `internal/install` | 目标路径解析、目录创建、复制文件、覆盖策略 | package、file map、install target | 文件安装结果 | RQ-004 |
| `internal/transform` | 将 Codex `sub-agents/*.md` 转换为 `.toml` | markdown front matter 与正文 | `.toml` 内容 | RQ-005 |
| `internal/config` | 读取并校验 `file_map.json` 数据源和 target 映射规则 | config path、工作目录 | 数据源和映射计划 | RQ-004, RQ-012 |
| `.github/workflows/release.yml` | tag 触发跨平台构建与 Release 上传 | tag、Go module、AISupport 资产 | Release artifacts | RQ-006, RQ-007 |
| `docs/`、`README.md`、`AGENTS.md` | 项目级说明与协作规范 | 实现事实、spec、plan | 文档交付物 | RQ-008, RQ-010, RQ-011 |

### 协议 / 接口设计（若适用）

| 服务或模块 | 方法或动作 | 请求关键字段 | 响应关键字段 | 来源 RQ |
|------------|-----------|-------------|-------------|---------|
| CLI | `ferryPilot -g` | install mode = global, optional target agent | package 选择提示、安装摘要 | RQ-002 |
| CLI | `ferryPilot -p` | install mode = project, optional target agent | package 选择提示、安装摘要 | RQ-002 |
| CLI | `ferryPilot -g -t codex` | install mode、target agent | 指定 agent 的全局安装结果 | RQ-002 |
| CLI | `ferryPilot -p -t cursor` | install mode、target agent | 指定 agent 的项目安装结果 | RQ-002 |
| Release workflow | tag push | tag ref、GOOS、GOARCH | `ferryPilot-<goos>-<goarch>` 产物与 Windows `.exe` | RQ-006, RQ-007 |

### 数据或配置设计（若适用）

| 对象 | 字段或配置项 | 类型 | 约束 | 说明 |
|------|--------------|------|------|------|
| CLI Options | `mode` | enum | `global` 或 `project` 二选一 | 对应 `-g`、`-p` |
| CLI Options | `target` | string | 可选；未知 target 应给出明确错误或可选列表 | 对应 `-t/--target` |
| Package | `name` | string | 来自 `AISupport/` 一级目录 | 用于用户选择和安装摘要 |
| File Map | `data_source` | object | 必须描述 git repository，可选 ref | 运行时 clone 到 tmp |
| File Mapping | `source` | string | 仅描述 AISupport package 内相对路径 | 不引用旧 Python 配置文件 |
| File Mapping | `destination` | string | 根据 install mode 和 target agent 解析 | 输出到用户目录或项目目录 |
| SubAgent Front Matter | `name`、`description` | string | 从 markdown YAML front matter 读取 | 转换为 TOML 字段 |
| Release Matrix | `goos`、`goarch` | string | `windows/amd64`、`linux/amd64`、`darwin/amd64`、`darwin/arm64` | 与 CQ-001 对齐 |

### 错误 / 枚举 / 文案（若适用）

| 类型 | 名称 | 说明 | 来源 RQ |
|------|------|------|---------|
| Enum | `InstallModeGlobal` | 安装到当前用户 home 下的目标运行时目录 | RQ-002 |
| Enum | `InstallModeProject` | 安装到当前工作目录对应的项目目标目录 | RQ-002 |
| Error | `ErrNoAISupportPackages` | 未发现可安装 package | RQ-003 |
| Error | `ErrInvalidInstallMode` | 未指定或同时指定 `-g`、`-p` | RQ-002 |
| Error | `ErrUnknownTarget` | 指定 target agent 不受支持 | RQ-002, RQ-004 |
| Error | `ErrAssetRead` | 读取 AISupport 资产失败 | RQ-003, RQ-004 |
| Error | `ErrInstallWrite` | 创建目录或写入文件失败 | RQ-004 |
| Error | `ErrSubAgentTransform` | sub-agent markdown 转 TOML 失败 | RQ-005 |

---

## 测试与验证

### 测试目标

| 目标 | 层次 | 类型 | 优先级 | 来源 RQ | 说明 |
|------|------|------|--------|---------|------|
| CLI 参数解析 | CLI | Unit | P0 | RQ-002 | 覆盖 `-g`、`-p`、`-t`、冲突参数和缺失模式 |
| Package 发现 | Core | Unit | P0 | RQ-003 | 使用测试 FS 验证只扫描 AISupport 一级目录 |
| 安装映射 | Core | Unit | P0 | RQ-004 | 验证不同 mode/target 下目标路径解析正确 |
| 文件复制 | Core | Integration | P0 | RQ-004 | 在临时目录执行安装，验证目录和文件产物 |
| sub-agent 转换 | Core | Unit | P0 | RQ-005 | 使用真实格式样例验证 markdown front matter 到 TOML |
| git 数据源下载 | Core | Integration | P1 | RQ-003, RQ-004, RQ-012 | 验证根据 `file_map.json` 获取数据源并从 tmp 安装 |
| 跨平台构建 | Release | CI | P0 | RQ-006 | workflow 构建四个平台矩阵 |
| GitHub Release 上传 | Release | CI | P0 | RQ-007 | tag 触发后产物附加到 Release |
| 项目级文档 | Docs | Review | P1 | RQ-008, RQ-010, RQ-011 | 验证 `docs/`、`README.md`、`AGENTS.md` 与实现一致 |

### 验证要求

- `go test ./...` 在 `utils/ferryPilotGo` 通过。
- `go build` 可在本机生成名为 `ferryPilot` 的可执行文件。
- GitHub Actions workflow 包含 `windows/amd64`、`linux/amd64`、`darwin/amd64`、`darwin/arm64` 构建矩阵，并将 `config/file_map.json` 与可执行文件一起打包。
- 创建 tag 时 workflow 能生成并上传 Release artifacts。
- 仓库根存在 `README.md`、`AGENTS.md` 和 `docs/`，且内容描述项目整体而不仅是单个工具。
- 实现与 workflow 不依赖、导入、复制或运行 `utils/ferryPilot` 下的 Python 实现文件。

---

## 风险与边界

| 风险项 | 影响 | 缓解措施 |
|--------|------|---------|
| 运行环境缺少 git | 无法下载 AISupport 数据源 | 在文档中声明运行时依赖 git；错误信息明确指向数据源 clone 失败 |
| 不允许引用旧 Python 实现导致行为细节缺口 | 可能遗漏旧工具隐含行为 | 以 spec、当前 README 和 AISupport 资产结构作为事实源；缺口通过测试和用户验收补齐 |
| 不同操作系统路径和权限差异 | 全局安装路径、可执行权限、换行符处理可能不一致 | 使用 Go 标准路径 API；针对 Windows/Linux/macOS 在 CI 做构建，关键路径逻辑用表驱动测试覆盖 |
| sub-agent markdown 转 TOML 细节不严谨 | Codex agent 产物不可用 | 为 front matter、正文保留、特殊字符转义建立单元测试 |
| Release workflow 权限不足 | tag 构建成功但无法创建或更新 Release | workflow 显式配置 `contents: write`，使用成熟 action 上传 Release 资产 |
| 项目级文档与实现不同步 | 使用者和 agent 获取错误指引 | 文档作为完成验收项，最终 review 时对照 CLI 行为、目录结构和 workflow |

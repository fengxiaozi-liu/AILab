# AI 工程体系

这是一个安装 AI 能力包的仓库，不是业务服务仓库。

仓库维护的是一套可分发的 `prompt + skill + rule + reference + template` 资产，以及用于生成项目说明文件的 agent 模板源。它的职责不是承载某个业务系统，而是把一套标准化研发流程安装到 Codex、Gemini、Copilot 等运行时。

## 核心流程

```text
/specify -> /clarify -> /plan -> /tasks -> /analyze -> /implement -> /code-review -> /summary
                    \-> /checklist
```

## 仓库定位

- `speckit-*` 负责标准化研发流程。
- 项目适配 skill 负责识别并补充具体项目约束。
- 语言 skill 负责补充语言级实现约束。
- 安装脚本负责把这些资产分发到不同 AI 运行时，并生成项目级说明文件。

这意味着：

- 流程 skill 不应静态绑定具体项目 skill。
- 是否加载某个项目 skill，应由 agent 在运行时根据当前项目决定。
- `.shared/agents/*.md` 只用于生成项目说明文件，不是运行时共享 `agents/` 目录。

## 资产边界

### 源资产

下面这些目录是维护入口，应直接在仓库内修改：

```text
.shared/
├── agents/                          # 项目说明文件模板源
├── skills/                          # 技能源文件
└── ...

.codex/prompts/                      # Codex 命令入口源
.gemini/.agents/workflows/           # Gemini 工作流入口源
.github/prompts/                     # GitHub/Copilot prompt 入口源
.github/agents/                      # GitHub/Copilot agent 入口源
scripts/link_ai.py                   # 安装脚本
```

### 运行时产物

下面这些文件或目录是安装脚本生成或覆盖的运行时产物：

- `~/.codex/*`
- `~/.gemini/antigravity/*`
- 项目内 `AGENTS.md`
- 项目内 `GEMINI.md`
- 项目内 `.github/copilot-instructions.md`
- 项目内 `.agent/*`
- 项目内 `.github/*` 中由安装脚本管理的部分

默认原则：

- 修改行为规范、模板、skill、reference 时，改源资产。
- 运行时目录只作为安装目标，不作为主要维护入口。
- 项目说明文件属于安装产物，可重新生成覆盖。

## 目录结构

```text
.shared/
├── agents/                          # 项目说明文件模板源
├── skills/                          # 技能源文件
│   ├── speckit-*/                   #   标准化流程 skill
│   ├── kratos-patterns/             #   Kratos 项目适配器
│   ├── kratos-domain/               #   聚合、分层、命名、Repo、UseCase、测试
│   ├── kratos-components/           #   Ent、EventBus、Crontab、Depend、Config
│   ├── kratos-pkg/                  #   internal/pkg 通用基础能力
│   ├── kratos-conventions/          #   Error、Enum、i18n
│   ├── kratos-entry/                #   Server、Gateway、Proto、Wire、Codegen
│   ├── golang-patterns/             #   Go 语言模式
│   └── project-skill-generator/     #   从项目抽取并生成项目 skill
└── ...
```

## prompt / skill / agent 的关系

### prompt

- prompt 是运行时命令入口。
- 它决定用户触发某个命令时，如何进入对应流程。
- 例如 Codex 的 `/plan`、`/tasks`、`/implement` 入口来自 `.codex/prompts/*`。

### skill

- skill 是可复用能力单元。
- 一个 skill 通常包含 `SKILL.md`，以及可选的 `rules/`、`reference/`、`templates/`。
- 进入 skill 后，优先读取 `rules`，再按需读取 `reference`。

### agent

- agent 在这里主要指“项目级说明文件模板源”。
- `.shared/agents/*.md` 不会被直接复制到运行时目录。
- 安装脚本会基于这些模板生成项目内的说明文件，例如 `AGENTS.md`、`GEMINI.md`、`.github/copilot-instructions.md`。

### 三者的协作关系

1. 用户通过 prompt 触发一个流程。
2. agent 负责当前项目的总体协作约束与入口说明。
3. agent 再按任务类型加载对应 skill。
4. 若当前项目可识别，再由 agent 决定是否补充加载项目 skill 或语言 skill。

## Commands

| 命令 | 说明 |
|------|------|
| `/specify` | 需求捕获 |
| `/clarify` | 深度澄清 |
| `/checklist` | 需求质量检查 |
| `/plan` | 架构设计 |
| `/tasks` | 任务拆解 |
| `/analyze` | 一致性校验 |
| `/implement` | 代码实施 |
| `/code-review` | 代码评审 |
| `/summary` | 完工总结 |

## Skills

| 类型 | Skill | 说明 |
|------|-------|------|
| 流程型 | `speckit-*` | 标准化需求、计划、任务、分析、实施、评审、总结流程 |
| 项目适配器 | `kratos-patterns` | Kratos 项目适配器，负责路由到对应 Kratos 技能 |
| 领域型 | `kratos-domain` | 聚合、分层、命名、Repo、UseCase、测试 |
| 组件型 | `kratos-components` | Ent、EventBus、Crontab、Depend、Config |
| 基础包 | `kratos-pkg` | internal/pkg 通用能力 |
| 规范型 | `kratos-conventions` | Error、Enum、i18n |
| 接入型 | `kratos-entry` | Server、Gateway、Proto、Wire、Codegen |
| 语言型 | `golang-patterns` | Go 语言模式与惯用法 |
| 元能力 | `project-skill-generator` | 从任意项目抽取并生成项目专属 skill |

## Rules 与 Reference

- `skills/*/rules/*.md`：治理规则、硬约束、禁止项、检查项。
- `skills/*/reference/*.md`：实现说明、模板、示例、常见模式与常见坑。
- 推荐顺序：先进入对应 skill，再先读 `rules`，后读 `reference`。

## Speckit 与项目适配器

- `speckit-*` 负责标准化研发流程，本身不绑定具体语言或框架。
- 遇到具体项目时，再由 agent 在运行时决定是否加载对应项目适配器 skill。
- 当前仓库内置的是 Kratos 适配器：`kratos-patterns`。
- 后续接入其他语言或技术栈时，应新增新的项目适配器 skill，而不是继续把 speckit 流程写死到 Kratos 规则里。

## Templates

Speckit 产物模板跟随各自 skill 维护：

| 模板 | 路径 |
|------|------|
| spec | `skills/speckit-specify/templates/spec-template.md` |
| plan | `skills/speckit-plan/templates/plan-template.md` |
| tasks | `skills/speckit-tasks/templates/tasks-template.md` |
| checklist | `skills/speckit-checklist/templates/checklist-template.md` |
| analyze | `skills/speckit-analyze/templates/analyze-template.md` |
| summary | `skills/speckit-summary/templates/summary-template.md` |

`.shared/templates/` 只保留安装类模板。

## 安装矩阵

安装脚本是 [scripts/link_ai.py](/c:/Users/Lidabao/Desktop/Golang/scripts/link_ai.py)。

| Target | 用户目录安装 | 项目目录安装 | 项目说明文件 | 说明 |
|--------|--------------|--------------|--------------|------|
| `codex` | `~/.codex/prompts`、`~/.codex/skills` | 无额外项目目录 | 项目根 `AGENTS.md` | 不复制 `.shared/agents` 到 `~/.codex/agents` |
| `gemini` | `~/.gemini/antigravity/workflows`、`~/.gemini/antigravity/skills` | 项目内 `.agent/*` | 项目根 `GEMINI.md` | 项目目录会复制 `.gemini/.agents` 整棵树 |
| `copilot` | 无 | 项目内 `.github`、`.github/skills` | `.github/copilot-instructions.md` | 不复制 `.shared/agents` 到项目目录 |
| `github` | 无 | 项目内 `.github`、`.github/skills`、`.github/agents` | 无 | 当前仍保留 `agents` 同步逻辑 |

说明：

- `codex/gemini/copilot` 会生成项目说明文件。
- `.shared/agents` 只作为模板源，不会被复制到运行时 `agents/` 目录。
- 历史残留的运行时 `agents/` 目录不由脚本自动清理。

## 安装脚本

当前支持命令：

| 命令 | 说明 |
|------|------|
| `install` | 安装或覆盖运行时资产，并生成项目说明文件 |
| `delete` | 删除脚本管理的运行时资产 |

示例：

```powershell
python scripts/link_ai.py install --target codex
python scripts/link_ai.py install --target gemini
python scripts/link_ai.py install --target copilot
python scripts/link_ai.py delete --target gemini
```

常用参数：

- `--target`：目标运行时，当前支持 `codex`、`gemini`、`copilot`、`github`
- `--agent`：项目说明文件模板名，默认是 `kratos`
- `--project-dir`：显式指定项目目录
- `--home`：显式指定用户目录
- `--repo-url` / `--branch`：从远程仓库拉取后执行安装

## 最小使用示例

假设当前在一个待接入的项目根目录下，需要给 Codex 安装运行时资产并生成项目说明文件：

```powershell
python scripts/link_ai.py install --target codex --agent kratos
```

执行后应得到：

- 用户目录 `~/.codex/prompts/*`
- 用户目录 `~/.codex/skills/*`
- 项目根 `AGENTS.md`

不会得到：

- `~/.codex/agents/*`

后续流程是：

1. AI 从项目根 `AGENTS.md` 获取当前项目的总说明。
2. 用户通过 `/plan`、`/tasks`、`/implement` 等命令进入 prompt。
3. prompt 再驱动 agent 加载对应 skill。
4. 若当前项目是 Kratos 项目，则 agent 在运行时决定加载 `kratos-patterns` 及其下游 skill。

## 扩展方式

### 新增一个流程 skill

- 放到 `.shared/skills/<skill-name>/`
- 在 `SKILL.md` 中定义触发条件、上下文加载、执行规则
- 如有硬约束，放到 `rules/`
- 如有示例与说明，放到 `reference/`
- 如有产物模板，放到 `templates/`

### 新增一个项目适配器

- 新增一个项目适配 skill，而不是修改 `speckit-*` 去写死项目规则
- 让 agent 在运行时识别项目后决定是否加载它
- 若需要项目说明文件模板，再在 `.shared/agents/` 新增对应模板

### 新增一个 AI target

- 在 [scripts/link_ai.py](/c:/Users/Lidabao/Desktop/Golang/scripts/link_ai.py) 的 `TARGETS` 中新增配置
- 明确用户目录安装项、项目目录安装项、是否生成项目说明文件
- 保持 `.shared/agents` 只作为模板源，不把它重新变成运行时共享目录

## 事实源

- `.shared/*` 才是源资产。
- `.shared/agents/*` 只用于渲染项目级说明文件，不作为运行时共享目录安装。
- `AGENTS.md`、`GEMINI.md`、`.github/copilot-instructions.md` 这类文件是安装脚本渲染出来的运行时产物。

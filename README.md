# AI 工程体系

这是一个安装 AI 能力包的仓库，不是业务服务仓库。
仓库维护的是一套可分发的 `agent + prompt + skill + rule + reference + template` 资产，用于把标准化研发流程安装到 Codex、Gemini、Copilot 等运行时。

## 核心流程

```text
/specify -> /clarify -> /plan -> /tasks -> /analyze -> /implement -> /code-review -> /summary
                    \-> /checklist
```

## 目录结构

```text
.shared/
├── agents/                          # agent 模板源文件
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
├── templates/                       # 安装类模板，仅保留项目指令模板
└── ...

.codex/prompts/                      # Codex 命令入口
.gemini/.agents/workflows/           # Gemini 工作流入口
scripts/link_ai.py                   # 运行时安装脚本
```

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
- 使用顺序：先进入对应 skill，再先读 `rules`，后读 `reference`。

## Speckit 与项目适配器

- `speckit-*` 负责标准化研发流程，本身不绑定具体语言或框架。
- 遇到具体项目时，再加载对应项目适配器 skill。
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

`.shared/templates/` 只保留安装类模板，例如 `project-instructions-template.md`。

## 安装脚本

运行时安装脚本是 [scripts/link_ai.py](/c:/Users/Administrator/Desktop/repo/Golang/scripts/link_ai.py)。

当前支持命令：

| 命令 | 说明 |
|------|------|
| `install` | 安装或覆盖运行时资产 |
| `delete` | 删除脚本管理的运行时资产 |

示例：

```powershell
python scripts/link_ai.py install --target codex
python scripts/link_ai.py install --target gemini
python scripts/link_ai.py delete --target gemini
```

## 事实源

- `.shared/*` 才是源资产。
- `AGENTS.md`、`GEMINI.md`、`.github/copilot-instructions.md` 这类文件是安装脚本渲染出来的运行时产物。

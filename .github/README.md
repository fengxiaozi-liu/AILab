# AI 工程体系

面向 **Go/Kratos 微服务**的 AI 辅助研发系统，基于 SpecKit 规范驱动开发范式，通过 GitHub Copilot Agent 实现从需求到代码的结构化推进。

## 核心流程

```
/constitution ─→ /specify ─→ /clarify ─→ /plan ─→ /tasks ─→ /analyze ─→ /implement ─→ /code-review
   治理原则       需求捕获     深度澄清    技术方案    任务拆解    一致性校验     代码实施       代码审查
                                  ↑                                              ↑              │
                             /checklist                                          └──────────────┘
                            需求质量检查                                          (有问题则修复后再审)
```

## 目录结构

```
.github/
├── copilot-instructions.md          # 主 Agent 指令（AI 进入仓库自动加载）
├── agents/                          # 子 Agent 定义（9 个）
├── prompts/                         # Slash Commands 入口（9 个）
├── skills/                          # 技能（流程型 + 知识型）
│   ├── speckit-*/                   #   流程型：各阶段执行方法论
│   ├── kratos-patterns/             #   知识型：Kratos 框架编码规范
│   │   └── reference/               #     16 份参考规范文档
│   └── golang-patterns/             #   知识型：Go 语言最佳实践
├── rules/                           # 项目规则
│   └── project/                     #   编码约束、安全、Go 语言规范
├── .specify/
│   ├── memory/
│   │   └── constitution.md          # 宪法（项目治理基线）
│   └── templates/                   # 产物模板
│       ├── spec-template.md
│       ├── plan-template.md
│       ├── tasks-template.md
│       └── checklist-template.md
└── docs/                            # 工程文档
```

## 组件说明

### 主 Agent 指令

| 文件 | 说明 |
|------|------|
| `copilot-instructions.md` | 总控与路由，AI 进入仓库自动加载；定义角色、子 Agent 调度、安全规则 |

### Commands × Agents

每个命令对应一个 Agent，Agent 加载对应 Skill 执行。

| 命令 | Agent | 加载 Skill | 职责 | 产出物 |
|------|-------|-----------|------|--------|
| `/constitution` | `constitution` | speckit-constitution | 创建/修订项目治理原则 | `constitution.md` |
| `/specify` | `specify` | speckit-specify | 将自然语言需求结构化 | `spec.md` |
| `/clarify` | `clarify` | speckit-clarify | 对 spec 中的 CQ 逐题深度澄清 | 更新 `spec.md` |
| `/checklist` | `checklist` | speckit-checklist | 生成需求质量检查清单（独立工具） | `checklist.md` |
| `/plan` | `plan` | speckit-plan + kratos-patterns | 技术方案设计 | `plan.md` |
| `/tasks` | `tasks` | speckit-tasks + kratos-patterns | 按 Kratos 工作流拆解可执行任务 | `tasks.md` |
| `/analyze` | `analyze` | speckit-analyze | 对 spec/plan/tasks 三件套做只读一致性检查 | 只读报告 |
| `/implement` | `implement` | speckit-implement + kratos-patterns | 按 tasks.md 逐 Phase 执行代码实现 | 代码变更 |
| `/code-review` | `code-review` | speckit-code-review + kratos-patterns | 完成度/架构合规/性能/安全审查 | `review.md` |

### Skills（技能）

分为两类，职责互补：

| 类型 | Skill | 说明 |
|------|-------|------|
| **流程型** | `speckit-constitution` | 宪法治理方法论 |
| | `speckit-specify` | 需求捕获方法论（9 维度 CQ 评估） |
| | `speckit-clarify` | 深度澄清方法论 |
| | `speckit-checklist` | 需求质量清单方法论 |
| | `speckit-plan` | 技术方案设计方法论（Phase 0/1 + 框架适配 + 测试策略） |
| | `speckit-tasks` | 任务拆解方法论（RQ 追踪 + Phase 结构） |
| | `speckit-analyze` | 一致性分析方法论（7 类检测 + 4 级严重度） |
| | `speckit-implement` | 实施执行方法论（逐 Phase + 断点续做） |
| | `speckit-code-review` | 代码审查方法论（4 维度评分 + 40 分制） |
| **知识型** | `kratos-patterns` | Kratos 框架编码规范，按项目类型分流（BaseService/业务/网关） |
| | `golang-patterns` | Go 语言惯用法与最佳实践 |

#### kratos-patterns 参考文档（16 份）

| 参考文档 | 说明 |
|---------|------|
| `naming-spec.md` | 命名规范 |
| `project-spec.md` | 项目结构规范 |
| `proto-spec.md` | Proto 文件规范 |
| `codegen-spec.md` | 代码生成规范（Wire 等） |
| `ent-spec.md` | Ent ORM Schema 规范 |
| `enum-spec.md` | 枚举规范 |
| `error-spec.md` | 异常规范 |
| `layer-spec.md` | 层规范（biz/data/service） |
| `depend-spec.md` | InnerRPC 依赖包装规范 |
| `server-spec.md` | 服务注册规范 |
| `gateway-spec.md` | 网关层规范 |
| `config-spec.md` | 配置规范 |
| `i18n-spec.md` | 国际化规范 |
| `crontab-spec.md` | 定时任务规范 |
| `wire-spec.md` | Wire DI 规范 |
| `testing-spec.md` | 测试规范（mockgen + enttest） |

### Rules（项目规则）

| 文件 | 说明 |
|------|------|
| `project/constraints.md` | 全局约束（契约优先、变更确认） |
| `project/security.md` | 安全规则（输入校验、脱敏） |
| `project/coding-conventions.md` | 工程治理（生成物、错误、日志） |
| `project/go-language.md` | Go 语言规范 |

### Templates（产物模板）

所有 Agent 输出遵循标准模板，存放于 `.specify/templates/`：

| 模板 | 使用者 | 说明 |
|------|--------|------|
| `spec-template.md` | specify agent | 需求规格模板（需求描述 + CQ 块） |
| `plan-template.md` | plan agent | 技术方案模板（Phase 0/1 + 框架适配 + 测试策略 + 风险） |
| `tasks-template.md` | tasks agent | 任务清单模板（Phase 分组 + RQ 追踪） |
| `checklist-template.md` | checklist agent | 质量清单模板 |

### 产出物目录

功能需求的产出物存放在 `specs/<feature>/` 目录下：

```
specs/<feature>/
├── spec.md           # 需求规格
├── plan.md           # 技术方案
├── tasks.md          # 任务清单
├── checklist.md      # 质量清单（可选）
└── review.md         # 代码审查报告（可选）
```

## 架构设计原则

| 原则 | 说明 |
|------|------|
| **SKILL = 方法论** | 详细的领域执行指南（"怎么做"），由 Agent 加载 |
| **Agent = 编排** | 简洁的流程调度（"做什么"），引用 Skill 执行 |
| **Prompt = 入口** | 极简触发点，关联 Agent |
| **Template = 契约** | 标准化产出格式，减少自由偏差 |
| **宪法优先** | `constitution.md` 是治理基线，高于具体执行偏好 |
| **证据优先** | 禁止臆测仓库事实，必须用工具验证 |

## 新增组件规范

### 新增 Agent

1. 在 `agents/` 创建 `<name>.agent.md`，frontmatter 包含 `name`、`description`、`handoffs`
2. 在 `prompts/` 创建 `<name>.prompt.md`，frontmatter 指定 `agent: <name>`
3. 在 `skills/` 创建 `speckit-<name>/SKILL.md`，定义执行方法论
4. 更新 `copilot-instructions.md` 的子 Agent 表和命令映射

### 新增 Skill

1. 创建目录 `skills/<skill-name>/`
2. 入口文件 `SKILL.md`，详细参考放在 `reference/` 或 `example/` 子目录
3. 在对应 Agent 中声明加载

### 新增 Rule

- 项目级硬规则 → `rules/project/`
- 写法：包含原则/规范/禁止/示例，保持简洁

## 编码注意（Windows）

本仓库文档以 UTF-8 为准。Windows/PowerShell 下请使用支持 UTF-8 的编辑器/终端。

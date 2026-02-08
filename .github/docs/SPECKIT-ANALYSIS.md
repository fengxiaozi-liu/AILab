# SpecKit 范式解析 & 当前实现差异分析

> 本文档包含三部分：  
> 1. SpecKit 标准范式详解  
> 2. 当前 `.github/` 实现的映射关系  
> 3. 差异分析与改进建议

---

## 一、SpecKit 标准范式详解

### 1.1 什么是 SpecKit

SpecKit 是一套**规范驱动开发（Spec-Driven Development）**方法论框架，核心理念：

| 理念 | 说明 |
|------|------|
| 意图驱动 | 先定义需求（What/Why），再实现（How） |
| 多阶段细化 | 从模糊想法到可执行任务的结构化推进 |
| 技术无关 | 规范层不绑定语言/框架，实现层才涉及技术选型 |
| AI 协助执行 | 每个阶段由 AI Agent 驱动，规范文档作为 Agent 的输入/输出契约 |
| 多运行时适配 | 同一套 Skills 可在 Claude Code / Codex / Copilot / Gemini 上运行 |

### 1.2 核心流程（7 阶段）

```
1. Constitution（治理原则）
      ↓
2. Specify（需求规格）  ──或──  Baseline（从代码逆向生成规格）
      ↓
3. Clarify（澄清，可选）
      ↓
4. Plan（技术方案）
      ↓
5. Tasks（任务拆解）
      ↓
6. Analyze（一致性校验，可选）
      ↓
7. Implement（实施）
```

#### 附加阶段

| 阶段 | 说明 | 使用时机 |
|------|------|----------|
| Baseline | 从现有代码逆向生成规格 | 旧系统文档化、重构前 |
| Checklist | 生成质量检查清单 | 任意阶段，"需求的单元测试" |
| TasksToIssues | 将 tasks.md 转为 GitHub Issues | 项目管理集成 |

### 1.3 各阶段 Skill 详解

#### 阶段 1：speckit-constitution（治理原则）

- **目标**：定义项目不可违反的原则与治理规则
- **输入**：用户提供的原则或修正
- **输出**：`.specify/memory/constitution.md`
- **关键机制**：
  - 模板中的 `[PLACEHOLDER]` 替换为具体值
  - 版本语义化管理（MAJOR/MINOR/PATCH）
  - 一致性传播：更新后自动校验 plan-template、spec-template、tasks-template 的对齐
  - 生成 Sync Impact Report（变更影响报告）
- **执行频率**：项目初始化时执行一次，后续仅在治理变更时修订

#### 阶段 2：speckit-specify（需求规格）

- **目标**：将自然语言需求转化为结构化规格文档
- **输入**：用户的功能描述
- **输出**：`specs/N-feature-name/spec.md` + `checklists/requirements.md`
- **关键机制**：
  - 自动生成短名称（2-4 词，action-noun 格式）
  - 检查已有分支编号，自动递增
  - 运行 `create-new-feature.sh` 创建 Git 分支 + 目录
  - 基于模板填充：用户故事、功能需求、成功标准、关键实体
  - `[NEEDS CLARIFICATION]` 标记（最多 3 个）
  - 规格质量验证 Checklist 自动生成
- **关键规则**：
  - 规格必须技术无关（不放实现细节）
  - 需求必须可测试

#### 阶段 2-alt：speckit-baseline（逆向规格）

- **目标**：从现有代码分析逆向生成需求规格
- **输入**：代码文件/目录/glob 模式
- **输出**：同 specify（spec.md + checklist）
- **关键机制**：
  - 自动识别入口点、公共接口、数据模型、API 端点
  - 将技术实现转化为业务需求（HOW → WHAT/WHY）
  - 示例：`if (user.role === 'admin')` → "System MUST restrict action to administrator users"

#### 阶段 3：speckit-clarify（澄清）

- **目标**：识别规格中的模糊点，通过交互式提问消除歧义
- **输入**：已有的 spec.md
- **输出**：更新后的 spec.md（含 Clarifications 章节）
- **关键机制**：
  - 9 维度歧义扫描分类法（功能范围、领域模型、交互流程、非功能属性、集成、边缘情况、约束、术语、完成信号）
  - 每次只问 1 个问题，最多 5 个（总上限 10 个）
  - 每个问题提供推荐答案 + 选项表格
  - 增量写入 spec.md（每回答一个立即更新）
  - 按 Impact × Uncertainty 启发式排序

#### 阶段 4：speckit-plan（技术方案）

- **目标**：基于规格生成技术实施策略
- **输入**：spec.md + constitution.md
- **输出**：plan.md, research.md, data-model.md, contracts/, quickstart.md
- **关键机制**：
  - Phase 0 Outline & Research：提取未知项 → 调研 → 形成决策记录
  - Phase 1 Design & Contracts：实体提取 → 数据模型 → API 契约（OpenAPI/GraphQL）
  - Constitution Check：方案必须通过原则门禁
  - agent-context 更新脚本自动执行

#### 阶段 5：speckit-tasks（任务拆解）

- **目标**：将方案拆分为可执行的依赖有序任务列表
- **输入**：plan.md + spec.md + 可选制品
- **输出**：tasks.md
- **关键机制**：
  - 按用户故事组织任务（非按技术层）
  - 严格 Checklist 格式：`- [ ] [TaskID] [P?] [Story?] Description with file path`
  - `[P]` 标记可并行任务
  - Phase 结构：Setup → Foundation → UserStory(P1→P2→P3) → Polish
  - MVP 范围建议
  - 依赖图 + 并行执行示例

#### 阶段 6：speckit-analyze（一致性分析）

- **目标**：在实施前进行跨制品一致性和质量分析
- **输入**：spec.md, plan.md, tasks.md, constitution.md
- **输出**：只读分析报告（Markdown，不修改任何文件）
- **关键机制**：
  - 6 类检测：重复、歧义、欠规范、原则对齐、覆盖缺口、不一致
  - 4 级严重度：CRITICAL / HIGH / MEDIUM / LOW
  - 需求到任务的覆盖映射表
  - Constitution 违规自动标记为 CRITICAL
  - 最多 50 个发现，超出部分汇总

#### 阶段 7：speckit-implement（实施）

- **目标**：按 tasks.md 逐条执行实现
- **输入**：tasks.md + plan.md + 可选制品
- **输出**：代码变更 + tasks.md 状态更新
- **关键机制**：
  - Checklist 前置检查（不完整则询问是否继续）
  - 自动创建/校验 ignore 文件（.gitignore, .dockerignore 等）
  - Phase-by-phase 执行，尊重依赖顺序
  - TDD 导向：有测试时先写测试再写实现
  - 完成一个任务立即标记 `[X]`
  - 错误处理：非并行任务失败时暂停

#### 附加：speckit-checklist（定制检查清单）

- **核心理念**："需求的单元测试"——不是校验实现是否正确，而是校验需求描述是否完整、清晰、一致
- **动态澄清**：从用户请求中提取信号，生成 1-3 个上下文相关问题
- **可在任意阶段使用**

#### 附加：speckit-taskstoissues（任务转 Issue）

- **目标**：将 tasks.md 转为 GitHub Issues
- **安全约束**：仅在确认是 GitHub 远端仓库时才执行

### 1.4 文件产出结构

```
.specify/
├── memory/
│   └── constitution.md           # 项目原则
└── templates/
    ├── spec-template.md
    ├── plan-template.md
    └── tasks-template.md

specs/
└── N-feature-name/
    ├── spec.md                   # 需求规格
    ├── plan.md                   # 技术方案
    ├── research.md               # 技术调研
    ├── data-model.md             # 实体定义
    ├── tasks.md                  # 任务拆解
    ├── quickstart.md             # 快速启动指南
    ├── contracts/                # API 契约
    │   ├── api.yaml
    │   └── schema.graphql
    └── checklists/               # 质量检查
        └── requirements.md
```

### 1.5 多运行时架构

SpecKit 采用符号链接策略实现"一次编写，多端运行"：

```
skills/                 ← 单一事实源
  ├── speckit-*/SKILL.md
.claude/skills  → ../skills
.codex/skills   → ../skills
.github/skills  → ../skills
```

| 运行时 | Agent 定义 | 命令/Prompt | Skills |
|--------|-----------|-------------|--------|
| Claude Code | — | .claude/commands/ | .claude/skills → ../skills |
| Codex CLI | — | .codex/prompts/ | .codex/skills → ../skills |
| GitHub Copilot | .github/agents/ | .github/prompts/ | .github/skills → ../skills |
| Gemini CLI | — | .gemini/commands/ | — |

### 1.6 Copilot 运行时适配（标准方式）

在标准 SpecKit 中，Copilot 的适配方式：

| 组件 | 文件模式 | 数量 |
|------|---------|------|
| Agent | `.github/agents/speckit.{phase}.agent.md` | 10 个（每阶段一个） |
| Prompt | `.github/prompts/speckit.{phase}.prompt.md` | 10 个（极简，仅关联 agent） |
| Skill | `.github/skills/speckit-{phase}/SKILL.md` | 10 个（共享源） |

Prompt 文件极简（通常 3-5 行），仅指定 `agent:` 字段，所有逻辑在 Agent 文件中。  
Agent 文件包含完整工作流，通过 `handoffs` 声明阶段间的衔接。

---

## 二、当前实现映射关系

### 2.1 组件对照表

| SpecKit 标准范式 | Copilot 标准实现方式 | 当前实现 | 当前文件 |
|-----------------|-------------------|---------|---------|
| AGENTS.md | copilot-instructions.md | ✅ 已实现 | `.github/copilot-instructions.md` |
| speckit.plan.agent.md | agents/planner.agent.md | ⚠️ 已实现（自定义流程） | `.github/agents/planner.agent.md` |
| speckit.implement.agent.md | agents/task-executor.agent.md | ⚠️ 已实现（自定义流程） | `.github/agents/task-executor.agent.md` |
| speckit.plan.prompt.md | prompts/plan.prompt.md | ✅ 已实现 | `.github/prompts/plan.prompt.md` |
| speckit.implement.prompt.md | prompts/task.prompt.md | ✅ 已实现 | `.github/prompts/task.prompt.md` |
| speckit-*/SKILL.md (10 个) | skills/ (共享) | ⚠️ 部分 | 仅 3 个技能（kratos/golang/backend） |
| speckit.constitution | — | ❌ 缺失 | — |
| speckit.specify | — | ❌ 缺失 | — |
| speckit.clarify | — | ❌ 融入 planner 阶段 A | — |
| speckit.tasks | — | ❌ 融入 planner 阶段 C | — |
| speckit.analyze | — | ❌ 缺失 | — |
| speckit.baseline | — | ❌ 缺失 | — |
| speckit.checklist | — | ❌ 缺失 | — |
| speckit.taskstoissues | — | ❌ 缺失 | — |
| .specify/templates/ | — | ❌ 缺失 | — |
| .specify/memory/ | — | ❌ 缺失 | — |
| .specify/scripts/ | — | ❌ 缺失 | — |
| rules/ | — | ✅ 已实现 | `.github/rules/` (5 规则文件) |
| docs/ | — | ✅ 已实现 | `.github/docs/AI-ENGINEERING-GUIDE.md` |

### 2.2 流程对照

```
SpecKit 标准流程：
Constitution → Specify → Clarify → Plan → Tasks → Analyze → Implement
     1           2         3        4       5        6          7

当前实现流程：
（无 Constitution）→（无 Specify）→ planner(A:澄清 → B:设计 → C:生成文档+任务) → task-executor
                                    ▲ 将 Clarify+Plan+Tasks 压缩到 1 个 Agent
```

---

## 三、差异分析与改进建议

### 3.1 架构层面差异

#### 差异 1：阶段压缩过度——planner 承担了 3 个 SpecKit 阶段

| 维度 | SpecKit 标准 | 当前实现 |
|------|-------------|---------|
| 澄清 | speckit-clarify（独立 Agent，9 维度歧义扫描） | planner 阶段 A（仅"关键决策点"检查） |
| 方案设计 | speckit-plan（含 research.md, data-model.md, contracts/） | planner 阶段 B（方案对比，无独立产物） |
| 任务拆解 | speckit-tasks（严格 Checklist 格式 + 用户故事维度） | planner 阶段 C（按技术层拆解） |

**影响**：
- 澄清深度不够（缺少 9 维度扫描、推荐答案机制、增量更新 spec）
- 方案设计无中间产物（research.md, data-model.md, contracts/ 均缺失）
- 任务拆解无用户故事关联、无并行标记 `[P]`、无 MVP 范围建议

**改进建议**：
- **短期**：在 planner 的阶段 A 中增加更系统的歧义检测维度（参考 SpecKit 9 维度）；阶段 B 输出时增加 research 产物；阶段 C 任务格式对齐 SpecKit checklist 格式
- **中期**：将 planner 拆分为 clarify + plan + tasks 三个独立 Agent，职责更清晰

#### 差异 2：缺少 Constitution（治理原则）

SpecKit 要求所有方案设计必须通过 Constitution 门禁校验。当前实现没有这一层：
- 没有 `.specify/memory/constitution.md`
- Plan 没有 "Constitution Check" 环节
- 没有项目级不可违反的原则声明

**影响**：方案设计缺少"红线约束"，可能生成违反项目原则的方案。

**改进建议**：
- 创建 `.github/docs/constitution.md` 或复用规则体系（`rules/project/` 已部分承担此角色）
- 在 planner 阶段 B 增加 "规则校验" 步骤，交叉检查 `rules/project/` 中的约束

#### 差异 3：缺少 Specify（需求规格化）

SpecKit 在方案设计之前有一个专门的需求结构化阶段，输出标准化的 spec.md。当前实现直接从自然语言需求跳到方案设计。

**影响**：
- 没有标准化的需求文档格式
- 需求理解仅作为 planner 输出的一部分，无独立存档
- 无法复用 spec 进行后续的 analyze / checklist 等质量校验

**改进建议**：
- 新增 `specifier` Agent 或在 planner 阶段 A 前增加一个 "需求结构化" 预处理步骤
- 输出标准化的 spec 文件到 `plan/` 或 `specs/` 目录

#### 差异 4：缺少 Analyze（一致性校验）

SpecKit 有专门的跨制品一致性分析阶段，在实施前验证 spec ↔ plan ↔ tasks 的一致性。

**影响**：实施前无法发现需求-方案-任务之间的遗漏、冲突和术语漂移。

**改进建议**：
- 新增 `analyzer` Agent 或在 task-executor 执行前增加一个校验步骤
- 可利用现有 rules 体系作为校验基准

#### 差异 5：缺少模板和产出物标准化

| SpecKit 模板/产物 | 当前实现 |
|------------------|---------|
| spec-template.md | ❌ 无 |
| plan-template.md | ❌ 无（planner 自由格式输出） |
| tasks-template.md | ❌ 无（planner 自有格式） |
| research.md | ❌ 无 |
| data-model.md | ❌ 无 |
| contracts/ | ❌ 无 |
| quickstart.md | ❌ 无 |
| constitution.md | ❌ 无 |

**改进建议**：
- 创建 `plan/templates/` 目录，定义设计文档和任务列表的标准模板
- planner 输出时强制使用模板（减少自由格式偏差）

### 3.2 Agent 设计差异

#### 差异 6：Agent 数量与粒度

| 维度 | SpecKit 标准 | 当前实现 |
|------|-------------|---------|
| Agent 数量 | 10 个（每阶段一个） | 2 个（planner + task-executor） |
| 单 Agent 职责 | 单一阶段职责 | planner 承担 3 个阶段 |
| Agent 间衔接 | handoffs 声明式 | 主 Agent 手动编排 |

**影响**：
- planner 过于复杂（3 段式流程在单 Agent 中）
- 缺少 handoffs 声明，阶段转换依赖主 Agent 编排

**改进建议**：
- 对 planner 最少拆分为 2 个 Agent（clarifier + planner），或保持现状但增加阶段间的产物标准化
- 利用 Copilot Agent 的 `handoffs` 功能声明阶段衔接

#### 差异 7：Agent 无 Skill 引用

SpecKit 标准中 Agent 文件通过 SKILL.md 获取详细工作流指导。当前实现的 Agent 直接内嵌全部逻辑。

| 维度 | SpecKit 标准 | 当前实现 |
|------|-------------|---------|
| 工作流详情 | 在 SKILL.md 中（~100-300 行） | 直接写在 Agent 文件中（~140 行） |
| Skill 引用 | Agent 引用 SKILL.md | Agent 仅通过 `rule` 引用规则文件，不引用 Skill |

**影响**：Agent 文件中逻辑与工作流混合，不利于复用和维护。

**改进建议**：
- 在 Agent 的 front-matter 或工作流中显式引用对应的 Skill 文件
- 将详细的阶段执行逻辑迁移到 Skill 中，Agent 只保留流程编排

### 3.3 Skills 层面差异

#### 差异 8：Skills 定位不同

| 维度 | SpecKit Skills | 当前 Skills |
|------|---------------|-------------|
| 性质 | 工作流驱动（每阶段的详细执行指南） | 技术知识库（Go/Kratos/后端架构的最佳实践） |
| 触发方式 | 随阶段自动激活 | 按需手动加载 |
| 数量 | 10 个流程 Skill | 3 个知识 Skill |
| 内容 | 步骤化脚本（1. 运行命令 2. 读取文件 3. 填充模板...） | 模式/规范文档 |

当前 Skills 实际是"技术知识库"，不是"工作流 Skill"。两者并不矛盾，可以共存。

**改进建议**：
- 保留现有 3 个知识型 Skill（它们是项目特有的技术规范）
- 增加流程型 Skill（参考 SpecKit 的 speckit-plan/speckit-tasks 等），让 Agent 在执行时有标准化步骤指导

### 3.4 Rules 层面差异

#### 差异 9：Rules 体系是当前实现的优势

当前实现在 Rules 方面比 SpecKit 标准范式做得更细致：

| 维度 | SpecKit 标准 | 当前实现 |
|------|-------------|---------|
| 规则分层 | 无独立规则体系（靠 Constitution + SKILL.md 内嵌） | 3 层分离：ai/ + agents/ + project/ |
| AI 编排规则 | 无（靠 AGENTS.md 总纲） | `ai/agents.md` + `ai/interaction.md` |
| Agent 专属规则 | 无 | `agents/planner.md` + `agents/task-executor.md` |
| 项目规则 | 无（靠 Constitution） | `project/constraints.md` + `project/security.md` + `project/coding-conventions.md` + `project/naming.md` + `project/go-language.md` |

**这是当前实现的一个亮点**：SpecKit 的 Constitution 是一个通用容器，而当前的 Rules 体系按受众分片（AI/Agent/项目），更精准。

### 3.5 产出物管理差异

#### 差异 10：产出物目录与命名

| 维度 | SpecKit 标准 | 当前实现 |
|------|-------------|---------|
| 目录结构 | `specs/N-feature-name/` | `plan/[feature-name]-*.md` |
| 分支管理 | 自动编号分支 (`N-feature-name`) | 无自动分支管理 |
| 文件命名 | `spec.md`, `plan.md`, `tasks.md` | `[feature]-design.md`, `[feature]-task.md` |
| 完成报告 | 直接在 tasks.md 上标记 | 独立 `[feature]-done.md` |

**改进建议**：
- 考虑采用 SpecKit 的 `specs/N-feature/` 目录结构，每个 feature 独立目录，更清晰
- 但 `done.md` 完成报告是当前实现的优势，SpecKit 没有这个机制

### 3.6 工程化差异

#### 差异 11：缺少辅助脚本

SpecKit 提供了一套 `.specify/scripts/bash/` 辅助脚本：
- `create-new-feature.sh`：创建特性分支 + 目录
- `setup-plan.sh`：初始化方案目录
- `check-prerequisites.sh`：验证前置条件
- `update-agent-context.sh`：更新 Agent 上下文

当前实现完全依赖 Agent 手动操作文件系统。

**改进建议**：
- 创建简单的辅助脚本（PowerShell/bash），自动化分支创建、目录初始化等重复操作
- 至少实现 `check-prerequisites` 和 `create-feature` 两个脚本

#### 差异 12：缺少 CI 验证

SpecKit 有 `.github/workflows/ci.yml` 验证规范文件的一致性（symlinks、格式等）。

当前实现无 CI 验证环节。

### 3.7 文档与指南差异

#### 差异 13：AI-ENGINEERING-GUIDE 是优势

当前实现有一份非常完善的 `AI-ENGINEERING-GUIDE.md`（616 行），定义了：
- 短规则原则（单文件 ≤150 行）
- 单一事实来源原则
- 表格优于段落原则
- 示例精简原则
- 信号优于口令原则
- 正向指令原则
- 各组件结构模板与好坏示例
- 维护检查清单

**这是 SpecKit 标准范式中没有但非常有价值的内容**。SpecKit 侧重"做什么"，而 AI-ENGINEERING-GUIDE 侧重"怎么写好 Agent 配置"。

---

## 四、改进优先级建议

### P1 — 高优（结构性缺陷）

| 编号 | 改进项 | 说明 | 复杂度 |
|------|--------|------|--------|
| I-01 | 增加需求结构化阶段 | 在 planner 之前增加 Specify 步骤，输出标准 spec 文件 | 中 |
| I-02 | 标准化产物模板 | 创建 design/task 的 Markdown 模板，planner 使用模板填充 | 低 |
| I-03 | 标准化任务格式 | 采用 SpecKit 的 Checklist 格式（TaskID + 并行标记 + 故事关联） | 低 |
| I-04 | 增加 Constitution / 红线原则 | 汇总 rules/project/ 形成一份"不可违反原则"清单，方案设计时强制校验 | 低 |

### P2 — 中优（质量提升）

| 编号 | 改进项 | 说明 | 复杂度 |
|------|--------|------|--------|
| I-05 | 增加分析/校验阶段 | 在 task-executor 执行前增加 spec ↔ plan ↔ task 一致性检查 | 中 |
| I-06 | Planner 澄清深度增强 | 引入 SpecKit 的 9 维度歧义扫描 + 推荐答案机制 | 中 |
| I-07 | 增加方案中间产物 | planner 阶段 B 额外输出 research.md + data-model.md | 中 |
| I-08 | Agent 引用 Skill | Agent front-matter 中声明引用的 Skill 文件 | 低 |

### P3 — 低优（锦上添花）

| 编号 | 改进项 | 说明 | 复杂度 |
|------|--------|------|--------|
| I-09 | 增加 Baseline Agent | 从现有代码逆向生成需求规格 | 高 |
| I-10 | 增加 Checklist Agent | 实现"需求的单元测试"机制 | 中 |
| I-11 | 增加辅助脚本 | 自动化分支创建、目录初始化 | 低 |
| I-12 | 产出物目录重组 | 从 `plan/` 改为 `specs/N-feature/` 结构 | 中 |
| I-13 | 增加 TasksToIssues | 将任务列表转为 GitHub Issues | 低 |
| I-14 | 拆分 Planner 为多个 Agent | clarifier + planner + taskgen | 高 |

---

## 五、当前实现的独特优势（SpecKit 未覆盖）

不应丢弃这些优势，它们是对 SpecKit 范式的有价值补充：

| 优势 | 说明 |
|------|------|
| 精细化 Rules 体系 | 3 层分离（ai/agents/project），比 Constitution 更精准 |
| AI-ENGINEERING-GUIDE | 如何编写好的 Agent 配置的元指南 |
| 项目特化 Skills | kratos-patterns 按项目类型（BaseService/业务/网关）分流 |
| golang-patterns 按需加载 | 8 个子文件按场景加载，避免上下文爆炸 |
| done.md 完成报告 | 实施完成后有独立输出报告 |
| 防漂移锚点 | Agent 编排中的 Anti-Drift 机制 |
| 信号机制 | WAIT_CONFIRM / DONE 简洁有效 |

---

## 六、总结

当前实现是一个**面向 Go/Kratos 微服务定制化的 Agent 系统**，采用了 SpecKit 范式的部分思想（阶段化、Agent 分工、信号机制），但与完整 SpecKit 范式相比，主要差距在于：

1. **阶段覆盖不全**：7 个标准阶段中仅实现了约 3 个（Clarify/Plan/Tasks 压缩在 planner 中 + Implement）
2. **产出物不标准**：缺少模板驱动、缺少中间产物（research/data-model/contracts）
3. **质量门禁不足**：缺少 Constitution 门禁和 Analyze 阶段

同时，当前实现也有 SpecKit 标准范式未覆盖的亮点（Rules 分层、AI 工程指南、项目特化 Skills、防漂移机制）。

**最佳路径**：以现有架构为基础，按 P1→P2→P3 优先级逐步补齐 SpecKit 缺失的阶段和产物，同时保留当前的差异化优势。

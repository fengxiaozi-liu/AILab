<!--
Sync Impact Report
- Version change: 1.4.0 -> 1.4.1
- Change type: PATCH（修正原则结构：补齐 III，并将每条原则收敛为单一 Rule Source）
- Principles changed:
  - [KRATOS_PRINCIPLE_I]
  - [KRATOS_PRINCIPLE_II]
  - [KRATOS_PRINCIPLE_III]
  - [KRATOS_PRINCIPLE_IV]
  - [KRATOS_PRINCIPLE_V]
- Governance changed:
  - 明确每条原则仅允许一个 Rule Source
- Deferred TODOs: None
-->

# [项目宪法]

## 核心原则

### [Contract-First（契约优先）]

- Principle-ID: `[KRATOS_PRINCIPLE_I]`
- 说明：跨边界能力与高风险变更必须遵循“先契约、后实现、可追溯”的治理原则。
- Rule Source: `.github/rules/project/constraints.md`

### [Security-By-Default（默认安全）]

- Principle-ID: `[KRATOS_PRINCIPLE_II]`
- 说明：所有输入默认不可信，安全校验必须前置，敏感信息必须边界内处理。
- Rule Source: `.github/rules/project/security.md`

### [Engineering-Consistency（工程一致性）]

- Principle-ID: `[KRATOS_PRINCIPLE_III]`
- 说明：工程产物必须保持一致、可维护、可审计。
- Rule Source: `.github/rules/project/coding-conventions.md`

### [Unified-Errors-and-Observability（统一错误与可观测性）]

- Principle-ID: `[KRATOS_PRINCIPLE_IV]`
- 说明：错误语义与日志结构必须统一，确保跨层可诊断与线上可观测。
- Rule Source: `.github/rules/project/go-language.md`

### [Reliability-and-Performance-Baseline（可靠性与性能基线）]

- Principle-ID: `[KRATOS_PRINCIPLE_V]`
- 说明：可靠性与性能属于功能正确性，不应作为后置优化项。
- Rule Source: `.github/rules/project/constraints.md`

## 治理

- 宪法优先级：本宪法高于流程偏好；若冲突，先修流程文件以满足宪法。
- 作用域：本宪法仅定义项目治理原则，约束 `specify`、`plan`、`tasks`、`implement` 全流程。
- Rule Application：
  - 实施细则以 `.github/rules/project/*.md` 为准。
  - 当原则抽象与落地细则有歧义时，以 rules 中可执行条款和 Good/Bad demo 作为执行依据。
  - 每条 Principle-ID 仅绑定一个 Rule Source 文件。
- 更新联动要求：
  - 每次宪法修订 MUST 同步评估并更新 `rules/project/*`。
  - 每次宪法修订 MUST 同步评估并更新运行资产：`.github/agents/*`、`.github/prompts/*`（命令入口）。
- 跨 Agent 引用规范：
  - 其他 agent 引用宪法时 MUST 使用 Principle-ID（如 `[KRATOS_PRINCIPLE_I]`、`[KRATOS_PRINCIPLE_II]`）。
- 修订流程：
  1) 提交修订动机与影响范围。
  2) 更新本文件与 Sync Impact Report。
  3) 同步核对依赖入口文档与 prompts。
- 版本策略：
  - MAJOR：原则删除/重定义/不兼容治理变更。
  - MINOR：新增原则/章节或实质扩展约束。
  - PATCH：措辞澄清、排版修复、非语义调整。

**版本**: 1.4.1 | **批准日期**: 2026-02-08 | **最后修订**: 2026-02-08

---
name: kratos-conventions
description: Kratos 横切语义约定与共享字典收敛。用于 error semantic、not found、enum/常量、i18n key、logging 字段、comment 规范等跨模块共享约定。不要用于业务编排、组件接入、proto/wire/codegen 或纯代码审查。
---

# Kratos Conventions

## 何时使用

- 修改统一错误语义、`not found` 语义或 error helper。
- 收敛 enum、稳定常量值域、reason 命名和共享 key。
- 调整 i18n key、多语言文案、日志字段或共享注释规范。

## 职责边界

- 本技能负责横切语义约定、共享字典和表达规范。
- 本技能不负责业务编排、聚合建模、组件接入、proto/wire/codegen，也不负责纯风格自检。

## 输入

- 必需：当前横切语义改动的范围，例如 `error`、`enum`、`i18n`、`logging`、`comment`。
- 可选：涉及的模块、文件路径、已有命名模式或共享前缀。
- 可选：需要重点治理的风险点，例如裸业务字符串、reason 漂移、日志字段散落、无效注释。

缺少必需输入时，MUST 先从工作区和任务上下文补齐；仍无法判断改动属于哪类横切约定时，再向用户提问，不得直接套模板。

## 工作流

### 收集证据与补齐输入

- 优先使用：改动文件路径、模块归属、现有命名模式（reason/key/module/enum 名）、已有共享定义位置。
- 输入不足时：先在仓库中检索现有实现与共享定义；仍不足再向用户追问（不要猜）。

### 判定约定类型（只选一个主类型）

- `error`：错误函数、reason、not found、错误模板、错误语义收敛。
- `enum`：状态/原因码/稳定值域/类型化常量与跨层联动。
- `i18n`：i18n key 与文案、localize 调用、生成流程约束。
- `logging`：logger/helper 约定、module 字段、WithContext(ctx)、日志级别与脱敏。
- `comment`：对象/字段/流程注释规范、proto/ent 注释边界、生成物注释约束。

若同时命中多个类型：以“共享定义的主改动点”为主类型，其他作为联动点处理。

### 先复用检查再读 references

- 先检索项目是否已有同类共享定义与命名模式（error/enum/i18n/logging/comment）。
- 能复用就复用，避免新造平行体系（新 reason/key/module/enum 前缀）。

### 按类型按需加载 references

- `error` -> 读 `references/error-spec.md`
- `enum` -> 读 `references/enum-spec.md`
- `i18n` -> 读 `references/i18n-spec.md`
- `logging` -> 读 `references/logging-spec.md`
- `comment` -> 读 `references/comment-spec.md`

### 实施共享约定改动（只做约定职责）

- 先收敛业务语义，再决定 error/enum/i18n/logging/comment 的表达方式，不让横切层反向定义业务。
- 优先复用已有命名模式与共享定义位置；避免散落裸字符串与魔法值。
- 禁止通过手改派生产物/生成物注释来替代修改源定义。

### 边界自检（不做测试验收）

- 是否以项目既有共享定义为主（给出证据：找到的定义位置/命名模式/相似实现）。
- 是否对照了本次类型对应的 references（至少确认 MUST/MUST NOT 没踩）。
- 是否引入了平行体系（reason/key/enum/module 变体）、是否让错误语义与领域语义脱节。

## 约束

### MUST

- MUST 先定义业务语义，再决定 error、enum、i18n、logging 和 comment 的表达方式。
- MUST 优先复用现有 enum、常量、reason、i18n key 前缀和日志字段约定，不新造平行体系。
- MUST 让 not found、状态流转、日志字段和文案表达与当前领域语义保持一致。
- MUST 避免裸业务字符串、魔法值和临时文案散落在业务流程中。
- MUST 只加载当前横切主题直接相关的 references。

### MUST NOT

- MUST NOT 用 `nil, nil` 或空对象伪装 not found 成功。
- MUST NOT 把错误描述、中文文案或状态值直接散落在业务逻辑里。
- MUST NOT 手改派生产物或生成物注释来替代修改源定义。
- MUST NOT 在高并发循环里无节制打印日志或记录敏感信息。
- MUST NOT 用注释重复显而易见的代码行为。

### SHOULD

- SHOULD 在新增错误语义时同步检查 enum、i18n 和日志检索字段是否需要联动。
- SHOULD 在参与 `switch`、状态流转和筛选的值域上优先类型化。
- SHOULD 让注释优先解释“为什么存在”和“负责什么”，而不是“这一行做了什么”。

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `references/error-spec.md` | 修改错误函数、reason、not found、错误模板和业务错误表达 | 涉及错误语义治理时加载 |
| `references/enum-spec.md` | 修改状态、原因码、稳定值域和类型化常量 | 涉及 enum/常量治理时加载 |
| `references/i18n-spec.md` | 修改 i18n key、文案模板、多语言源文件和 localize 调用 | 涉及多语言与文案语义时加载 |
| `references/logging-spec.md` | 修改 logger、module 字段、上下文日志、生命周期日志和日志级别 | 涉及日志治理时加载 |
| `references/comment-spec.md` | 修改对象、字段、流程、proto 和 ent 注释的写法与边界 | 涉及注释质量治理时加载 |

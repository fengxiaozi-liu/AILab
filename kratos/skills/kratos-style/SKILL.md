---
name: kratos-style
description: Kratos 代码风格与防御式实现审查。用于 code review、自检、收尾治理、`nil, nil`、吞错、duplicate helper、命名/术语收敛（含跨层命名漂移）、运行期兜底等坏味修正。它是审查型技能，不替代 domain、service、components、pkg 等主执行技能。
---

# Kratos Style

## 何时使用

- 对已有业务代码或即将提交的改动做 Code Review、自检或收尾打磨。
- 集中治理 `nil, nil`、吞错、重复 helper、relation 补查、命名漂移。
- 在功能已基本落地后，检查是否存在不必要的运行期兜底和重复封装。

## 职责边界

- 本技能是审查型技能，负责风格、错误路径和实现洁净度治理。
- 本技能不主导业务开发，不负责组件接入、proto/wire/codegen、服务注册或业务编排设计。

## 输入

- 必需：待审查的改动范围、模块、文件路径或风险点。
- 可选：重点检查主题，例如错误路径、重复 helper、命名收敛、relation 补查。
- 可选：已有实现文件路径或 PR 范围。

缺少必需输入时，MUST 先从当前工作区和任务上下文补齐；仍无法判断审查范围时，再向用户提问，不得无依据猜测。

## 工作流

### 收集证据与补齐输入

- 优先使用：待审查的改动范围、文件路径、风险点描述、已有实现与 helper 位置。
- 输入不足时：先在仓库中定位改动面与相似实现；仍不足再向用户追问（不要猜）。

### 判定审查主题（只选一个主类型）

- `defensive`：错误路径、吞错、重复 helper、relation 补查、运行期兜底等失败路径治理。
- `naming`：领域术语、跨层对象名、文件名、proto/message/service 命名收敛。

若同时命中多个主题：以“最可能导致线上风险/架构腐化的点”为主类型，其它作为联动点处理。

### 按类型按需加载 references

- `defensive` -> 读 `references/defensive-spec.md`
- `naming` -> 读 `references/naming-spec.md`

### 实施审查与修正（只做审查职责）

- 只做必要的风格治理与风险修正；不主导业务实现、不替代主执行技能。
- 避免把框架保证或主流程保证重复写回业务代码；避免新增一次性 helper。

### 边界自检（不做测试验收）

- 是否只在审查/收尾阶段使用本技能（未替代 domain/service/components/pkg）。
- 是否避免吞错、`nil, nil`、重复 helper、越界 relation 补查与运行期兜底重建。

## 约束

### MUST

- MUST 只在风格审查、风险治理或收尾阶段使用本技能。
- MUST 优先加载与当前风险最相关的 references。
- MUST 先检查 Wire、框架生命周期和主流程是否已提供保证，再决定是否保留防御代码。
- MUST 先检查 `internal/pkg`、现有 helper 和已有实现，再决定是否删除或复用 helper。
- MUST 在函数签名包含 `error` 时返回真实错误，不得把失败吞掉改成成功路径。
- MUST 保持同一领域概念在 aggregate、repo、usecase、service、proto 间使用同一术语。

### MUST NOT

- MUST NOT 把本技能当作常规功能开发的主技能。
- MUST NOT 在 Wire 已保证依赖的前提下做运行期二次装配、判空兜底或“防御性”重建。
- MUST NOT 对主流程已保证完整的数据做下游 relation 补查。
- MUST NOT 为单次调用或单点逻辑新增包装 helper。
- MUST NOT 新增与现有 `internal/pkg` 或项目内能力重复的 helper。
- MUST NOT 返回 `nil, nil`、`_ = err`、只打日志不返回错误，或通过前置短路掩盖真实失败路径。

### SHOULD

- SHOULD 在命名修改时同步检查文件名、结构体名、接口名、proto message 和 service 名是否一起收敛。
- SHOULD 在收尾时明确说明删除了哪些保护代码、复用了哪些 helper、修复了哪些真实错误路径。
- SHOULD 把本技能作为主执行技能后的补充审查，而不是与主技能并列抢占任务。

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `references/defensive-spec.md` | 错误路径、吞错、`nil, nil`、重复 helper、relation 补查、运行期兜底 | 涉及失败路径或防御式治理时加载 |
| `references/naming-spec.md` | 统一领域术语、跨层对象命名、文件名与 proto/message/service 命名 | 涉及命名新增、重命名或历史漂移清理时加载 |

---
name: kratos-pkg
description: Kratos `internal/pkg` 通用公共能力沉淀。用于 util、context、metadata、middleware、proto-helper、schema、seata，以及重复 helper 的公共化收敛。不要用于组件接入、业务编排或单一业务场景抽象。
---

# Kratos Pkg

## 何时使用

- 需要沉淀或复用脱离具体业务语义的通用能力，例如 util、context/metadata 透传、中间件、协议辅助工具。
- 需要判断一段 helper 是否应该进入 `internal/pkg`，以及应落在哪个子目录。
- 需要收敛重复 helper，改为复用已有公共能力。

## 职责边界

- 本技能只负责 `internal/pkg` 级别的通用能力沉淀与复用判断。
- 本技能不负责基础设施组件接入，不负责业务流程编排，也不负责单一业务场景抽象。

## 输入

- 必需：当前计划沉淀或复用的公共能力范围，例如 `util`、`context`、`metadata`、`middleware`、`proto-helper`、`schema`、`seata`。
- 可选：候选代码片段、已有 helper 名称或目录位置。
- 可选：需要重点回答的问题，例如“是否值得下沉”“应放哪个子目录”“是否已有现成能力”。

缺少必需输入时，MUST 先从工作区和任务上下文补齐；仍无法判断是否真的是公共能力时，再向用户提问，不得为了“看起来通用”就直接下沉。

## 工作流

### 收集证据与补齐输入

- 优先使用：候选代码片段/调用点、预计复用点数量、已有 `internal/pkg` 目录结构与现存能力、目标子目录候选。
- 输入不足时：先在仓库中检索现有实现与复用情况；仍不足再向用户追问（不要猜）。

### 判定能力类型（只选一个主类型）

- `util`：无状态通用函数（字符串/URL/并发等）。
- `context`：请求级上下文 helper（typed `WithXxx/GetXxx`）。
- `metadata`：metadata key 与读写收口。
- `middleware`：通用中间件与链路聚合。
- `proto-helper`：协议辅助转换（paging/filter/time range 等）。
- `schema`：结构解析、反射辅助与 schema 工具。
- `seata`：分布式事务/seata 辅助。

若同时命中多个类型：以“对外导出 API 的主能力面”为主类型，其它作为联动点处理。

### 先复用检查再读 references

- 先检索 `internal/pkg` 是否已有可复用实现与相邻能力（避免重复造轮子）。
- 判断是否满足公共边界：脱离具体业务语义后仍成立，且预期会被多个模块稳定复用。

### 按类型按需加载 references

- `util` -> 读 `references/util-spec.md`
- `context` -> 读 `references/context-spec.md`
- `metadata` -> 读 `references/metadata-spec.md`
- `middleware` -> 读 `references/middleware-spec.md`
- `proto-helper` -> 读 `references/proto-helper-spec.md`
- `schema` -> 读 `references/schema-spec.md`
- `seata` -> 读 `references/seata-spec.md`
- 若需要新增 `internal/pkg` 子域（新目录）或新增对外导出 API -> 读 `references/subdomain-spec.md`

### 实施公共能力沉淀（只做 pkg 职责）

- 能复用就复用；不能复用就留在原层；只有确实满足公共边界时才新增 `internal/pkg` 子域。
- 每个 helper 只解决一个稳定能力，不承载业务流程编排；目录名直接表达边界。

### 边界自检（不做测试验收）

- 是否先检索并复用了现有 `internal/pkg` 实现（给出证据：相似能力位置/路径）。
- 是否避免引入业务语义、业务状态判断与 repo 访问。
- 是否避免过度抽象与一次性复用；子目录落点是否能表达边界。

## 约束

### MUST

- MUST 先检索现有 `internal/pkg` 实现，再决定是否新增 helper 或子目录。
- MUST 只把脱离具体业务语义后仍成立、且会被多个模块稳定复用的能力下沉到 `internal/pkg`。
- MUST 让一个 helper 只解决一个稳定能力，不承载业务流程编排。
- MUST 使用能直接表达边界的子目录名，避免 `common`、`helper`、`misc`。

### MUST NOT

- MUST NOT 把只在一个调用点成立的逻辑强行下沉到 `internal/pkg`。
- MUST NOT 在 `internal/pkg` 中访问 repo、业务 context、大型聚合对象或外部业务语义。
- MUST NOT 用公共包包装具体业务流程、状态判断或临时页面逻辑。
- MUST NOT 把组件接入、生命周期注册或 consumer/listener/event 逻辑塞进 `pkg`。

### SHOULD

- SHOULD 按主题拆分文件和子目录，让能力边界能从路径直接读出来。
- SHOULD 在 context、metadata、middleware 三者存在透传关系时一起检查边界是否一致。
- SHOULD 在 proto helper、schema、seata 这类工具层能力上优先保持最小能力面，避免抽象膨胀。

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `references/util-spec.md` | 新增无状态通用函数、字符串/URL/并发等小型 helper | 涉及 util 类通用能力时加载 |
| `references/context-spec.md` | 新增请求级 context helper 和 typed `WithXxx/GetXxx` 收口 | 涉及 context 边界时加载 |
| `references/metadata-spec.md` | 新增 metadata key、读写 helper 和来源优先级 | 涉及 metadata 透传时加载 |
| `references/middleware-spec.md` | 新增通用中间件和链路聚合函数 | 涉及 middleware 设计时加载 |
| `references/proto-helper-spec.md` | 新增协议辅助转换、paging/filter/time range 等 proto helper | 涉及协议工具时加载 |
| `references/schema-spec.md` | 新增结构提取、反射辅助和 schema 类工具 | 涉及 schema 工具时加载 |
| `references/seata-spec.md` | 新增分布式事务/seata 辅助和事务封装 | 涉及 seata 或事务工具时加载 |
| `references/subdomain-spec.md` | 新增 `internal/pkg` 子域（新目录）、命名边界、最小导出 API 与落位规范 | 需要新建子域或新增导出 API 时加载 |

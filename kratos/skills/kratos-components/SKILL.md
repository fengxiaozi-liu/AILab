---
name: kratos-components
description: Kratos 基础设施组件接入与生命周期装配。用于 ent/schema、consumer、listener、cron、eventbus/local event、provider 注册、生命周期启动停止、组件接线和最小联动检查。不要用于业务编排、`internal/pkg` 公共能力沉淀或纯代码审查。
---

# Kratos Components

## 何时使用

- 新增或修改基础设施组件，例如 Ent、MQ Consumer、Cron、Local Event。
- 调整组件的 provider、注册入口、生命周期启动/停止或联动链路。
- 判断基础设施组件应该放哪里、如何接线、如何做最小联动检查。

## 职责边界

- 本技能负责基础设施组件接入、注册、生命周期和相关联动。
- 本技能不负责业务流程编排，不负责通用公共能力沉淀，也不负责接入层协议设计。

## 输入

- 必需：当前组件改动的类型与范围，例如 `ent`、`cron`、`consumer`、`listener`、`event`。
- 可选：涉及的目录、文件路径或注册入口。
- 可选：需要重点关注的接线点，例如生命周期、provider、注册器、生成命令、定向测试。

缺少必需输入时，MUST 先从工作区和任务上下文补齐；仍无法判断属于哪类组件时，再向用户提问，不得把所有组件 references 一次性全读。

## 工作流

### 收集证据与补齐输入

- 优先使用：改动文件路径、目录结构、wire/provider 注册点、启动/生命周期入口、topic/queue/payload、ent/schema 目录。
- 输入不足时：先在仓库中检索现有实现与注册入口；仍不足再向用户追问（不要猜）。

### 判定组件类型（只选一个主类型）

- `ent`：改动集中在 `ent/schema`、ORM 结构、生成链路。
- `cron`：涉及 scheduler/定时任务注册、启动停止、幂等/并发策略。
- `consumer`：涉及 MQ 消费者、topic/queue、消费回调、重试/ack、消费失败处理。
- `listener`：涉及本地 listener、handler 注册聚合、事件回调/转发边界。
- `event`：涉及 eventbus/local event、发布订阅约定、topic、payload。

若同时命中多个类型：以“注册/生命周期归属的核心改动点”为主类型，其他类型作为联动点处理。

### 先复用检查再读 references

- 先检索项目是否已有同类组件、注册入口、生命周期管理与目录落点。
- 能复用就复用，避免重复注册、散落注册、重复 helper。

### 按类型按需加载 references

- `ent` -> 读 `references/ent-spec.md`
- `cron` -> 读 `references/cron-spec.md`
- `consumer` -> 读 `references/consumer-spec.md`
- `listener` -> 读 `references/listener-spec.md`
- `event` -> 读 `references/event-spec.md`

### 实施组件接入与联动（只做组件职责）

- 组件层只做：接入、注册、生命周期、必要的最小转发/适配。
- 业务流程编排、事务边界、状态机：交给 UseCase（不下沉到 component）。
- 按 references 约束完成落位、注册和联动；必要时同步更新 `wire`、`codegen` 或生命周期装配。

## 约束

### MUST

- MUST 只加载当前组件直接相关的 references。
- MUST 让组件代码只承担基础设施职责，业务流程、事务边界和状态机仍交给 UseCase。
- MUST 保持组件注册、生命周期启动/停止和 provider 接线统一收口，不散落在多个模块。
- MUST 在 ent/schema、provider/wire 或启动注册变化后补齐对应生成与最小构建检查
### MUST NOT

- MUST NOT 在 cron、consumer、listener 中直接写完整业务流程编排。
- MUST NOT 吞掉解析失败、消费失败、发布失败或启动失败的错误。
- MUST NOT 在没有唯一约束或明确幂等语义时滥用 upsert/OnConflict。
- MUST NOT 让组件自行分散管理生命周期或重复注册。
- MUST NOT 为单一业务场景伪装出通用组件抽象或把公共能力硬塞进组件层。

### SHOULD

- SHOULD 在组件接入前先检索项目现有实现，尽量沿用现有容器、注册器和目录结构。
- SHOULD 在 consumer/listener/event 场景中保持 topic、queue、payload 和触发事实的语义清晰。
- SHOULD 在新增组件时同步做最小联动检查，例如 codegen、wire、build。

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `references/ent-spec.md` | 修改 Ent schema、field、edge、index、upsert 和 ORM 生成 | 涉及 ent/schema 和 ORM 生成时加载 |
| `references/cron-spec.md` | 修改定时任务、注册器、启动停止、幂等和并发策略 | 涉及 cron/scheduler 时加载 |
| `references/consumer-spec.md` | 修改消息消费、topic/queue、注册器、生命周期和错误返回 | 涉及 consumer 接入时加载 |
| `references/listener-spec.md` | 修改本地 listener、回调 handler、聚合注册和最小转发逻辑 | 涉及 listener 接入时加载 |
| `references/event-spec.md` | 修改 local event/eventbus 的 topic、payload、发布订阅约定 | 涉及本地事件和 eventbus 时加载 |

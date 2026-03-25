---
name: kratos-service
description: Kratos 接入层协议设计与服务暴露。用于 proto、HTTP/gRPC 接口、gateway 代理、OpenAPI 注解、server register、provider/wire、codegen 联动。不要用于业务编排、聚合建模、`internal/pkg` 下沉或纯代码审查。
---

# Kratos Service

## 何时使用

- 编写或修改 proto 契约、HTTP/gRPC 接口层。
- 调整 `internal/service`、`internal/server`、gateway、OpenAPI 注解或服务注册。
- 修改 provider、wire 装配入口和相关 codegen 链路。

## 职责边界

- 本技能负责接入层协议适配、服务暴露、注册装配和生成联动。
- 本技能不负责业务编排、聚合建模、公共能力下沉或纯风格审查。

## 输入

- 必需：本次改动的接入层目标与范围，例如 `service/server`、`gateway`、`proto`、`wire`、`codegen`。
- 可选：涉及的文件路径，例如 `internal/service/*`、`internal/server/*`、`api/**/*.proto`。
- 可选：需要重点关注的联动点，例如路由变更、协议变更、生成命令、OpenAPI 展示。

缺少必需输入时，MUST 先从工作区和任务上下文补齐；仍无法判断变更类型时，再向用户提问，不得无依据猜测。

## 工作流

### 收集证据与补齐输入

- 优先使用：改动文件路径、接口目标（HTTP/gRPC/gateway）、proto 变更点、注册入口、wire/provider 位置、现有 codegen 入口。
- 输入不足时：先在仓库中检索现有同类接口/注册/装配实现；仍不足再向用户追问（不要猜）。

### 判定接入层类型（只选一个主类型）

- `proto`：proto 契约、service/message 结构、side 边界、公共协议抽取。
- `server`：`internal/server`、服务注册、路由暴露与启动装配。
- `service`：`internal/service` 协议适配、参数转换、调用 usecase。
- `gateway`：gateway 代理、聚合转发、参数映射与响应适配。
- `wire`：provider、构造函数依赖、ProviderSet、wire 装配入口。
- `codegen`：proto/wire/schema/注册链路变化后的生成联动。
- `openapi-v3`：OpenAPI 注解与文档展示。

若同时命中多个类型：以“对外契约/暴露链路的主改动点”为主类型，其它作为联动点处理。

### 先复用检查再读 references

- 先检索项目是否已有同类接口的 proto/注册/适配/装配实现与目录落点。
- 能复用就复用，避免新增平行注册入口或重复映射规则。

### 按类型按需加载 references

- `server` / `service` -> 读 `references/server-spec.md`
- `gateway` -> 读 `references/gateway-spec.md`
- `openapi-v3` -> 读 `references/openapi-v3-spec.md`
- `proto` -> 读 `references/proto-spec.md`
- `wire` -> 读 `references/wire-spec.md`
- `codegen` -> 读 `references/codegen-spec.md`

### 实施接入层改动（只做接入层职责）

- 接入层只做协议适配、注册、代理和生成联动，不把业务编排下沉到 `service`、`gateway`、`wire`。
- 禁止手改生成产物；源定义变更后按项目链路补齐生成联动。

### 边界自检（不做测试验收）

- 是否以项目既有实现为主（给出证据：相似接口/注册入口/映射规则的位置）。
- 是否对照了本次类型对应的 references（至少确认 MUST/MUST NOT 没踩）。
- 是否保持边界：`service` 不编排业务，`gateway` 不维护业务状态，注册入口不散落。

## 约束

### MUST

- MUST 只加载与当前改动直接相关的 references。
- MUST 让 `service` 层只承担协议适配、参数转换和调用 usecase，不承载业务编排。
- MUST 让 `gateway` 只承担代理、协议转换、参数映射和响应适配，不维护业务状态。
- MUST 在 proto、wire、schema、注册链路变化后补做对应 codegen 和最小 build 验证。
- MUST 保持服务注册入口统一，新增接口时同步检查 service 实现、server 注册、provider/wire 和生成链路。

### MUST NOT

- MUST NOT 手改生成产物来替代修改源定义。
- MUST NOT 在 `internal/service`、gateway 或 server 层编排完整业务流程。
- MUST NOT 把服务注册散落到多个非统一入口文件。
- MUST NOT 在 proto 未稳定时抢先用 wire 反推设计。
- MUST NOT 让 side 间直接互引业务 proto，公共协议必须放到 `base/*`。

### SHOULD

- SHOULD 在新增或修改接口时同步检查 HTTP/gRPC 注册、gateway 映射、OpenAPI 展示和 provider 装配是否一起收敛。
- SHOULD 在 codegen 失败时优先回看源定义和注册关系，而不是直接修补产物。
- SHOULD 在只涉及业务编排或聚合边界时切换到 `kratos-domain`，不要继续停留在本技能。

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `references/server-spec.md` | 修改 `internal/service`、`internal/server`、服务注册和接入层落位 | 涉及 service/server/register 时加载 |
| `references/gateway-spec.md` | 修改 gateway 代理、参数映射、聚合转发、协议适配 | 涉及 gateway 逻辑时加载 |
| `references/openapi-v3-spec.md` | 修改 OpenAPI v3 注解、中文展示、字段说明与文档一致性 | 涉及文档注解或 HTTP proto 展示时加载 |
| `references/proto-spec.md` | 修改 proto 契约、service/message 结构、side 边界、公共协议抽取 | 涉及 proto 设计时加载 |
| `references/wire-spec.md` | 修改 provider、构造函数依赖、ProviderSet、wire 装配入口 | 涉及依赖注入和装配时加载 |
| `references/codegen-spec.md` | 修改 proto、wire、schema 或注册链路后的生成和校验 | 涉及生成物联动时加载 |

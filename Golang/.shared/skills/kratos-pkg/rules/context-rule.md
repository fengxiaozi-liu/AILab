# Context Rule

## Principles

- 业务上下文使用统一的 typed helper 读写，避免散落的原始 key。
- `context` 负责聚合运行态基础信息，不负责承载业务对象图。

## Specification

- 通过统一的 `WithTenant`、`WithLocalize`、`WithViewer` 一类 helper 注入上下文。
- 新增上下文字段时，同步评估是否需要联动 `metadata` 与 `middleware`。
- Context 中承载的值应是稳定的运行态基础信息，如租户、语言、viewer、本次请求的元数据。

## Prohibit

- 禁止在业务代码中直接使用裸字符串 key 读写 context。
- 禁止把大对象、聚合关系、临时业务结果直接塞进 context。
- 禁止同一语义同时由 context struct、metadata key、middleware 重复维护而不做收口。

## Checklist

- 是否已有现成的 context helper 可复用。
- 是否需要联动 metadata 和 middleware。
- 是否只放入了请求级基础信息，而非业务对象。

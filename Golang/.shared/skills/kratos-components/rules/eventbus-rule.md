# EventBus Rule

## Principles

- 事件驱动只传递边界清晰的事件语义，不替代业务编排。

## Specification

- 事件定义、发布、监听职责清晰。
- Listener 保持幂等和可重试。
- 事件传递遵守 `kratos-domain` 的对象复用规则，优先复用已有稳定对象，不因传递场景额外新建近义壳结构。

正例：

```go
Payload: account
Payload: store
```

反例：

```go
Payload: &AccountAfterOpenEventPayload{Account: account}
Payload: &CommitPageEventData{Store: store}
```

## Prohibit

- 禁止在事件监听器中隐藏核心业务编排。
- 禁止仅为 EventBus 传递场景新建 `EventPayload`、`ContextDTO`、`XxxData` 这类近义传递结构，除非事件语义已脱离原对象。

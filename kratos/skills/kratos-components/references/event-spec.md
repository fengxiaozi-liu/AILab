# Event Spec

## 事件值域

| 条件 | 做法 |
|------|------|
| topic/key 是稳定值域 | 集中定义为 enum/常量 |
| 魔法字面量 topic | 不允许 |

```go
// ✅
const TopicAccountAfterOpen = "account.after_open"
```

```go
// ❌
Topic: "account.after_open"
```

---

## payload 选择

| 条件 | 做法 |
|------|------|
| 现有聚合根或稳定对象已足够 | 直接复用 |
| 需要裁剪字段或隔离边界 | 新增专门 payload |

```go
// ✅ 直接传稳定对象
Payload: account
```

```go
// ✅ 边界明确时单独 payload
type AccountOpenedPayload struct {
    ID uint32
}
```

```go
// ❌ 临时 map
Payload: map[string]any{"id": account.ID}
```

---

## 发布与订阅边界

| 条件 | 做法 |
|------|------|
| 发布方 | 只发布事实，不为订阅方额外拼装 |
| 订阅方 | 只解析事实，再调用 UseCase 或下游消息 |

```go
// ✅
if err := u.eventBus.Publish(ctx, &eventbus.Event{
    Topic:   TopicAccountAfterOpen,
    Payload: account,
}); err != nil {
    return err
}
```

```go
// ❌ 把事件当事务补偿机制乱用
Publish(...)
// 希望订阅方回补事务失败
```

---

## 组合场景

```text
enum topic
-> UseCase publish fact
-> listener register
-> listener callback
```

这个组合场景同时满足：

- topic 集中定义
- payload 语义稳定
- 发布与订阅边界清晰

---

## 常见错误模式

```text
// ❌ topic 名散落
```

```text
// ❌ payload 用匿名 map
```

```text
// ❌ 事件替代事务边界
```

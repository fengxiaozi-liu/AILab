# Listener Spec

## Listener 边界

| 条件 | 做法 |
|------|------|
| 本地事件 listener / callback handler | 做解析、最小校验、组装入参、调用下游 |
| repo 聚合装配、跨资源编排 | 不放在 listener |

```go
// ✅
type AccountListener struct {
    producer *amqp.Producer
}

func (l *AccountListener) afterOpen(ctx context.Context, account *openbiz.Account) error {
    return l.producer.Publish(&amqp.Event{Event: openenum.AmqpEventOpenFirstReviewCheck.Value()})
}
```

```go
// ❌ listener 里查 repo + 业务决策
func (l *AccountListener) afterOpen(ctx context.Context, account *openbiz.Account) error {
    page, _ := l.repo.GetPage(ctx, account.ID)
    ...
}
```

---

## 注册聚合

| 条件 | 做法 |
|------|------|
| 新增 listener | 通过统一 listener 聚合注册 |
| 到处手工挂载 | 不允许 |

```go
// ✅
func (l *AccountListener) register(eventBus *eventbus.EventBus) {
    eventBus.Register(&eventbus.CallBack{
        Topics: []string{openenum.LocalEventAccountAfterOpen.Value()},
        Func:   l.afterOpen,
    })
}
```

```go
// ❌ 分散注册
eventBus.Register(...)
eventBus.Register(...)
```

---

## 组合场景

```text
UseCase publish event
-> listener aggregate register
-> listener callback
-> producer / usecase
```

这个组合场景同时满足：

- listener 足够薄
- 注册统一
- 错误直接返回

---

## 常见错误模式

```text
// ❌ listener 吞错
```

```text
// ❌ listener 做 repo 编排
```

```text
// ❌ 注册散落
```

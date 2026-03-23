# EventBus Spec

## Payload 对象选择

| 条件 | 做法 |
|------|------|
| 已有聚合对象能完整表达事件语义 | 直接用已有对象作为 Payload |
| 需要裁剪/脱敏/跨协议投影 | 新建专用结构体 |
| 仅因为"传递"场景就想包一层 | 禁止，MUST NOT 新建近义壳结构 |

```go
// ✅ 直接用已有聚合对象
err = u.eventBus.Publish(ctx, &eventbus.Event{
    Topic:   openenum.LocalEventAccountAfterOpen.Value(),
    Payload: account,
})

// ✅ 多个对象用 store 聚合，直接传
err = u.eventBus.Publish(ctx, &eventbus.Event{
    Topic:   openenum.LocalEventPageCommitAfter.Value(),
    Payload: store,
})

// ❌ 仅为传递而包装的近义壳结构
type AccountAfterOpenEventPayload struct {
    Account *Account
}
err = u.eventBus.Publish(ctx, &eventbus.Event{
    Payload: &AccountAfterOpenEventPayload{Account: account},
})

// ❌ ContextDTO 壳
type CommitPageEventData struct {
    Store *AccountFlowPageStore
}
```

---

## Listener 职责

```go
// ✅ Listener 只做后置独立动作，不包含核心业务编排
func (l *AccountListener) handleAccountAfterOpen(ctx context.Context, event *eventbus.Event) error {
    account, ok := event.Payload.(*biz.Account)
    if !ok {
        return errors.New(400, "INVALID_PAYLOAD", "invalid payload type")
    }
    return l.notifyUseCase.SendOpenNotification(ctx, account.ID)
}

// ❌ Listener 承担完整业务编排（实际上是第二个 UseCase）
func (l *AccountListener) handleAccountAfterOpen(ctx context.Context, event *eventbus.Event) error {
    account := event.Payload.(*biz.Account)
    account.Status = openenum.AccountStatusOpened
    if err := l.repo.UpdateAccount(ctx, account); err != nil { ... }
    if err := l.repo.CreateLog(ctx, ...); err != nil { ... }
    if err := l.sendNotification(ctx, account); err != nil { ... }
    return nil
}
```

---

## Listener 注册

```go
// ✅ 统一注册，名称、Topic、Handle 清晰
func (l *AccountListener) register(eventBus *eventbus.EventBus) {
    listeners := []*eventbus.Listener{
        {
            Name:   "account.after.open",
            Topics: []string{openenum.LocalEventAccountAfterOpen.Value()},
            Handle: l.handleAccountAfterOpen,
        },
    }
    for _, listener := range listeners {
        eventBus.RegisterListener(listener)
    }
}
```

---

## 组合场景

```go
// UseCase 发布事件 + Listener 幂等消费
func (u *AccountUseCase) passReview(ctx context.Context, account *Account) error {
    if err := u.doOpenAccount(ctx, account); err != nil {
        return err
    }
    // ✅ Payload 直接用已有 account 对象
    return u.eventBus.Publish(ctx, &eventbus.Event{
        Topic:   openenum.LocalEventAccountAfterOpen.Value(),
        Payload: account,
    })
}

// ✅ Listener 幂等：重复执行不会破坏数据
func (l *AccountListener) handleAccountAfterOpen(ctx context.Context, event *eventbus.Event) error {
    account := event.Payload.(*biz.Account)
    return l.uc.EnsureNotificationSent(ctx, account.ID) // 幂等操作
}
```

---

## 常见错误模式

```go
// ❌ 新建近义 EventPayload
type AccountOpenedEvent struct { Account *biz.Account }

// ❌ Listener 内持有事务和多步写操作
func (l *Listener) handle(...) error {
    tx.InTx(ctx, func(ctx context.Context) error {
        // 大量业务逻辑...
    })
}

// ❌ 发布但不检查 err
u.eventBus.Publish(ctx, event)
```

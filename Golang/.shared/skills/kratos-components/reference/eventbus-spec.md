# EventBus Reference

## 这个主题解决什么问题

说明本地 EventBus 中事件、发布者、监听器如何组织。

## 适用场景

- 新增事件
- 发布业务事件
- 实现 Listener

## 设计意图

EventBus 参考解释的是事件如何承担“后置动作解耦”的角色，而不是变成另一条业务主流程。

- 事件适合表达主流程完成后需要被其他模块感知和继续处理的事实。
- 事件名称、载荷和监听器职责清晰后，更容易判断哪些逻辑适合拆成监听器。
- 解释清楚事件边界，也能减少监听器不断膨胀成第二个 UseCase 的风险。

## 实施提示

- 先明确主流程结果，再决定是否需要发事件。
- 先设计稳定的事件载荷，再安排监听器拆分。
- 如果一个监听器需要过多业务上下文，通常说明事件边界还可以再抽象。
- 传递对象时优先复用已有聚合对象、应用层对象、`store`、`filter` 等稳定结构，不要仅为了事件通道再包一层近义 payload。这条来自 `kratos-domain` 的对象复用规范，EventBus 只是其中一个典型场景。

## 推荐结构

- 事件定义围绕领域语义命名
- 发布者只负责在合适时机发事件
- Listener 负责消费事件并执行后续动作

## 标准模板

```go
type AccountCreated struct {
    AccountID uint32
}
```

## 代码示例参考

```go
func (u *AccountUseCase) Create(ctx context.Context, in *CreateAccountInput) error {
    account, err := u.accountRepo.CreateAccount(ctx, in)
    if err != nil {
        return err
    }
    u.eventBus.Publish(&event.AccountCreated{AccountID: account.ID})
    return nil
}
```

## 项目通用发布示例

```go
func (u *AccountUseCase) passReview(ctx context.Context, account *Account) error {
    if err := u.doOpenAccount(ctx, account); err != nil {
        return err
    }

    return u.eventBus.Publish(ctx, &eventbus.Event{
        Topic:   openenum.LocalEventAccountAfterOpen.Value(),
        Payload: account,
    })
}
```

## 页面提交前后事件示例

```go
if err := u.eventBus.Publish(ctx, &eventbus.Event{
    Topic:   openenum.LocalEventPageCommitBefore.Value(),
    Payload: store,
}); err != nil {
    return 0, err
}

return u.eventBus.Publish(ctx, &eventbus.Event{
    Topic:   openenum.LocalEventPageCommitAfter.Value(),
    Payload: store,
})
```

## Good Example

- 事件名和负载可直接表达业务含义
- Listener 输入输出边界清晰
- 已有对象已能表达语义时，直接作为 `Payload` 传递，不额外创建 `EventPayload`

## Listener 注册示例

```go
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

## 常见坑

- 事件名太泛，后续难区分
- 事件负载塞入过多与监听目标无关的字段
- 监听器承担复杂业务编排

## 相关 Rule

- `../rules/eventbus-rule.md`
- `../../kratos-domain/rules/aggregate-rule.md`

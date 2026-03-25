# Defensive Spec

## 依赖与运行期兜底

| 条件 | 做法 |
|------|------|
| 依赖由 Wire 或框架生命周期保证 | 直接使用，不做判空或二次装配 |
| 依赖来自外部输入或动态查找 | 在边界处校验并返回明确错误 |

```go
// ✅ Wire 保证的依赖直接使用
type OpenPageCommitListener struct {
    hookUseCase *openbiz.HookUseCase
}

func (l *OpenPageCommitListener) Handle() error {
    return l.hookUseCase.Do()
}
```

```go
// ❌ 对注入依赖做运行期兜底
func (l *OpenPageCommitListener) Handle() error {
    if l.hookUseCase == nil {
        return errors.New("hookUseCase is nil")
    }
    return l.hookUseCase.Do()
}
```

```go
// ⚠️ 只有动态来源才做边界校验
func (u *AccountUseCase) Sync(ctx context.Context, provider Provider) error {
    client := provider.Client("open")
    if client == nil {
        return errors.New("open client not found")
    }
    return client.Sync(ctx)
}
```

---

## 主流程数据与 relation 补查

| 条件 | 做法 |
|------|------|
| 主流程已保证 payload 完整 | 直接下游传递，不补查 relation |
| 当前函数负责组装聚合数据 | 在组装阶段一次性查齐 |

```go
// ✅ 主流程已准备完整数据，直接发布
return u.eventBus.Publish(ctx, &eventbus.Event{
    Topic:   openenum.LocalEventAccountAfterOpen.Value(),
    Payload: account,
})
```

```go
// ❌ 下游临时补查 relation
if account.CollectInfo == nil {
    collect, err := u.collectRepo.GetAccountCollectByFilter(ctx, account.ID)
    if err != nil {
        return err
    }
    account.CollectInfo = collect
}
```

```go
// ⚠️ 当前函数就是聚合装配点，可以一次性查齐
func (u *AccountUseCase) BuildAccount(ctx context.Context, id uint32) (*Account, error) {
    account, err := u.repo.GetAccount(ctx, id)
    if err != nil {
        return nil, err
    }
    collect, err := u.collectRepo.GetAccountCollectByFilter(ctx, id)
    if err != nil {
        return nil, err
    }
    account.CollectInfo = collect
    return account, nil
}
```

---

## 单点逻辑与 helper 抽象

| 条件 | 做法 |
|------|------|
| 逻辑只在一个调用点使用 | 就地展开 |
| 逻辑被多个模块复用且语义稳定 | 提取到已有公共层或新增最小 helper |

```go
// ✅ 单点逻辑直接展开
switch store.Account.OpenStatus {
case openenum.AccountOpenStatusFirstChecking:
    var data map[string]interface{}
    if err := helper.StructConvert(store.Account, &data); err != nil {
        return err
    }
    return l.producer.Publish(&amqp.Event{
        Event: openenum.AmqpEventOpenFirstReviewCheck.Value(),
        Data:  data,
    })
}
return nil
```

```go
// ❌ 为单次调用抽包装 helper
func (l *OpenPageCommitListener) Handle() error {
    return l.publishReviewCheck(store.Account, openenum.AmqpEventOpenFirstReviewCheck)
}
```

```go
// ⚠️ 被多个调用点复用且语义稳定时才值得抽象
func publishReviewEvent(producer *amqp.Producer, event string, payload any) error {
    return producer.Publish(&amqp.Event{Event: event, Data: payload})
}
```

---

## 错误返回与吞错

| 条件 | 做法 |
|------|------|
| 函数签名包含 `error` | 失败必须返回错误 |
| 失败可降级且已明确约定 | 返回降级结果并保留错误语义 |

```go
// ✅ 失败直接返回 error
func (r *accountRepo) GetAccount(ctx context.Context, id uint32) (*biz.Account, error) {
    info, err := r.query(ctx, id)
    if err != nil {
        return nil, err
    }
    return info, nil
}
```

```go
// ❌ 返回 nil, nil 掩盖失败
func (r *accountRepo) GetAccount(ctx context.Context, id uint32) (*biz.Account, error) {
    info, err := r.query(ctx, id)
    if err != nil {
        return nil, nil
    }
    return info, nil
}
```

```go
// ❌ 只记日志不返回
if err := u.repo.Save(ctx, account); err != nil {
    log.Error(err)
    return nil
}
```

```go
// ⚠️ 降级必须保留语义，而不是伪装成功
reply, err := u.remote.Query(ctx, req)
if err != nil {
    return &Reply{Status: "degraded"}, err
}
return reply, nil
```

---

## 公共能力复用

| 条件 | 做法 |
|------|------|
| `internal/pkg` 或项目内已有能力可用 | 直接复用 |
| 外部依赖已提供稳定实现 | 优先复用依赖 |
| 都没有 | 新增最小实现，并放到正确层级 |

```go
// ✅ 先复用 internal/pkg
if err := internalpkg.StructConvert(in, &out); err != nil {
    return err
}
```

```go
// ❌ 业务目录重复造轮子
func structToMap(v any) (map[string]any, error) {
    ...
}
```

---

## 组合场景

```go
type AccountListener struct {
    producer *amqp.Producer
}

func (l *AccountListener) Handle(ctx context.Context, account *biz.Account) error {
    var data map[string]interface{}
    if err := internalpkg.StructConvert(account, &data); err != nil {
        return err
    }
    return l.producer.Publish(&amqp.Event{
        Event: openenum.AmqpEventAccountOpened.Value(),
        Data:  data,
    })
}
```

这个组合场景同时满足：

- 不对 Wire 注入依赖判空
- 不补查 `account` 的 relation
- 不为单次发布逻辑抽 helper
- 复用 `internal/pkg`
- 失败直接返回 `error`

---

## 常见错误模式

```go
// ❌ 注入依赖判空
if u.repo == nil { ... }
```

```go
// ❌ 主流程后置补查 relation
if order.User == nil { ... }
```

```go
// ❌ 单点逻辑抽 helper
return u.handleSyncResult(...)
```

```go
// ❌ 吞错
_ = err
```

```go
// ❌ 伪成功返回
return nil, nil
```

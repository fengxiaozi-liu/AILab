# Seata Spec

## 全局事务配置决策

| 配置项 | 要求 |
|--------|------|
| `Name` | 必填，表示事务语义 |
| `Timeout` | 必填，结合业务链路时延合理设置 |
| `Propagation` | 默认 `REQUIRED`，嵌套事务用 `REQUIRES_NEW` |
| `LockRetryInternal / LockRetryTimes` | 按竞争情况设置，不省略 |

---

## 标准模板

```go
// ✅ GtxConfig 显式配置，通过 WithGlobalTx 包装业务回调
type GtxConfig struct {
    Timeout           time.Duration
    Name              string
    Propagation       tm.Propagation
    LockRetryInternal time.Duration
    LockRetryTimes    int16
}

func WithGlobalTx(ctx context.Context, gc *GtxConfig, business CallbackWithCtx) error {
    // begin → callback → commit / rollback
    return nil
}

// 调用方
err := seata.WithGlobalTx(ctx, &seata.GtxConfig{
    Timeout:           30 * time.Second,
    Name:              "open_account",
    Propagation:       tm.Required,
    LockRetryInternal: 10 * time.Millisecond,
    LockRetryTimes:    3,
}, func(ctx context.Context) error {
    if err := u.accountRepo.Create(ctx, req); err != nil { return err }
    return u.accountCollectRepo.Create(ctx, req.Collect)
})
```

---

## 禁止在 helper 中加业务逻辑

```go
// ❌ seata helper 内含业务补偿
func WithGlobalTx(ctx context.Context, gc *GtxConfig, business CallbackWithCtx) error {
    err := business(ctx)
    if err != nil {
        compensateAccount(ctx)  // ❌ 补偿逻辑不属于事务 helper
    }
    return err
}

// ❌ seata helper 内做状态判断
func WithGlobalTx(ctx context.Context, gc *GtxConfig, business CallbackWithCtx) error {
    account, _ := accountRepo.GetAccount(ctx, gc.AccountID)  // ❌ helper 不访问业务
    if account.Status != "pending" { return errors.New("invalid status") }
    ...
}
```

---

## 组合场景

```go
// 完整：UseCase 使用 seata.WithGlobalTx 包装多步写操作
func (u *AccountUseCase) OpenAccount(ctx context.Context, req *OpenAccountReq) error {
    return seata.WithGlobalTx(ctx, &seata.GtxConfig{
        Timeout:           60 * time.Second,
        Name:              "open_account_global_tx",
        Propagation:       tm.Required,
        LockRetryInternal: 10 * time.Millisecond,
        LockRetryTimes:    5,
    }, func(ctx context.Context) error {
        // 本地 DB 写入
        if err := u.accountRepo.Create(ctx, req); err != nil { return err }
        // 跨服务写入（通过 Depend 发起）
        if err := u.capitalDepend.OpenCapitalAccount(ctx, req.CapitalReq); err != nil { return err }
        // 发布事件
        return u.eventBus.Publish(ctx, &eventbus.Event{
            Topic:   openenum.LocalEventAccountAfterOpen.Value(),
            Payload: req,
        })
    })
}
```

---

## 常见错误模式

```go
// ❌ 每个模块复制一份全局事务模板
// service_a/pkg/seata.go  service_b/pkg/seata.go  ← 各自维护不一致

// ❌ GtxConfig 缺少必填字段
seata.WithGlobalTx(ctx, &seata.GtxConfig{}, func(ctx context.Context) error { ... })
// ❌ Name 和 Timeout 为零值，难以排查和监控

// ❌ 在 Repo 内开启全局事务（事务边界应在 UseCase）
func (r *accountRepo) CreateWithTx(ctx context.Context, req *biz.OpenAccountReq) error {
    return seata.WithGlobalTx(ctx, &seata.GtxConfig{...}, func(ctx context.Context) error {  // ❌
        return r.data.Db.Account(ctx).Create().SetUserCode(req.UserCode).Exec(ctx)
    })
}
```

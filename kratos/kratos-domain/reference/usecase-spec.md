# UseCase Spec

## 职责归属决策

| 逻辑 | 归属 |
|------|------|
| 状态校验、权限判断 | UseCase |
| 事务边界管理 | UseCase |
| 关联对象的查询装配 | Repo（通过 opts 委托） |
| 跨聚合根数据补充 | Repo.serviceRelation |
| 协议字段转换 | Service |
| 数据库 CRUD | Repo |

---

## opts 委托 relation 查询

```go
// ✅ opts 委托：UseCase 通知 Repo 加载哪些 relation
func (u *AccountUseCase) GetAccountDetail(ctx context.Context, id uint32) (*Account, error) {
    return u.accountRepo.GetAccount(ctx, id,
        filter.WithRelation(openenum.AccountCollectRelation),
        filter.WithRelation(openenum.AccountCheckUserRelation),
    )
}

// ❌ UseCase 自行调用 Repo/Depend 补查（破坏四段式）
func (u *AccountUseCase) GetAccountDetail(ctx context.Context, id uint32) (*Account, error) {
    account, err := u.accountRepo.GetAccount(ctx, id)
    if err != nil { return nil, err }
    user, err := u.adminUserDepend.GetUser(ctx, account.CheckAdminUserID)  // ❌ UseCase 补查
    account.CheckAdminUserInfo = user
    return account, nil
}
```

---

## 事务边界

```go
// ✅ 事务在 UseCase 管理
func (u *AccountUseCase) OpenAccount(ctx context.Context, req *OpenAccountReq) error {
    return u.data.WithTx(ctx, func(ctx context.Context) error {
        if err := u.accountRepo.Create(ctx, req); err != nil { return err }
        if err := u.accountCollectRepo.Create(ctx, req.Collect); err != nil { return err }
        return u.eventBus.Publish(ctx, &eventbus.Event{
            Topic:   openenum.LocalEventAccountAfterOpen.Value(),
            Payload: req,
        })
    })
}

// ❌ 事务在 Repo 中开启（Repo 不应管理事务边界）
func (r *accountRepo) Create(ctx context.Context, req *biz.OpenAccountReq) error {
    return r.data.WithTx(ctx, func(ctx context.Context) error { ... })  // ❌
}

// ⚠️ 只写一个表：不需要事务
func (u *AccountUseCase) UpdateStatus(ctx context.Context, id uint32, status openenum.AccountOpenStatus) error {
    return u.accountRepo.UpdateStatus(ctx, id, status)  // ✅ 单操作不需要 tx
}
```

---

## 状态流转守卫

```go
// ✅ UseCase 校验前置状态
func (u *AccountUseCase) Review(ctx context.Context, id uint32, pass bool) error {
    account, err := u.accountRepo.GetAccount(ctx, id)
    if err != nil { return err }
    if account.OpenStatus != openenum.AccountOpenStatusPending {
        return openerror.ErrorAccountStatusInvalid(ctx)
    }
    newStatus := openenum.AccountOpenStatusRejected
    if pass { newStatus = openenum.AccountOpenStatusOpened }
    return u.accountRepo.UpdateStatus(ctx, id, newStatus)
}

// ❌ 跳过状态校验直接更新
func (u *AccountUseCase) Review(ctx context.Context, id uint32, pass bool) error {
    return u.accountRepo.UpdateStatus(ctx, id, openenum.AccountOpenStatusOpened)  // ❌ 无守卫
}
```

---

## 组合场景

```go
// 完整 UseCase：包含状态守卫 + 事务 + ops 委托
func (u *AccountUseCase) PassReview(ctx context.Context, id uint32) error {
    // 1. 加载（不含 relation，仅需状态字段）
    account, err := u.accountRepo.GetAccount(ctx, id)
    if err != nil { return err }

    // 2. 状态守卫
    if account.OpenStatus != openenum.AccountOpenStatusPending {
        return openerror.ErrorAccountStatusInvalid(ctx)
    }

    // 3. 事务：更新状态 + 发送事件
    return u.data.WithTx(ctx, func(ctx context.Context) error {
        if err := u.accountRepo.UpdateStatus(ctx, id, openenum.AccountOpenStatusOpened); err != nil {
            return err
        }
        return u.eventBus.Publish(ctx, &eventbus.Event{
            Topic:   openenum.LocalEventAccountAfterOpen.Value(),
            Payload: account,  // ✅ 直接传聚合对象
        })
    })
}
```

---

## 常见错误模式

```go
// ❌ UseCase 构造查询直接访问 DB
func (u *AccountUseCase) GetAccount(ctx context.Context, id uint32) (*Account, error) {
    return u.data.Db.Account(ctx).Query().Where(entaccount.IDEQ(id)).First(ctx)
}

// ❌ UseCase 手动补查 relation（应通过 opts 委托给 Repo）
func (u *AccountUseCase) GetDetail(ctx context.Context, id uint32) (*Account, error) {
    account, _ := u.accountRepo.GetAccount(ctx, id)
    account.Collect, _ = u.accountCollectRepo.GetByAccountID(ctx, id)  // ❌
    return account, nil
}

// ❌ 业务重复查询（一次 Get 后又 Get 一次）
func (u *AccountUseCase) Review(ctx context.Context, id uint32, pass bool) error {
    if _, err := u.accountRepo.GetAccount(ctx, id); err != nil { return err }  // 第一次
    account, err := u.accountRepo.GetAccount(ctx, id)  // ❌ 重复
    ...
}
```

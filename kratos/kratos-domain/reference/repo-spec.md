# Repo Spec

## 四段式 Repo 骨架

```go
// ✅ 标准四段式结构
func (r *accountRepo) GetAccountList(ctx context.Context, opts ...filter.Option) ([]*biz.Account, error) {
    // 1. parseFilter - 解析 opts 构建过滤条件
    filterOpts := filter.ParseOptions(opts...)
    
    // 2. queryConfig - 构建 Ent 查询（Where / Order / Limit / WithEdges）
    query := r.data.Db.Account(ctx).Query()
    query = r.queryConfig(query, filterOpts)
    
    // 3. queryRelation - 装配 Ent Edge 关联（来自 WithEdges 预加载）
    list, err := query.All(ctx)
    if err != nil { return nil, err }
    result := make([]*biz.Account, 0, len(list))
    for _, item := range list {
        result = append(result, r.queryRelation(accountConvert(item), item.Edges))
    }
    
    // 4. serviceRelation - 通过 Depend/RPC 补查跨服务关联
    if err := r.serviceRelation(ctx, result, filterOpts); err != nil { return nil, err }
    return result, nil
}
```

---

## Update 方式决策

| 条件 | 方式 |
|------|------|
| 已有对象，单主键更新 | `UpdateOneID(id)` |
| 批量更新 / 无对象实例 / 条件更新 | `Update().Where(IDEQ(id))` |

```go
// ✅ 单主键更新，已知 ID
return r.data.Db.Account(ctx).UpdateOneID(id).
    SetStatus(openenum.AccountStatusOpened).
    SetUpdateTime(updateTime).
    Exec(ctx)

// ✅ 条件更新，无对象实例
return r.data.Db.Account(ctx).Update().
    Where(entaccount.IDEQ(id), entaccount.StatusEQ(openenum.AccountStatusPending)).
    SetStatus(openenum.AccountStatusOpened).
    Exec(ctx)

// ❌ 有主键时仍用 Where 更新，可读性差
r.data.Db.Account(ctx).Update().
    Where(entaccount.IDEQ(account.ID)).  // 已有 ID 时用 UpdateOneID 更清晰
    SetStatus(account.Status).Exec(ctx)
```

---

## Not Found 处理

```go
// ✅ not found 返回具体业务错误
info, err := query.First(ctx)
if err != nil {
    if ent.IsNotFound(err) {
        return nil, openerror.ErrorAccountNotFound(ctx)
    }
    return nil, err
}

// ❌ not found 返回 nil, nil（调用方无法区分是否存在）
info, err := query.First(ctx)
if ent.IsNotFound(err) {
    return nil, nil  // ❌ 禁止 nil-nil
}

// ⚠️ 列表查询，空集合返回空切片，不视为 not found
list, err := query.All(ctx)
if err != nil { return nil, err }
// len(list) == 0 是合法的空列表，不报错
```

---

## queryConfig / queryRelation / serviceRelation 私有方法

```go
// ✅ queryConfig：集中管理查询参数，保持 GetXxx/PageListXxx/ListXxx 一致
func (r *accountRepo) queryConfig(query *ent.AccountQuery, opts *filter.Options) *ent.AccountQuery {
    if opts.HasStatus() {
        query = query.Where(entaccount.StatusIn(opts.Statuses()...))
    }
    if opts.HasRelation(openenum.AccountCollectRelation) {
        query = query.WithAccountCollect()
    }
    return query
}

// ✅ serviceRelation：批量补查外部服务，避免 N+1
func (r *accountRepo) serviceRelation(ctx context.Context, list []*biz.Account, opts *filter.Options) error {
    if opts.HasRelation(openenum.AccountCheckUserRelation) {
        return r.adminUserDepend.FillUsers(ctx, list)
    }
    return nil
}
```

---

## 组合场景

```go
// 完整：PageListAccount 含 queryConfig + queryRelation + serviceRelation
func (r *accountRepo) PageListAccount(
    ctx context.Context, pg *page.Page, opts ...filter.Option,
) ([]*biz.Account, error) {
    filterOpts := filter.ParseOptions(opts...)

    query := r.data.Db.Account(ctx).Query()
    query = r.queryConfig(query, filterOpts)

    total, err := query.Count(ctx)
    if err != nil { return nil, err }
    pg.SetTotal(total)

    list, err := query.Offset(pg.Offset()).Limit(pg.Limit()).All(ctx)
    if err != nil { return nil, err }

    result := make([]*biz.Account, 0, len(list))
    for _, item := range list {
        result = append(result, r.queryRelation(accountConvert(item), item.Edges))
    }

    if err := r.serviceRelation(ctx, result, filterOpts); err != nil { return nil, err }
    return result, nil
}
```

---

## 常见错误模式

```go
// ❌ for 循环中 N+1 调用
for _, account := range accounts {
    user, _ := r.adminUserDepend.GetUser(ctx, account.CheckAdminUserID)  // ❌ N+1
    account.CheckAdminUserInfo = user
}

// ❌ nil-nil 隐藏 not found
info, err := query.First(ctx)
if ent.IsNotFound(err) { return nil, nil }  // ❌

// ❌ Service 或 UseCase 直接拼 Ent 查询
func (u *AccountUseCase) GetAccount(ctx context.Context, id uint32) (*Account, error) {
    return u.data.Db.Account(ctx).Query().Where(...).First(ctx)  // ❌ 越层
}
```

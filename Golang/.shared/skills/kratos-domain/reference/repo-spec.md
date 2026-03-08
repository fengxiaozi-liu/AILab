# Repo Reference

## 这个主题解决什么问题

说明 Repo 如何围绕聚合根组织查询、relation 装配和跨服务补全，以及常见的实现模板。

## 适用场景

- 新增或改造 Repo 接口
- 实现列表、分页、详情查询
- 区分本地关联和跨服务 relation

## 设计意图

Repo 参考的重点是解释查询、关系装配和跨服务补全为什么要组织成稳定骨架，而不是把每个接口都写成独立查询脚本。

- 四段式结构让过滤、查询配置、本地 relation 和远程 relation 各自演进。
- 列表和详情共享同一查询骨架后，更容易复用现有实现，而不是重新拼一套逻辑。
- relation 收口在 Repo 后，UseCase 可以持续保持业务编排角色，不需要关心装配细节。

## 实施提示

- 先画出“过滤 -> 查询配置 -> 本地 relation -> 远程 relation”的数据流。
- 优先复用已有聚合的查询骨架，再补本次差异化字段和 relation。
- 看到列表和详情在重复装配同一关系时，优先提炼共享装配函数。

## 推荐结构

Repo 查询通常按四段组织：

1. `parseFilter(...)`
2. `queryConfig(...)`
3. `queryRelation(...)`
4. `serviceRelation(...)`

## 典型实现方式

### 单对象查询

```go
func (r *accountRepo) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*biz.Account, error) {
    query := r.data.Db.Account(ctx).Query().Where(entaccount.IDEQ(id))
    query = r.parseFilter(query, nil)
    query = r.queryConfig(query, opts...)

    info, err := query.First(ctx)
    if err != nil {
        return nil, err
    }

    res := r.queryRelation(accountConvert(info), info.Edges)
    if err := r.serviceRelation(ctx, res, opts...); err != nil {
        return nil, err
    }

    return res, nil
}
```

### 列表查询

```go
func (r *accountRepo) ListAccount(ctx context.Context, f *biz.AccountFilter, opts ...filter.Option) ([]*biz.Account, error) {
    query := r.data.Db.Account(ctx).Query()
    query = r.parseFilter(query, f)
    query = r.queryConfig(query, opts...)

    list, err := query.All(ctx)
    if err != nil {
        return nil, err
    }

    res := make([]*biz.Account, 0, len(list))
    for _, item := range list {
        res = append(res, r.queryRelation(accountConvert(item), item.Edges))
    }

    if err := r.serviceRelation(ctx, res, opts...); err != nil {
        return nil, err
    }

    return res, nil
}
```

## 本地关联装配

```go
func (r *accountRepo) queryConfig(query *ent.AccountQuery, opts ...filter.Option) *ent.AccountQuery {
    cfg := filter.NewConfig(opts...)
    if _, ok := cfg.Relations[openenum.AccountCollectRelation]; ok {
        query = query.WithAccountCollect()
    }
    return query
}

func (r *accountRepo) queryRelation(info *biz.Account, edges ent.AccountEdges) *biz.Account {
    if info == nil {
        return nil
    }
    if len(edges.AccountCollect) > 0 {
        info.Collect = accountCollectConvert(edges.AccountCollect[0])
    }
    return info
}
```

## 项目通用完整骨架

```go
func (r *accountFlowPageRepo) GetAccountFlowPageByFilter(ctx context.Context, f *openbiz.AccountFlowPageFilter, opts ...filter.Option) (*openbiz.AccountFlowPage, error) {
    query := r.data.Db.AccountFlowPage(ctx).Query()
    query = r.parseFilter(query, f)
    query = r.queryConfig(query, opts...)

    info, err := query.First(ctx)
    if err != nil {
        if ent.IsNotFound(err) {
            return nil, openerror.ErrorPageNotFound(ctx)
        }
        return nil, err
    }

    result := r.queryRelation(flowPageConvert(info), info.Edges)
    if err := r.serviceRelation(ctx, result, opts...); err != nil {
        return nil, err
    }
    return result, nil
}

func (r *accountFlowPageRepo) ListAccountFlowPage(ctx context.Context, f *openbiz.AccountFlowPageFilter, opts ...filter.Option) ([]*openbiz.AccountFlowPage, error) {
    query := r.data.Db.AccountFlowPage(ctx).Query()
    query = r.parseFilter(query, f)
    query = r.queryConfig(query, opts...)

    list, err := query.All(ctx)
    if err != nil {
        return nil, err
    }

    result := make([]*openbiz.AccountFlowPage, 0, len(list))
    for _, info := range list {
        result = append(result, r.queryRelation(flowPageConvert(info), info.Edges))
    }

    if err := r.serviceRelation(ctx, result, opts...); err != nil {
        return nil, err
    }
    return result, nil
}
```

## 跨服务 relation 装配

```go
func (r *accountRepo) serviceRelation(ctx context.Context, data interface{}, opts ...filter.Option) error {
    cfg := filter.NewConfig(opts...)
    if _, ok := cfg.Relations[openenum.AccountCheckUserRelation]; !ok {
        return nil
    }

    if info, ok := data.(*biz.Account); ok {
        data = []*biz.Account{info}
    }
    list := data.([]*biz.Account)

    idSet := make(map[uint32]struct{})
    for _, item := range list {
        if item == nil {
            continue
        }
        if item.FirstCheckUserID > 0 {
            idSet[item.FirstCheckUserID] = struct{}{}
        }
        if item.SecondCheckUserID > 0 {
            idSet[item.SecondCheckUserID] = struct{}{}
        }
    }

    ids := make([]uint32, 0, len(idSet))
    for id := range idSet {
        ids = append(ids, id)
    }

    userMap, err := r.adminUserRepo.MapAdminUser(ctx, &adminbiz.AdminUserFilter{IDList: ids})
    if err != nil {
        return err
    }

    for _, item := range list {
        if item == nil {
            continue
        }
        item.FirstCheckUser = userMap[item.FirstCheckUserID]
        item.SecondCheckUser = userMap[item.SecondCheckUserID]
    }
    return nil
}
```

## Upsert 示例

```go
func (r *accountFlowPageRepo) UpsertAccountFlowPage(ctx context.Context, info *openbiz.AccountFlowPage) error {
    nowTime := uint32(time.Now().Unix())

    return r.data.Db.AccountFlowPage(ctx).
        Create().
        SetAccountID(info.AccountID).
        SetPageCode(info.PageCode).
        SetStatus(info.Status).
        SetReviewStage(info.ReviewStage).
        SetReasonText(info.ReasonText).
        SetCreateTime(nowTime).
        SetUpdateTime(nowTime).
        OnConflict().
        SetStatus(info.Status).
        SetReviewStage(info.ReviewStage).
        SetReasonText(info.ReasonText).
        SetUpdateTime(nowTime).
        Exec(ctx)
}
```

## 常见坑

- 本地 edge 和跨服务 relation 混在一个阶段处理
- 详情和列表各写一套装配逻辑
- 为某个接口单独返回临时 DTO，导致模型漂移
- 更新后立刻取回实体，但实际上并不需要更新后的完整对象

## 更新写法示例

```go
func (r *accountRepo) UpdateAccountStatus(ctx context.Context, id uint32, status openenum.AccountStatus) error {
    return r.data.Db.Account(ctx).
        Update().
        Where(entaccount.IDEQ(id)).
        SetStatus(status).
        Exec(ctx)
}
```

相比 `UpdateOneID(id)`，这类写法更适合作为默认更新模板，尤其是在只需要按条件更新并直接执行的 Repo 场景里。

反例：在单个 Repo 方法中删除整个聚合。

```go
func (r *accountRepo) DeleteAccountAggregate(ctx context.Context, id uint32) error {
    _, _ = r.data.Db.AccountCollect(ctx).Delete().Where(accountcollect.AccountID(id)).Exec(ctx)
    _, _ = r.data.Db.Account(ctx).Delete().Where(account.ID(id)).Exec(ctx)
    return nil
}
```

正例：Repo 只提供原子删除方法，由 UseCase 组合调用。

```go
func (r *accountRepo) DeleteAccount(ctx context.Context, id uint32) error {
    affected, err := r.data.Db.Account(ctx).
        Delete().
        Where(account.ID(id)).
        Exec(ctx)
    if err != nil {
        return err
    }
    if affected == 0 {
        return openerror.ErrorAccountNotFound(ctx)
    }
    return nil
}
```

```go
func (r *accountCollectRepo) DeleteByAccountID(ctx context.Context, id uint32) error {
    _, err := r.data.Db.AccountCollect(ctx).
        Delete().
        Where(accountcollect.AccountID(id)).
        Exec(ctx)
    return err
}
```

## 相关 Rule

- `../rules/repo-rule.md`
- `../rules/domain-security-rule.md`

## 相关 Reference

- `./usecase-spec.md`
- `../../kratos-components/reference/depend-spec.md`

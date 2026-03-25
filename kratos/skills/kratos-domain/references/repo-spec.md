# Repo Reference

## 约束先看
必须遵守：

- Repo 负责查询组织、relation 装配和数据访问边界
- relation 收口在 Repo，Service 和 UseCase 不补查
- 避免 N+1，远程 relation 必须批量收集、批量查询、批量回填
- Repo 不承载整聚合更新编排和业务状态流转

## 使用说明

说明 Repo 如何围绕聚合根组织查询、relation 装配和跨服务补全，以及常见实现骨架。

## 常见场景

- 新增或改造 Repo 接口
- 实现列表、分页、详情查询
- 区分本地 relation 和跨服务 relation
- 设计单条更新、upsert、删除的写入边界

## 推荐骨架

Repo 查询通常按四段组织：

1. `parseFilter(...)`
2. `queryConfig(...)`
3. `queryRelation(...)`
4. `serviceRelation(...)`

## 四段职责

| 阶段 | 作用 | 典型内容 |
|------|------|------|
| `parseFilter` | 把业务过滤条件收敛成查询条件 | `IDList`、状态、时间范围、排序、分页 |
| `queryConfig` | 根据 `opts ...filter.Option` 决定是否预加载本地 relation | `WithAccountCollect()`、`WithFields()` |
| `queryRelation` | 把 ent edges 转成业务对象字段 | `CollectInfo`、`PageList`、`FieldList` |
| `serviceRelation` | 补远程 relation 或依赖侧数据 | 管理员信息、用户信息、其他 depend 数据 |

## `parseFilter` 如何使用

`parseFilter` 只做一件事：把业务过滤对象翻译成 Ent Query 的 where / sort / paging，不负责 relation，不负责远程补全。

典型模板：

```go
func (r *accountRepo) parseFilter(query *ent.AccountQuery, f *openbiz.AccountFilter) *ent.AccountQuery {
    if f == nil {
        return query
    }
    if len(f.IDList) > 0 {
        query = query.Where(account.IDIn(f.IDList...))
    }
    if f.UserCode != "" {
        query = query.Where(account.UserCode(f.UserCode))
    }
    if len(f.OpenStatusList) > 0 {
        query = query.Where(account.OpenStatusIn(f.OpenStatusList...))
    }
    query.Modify(f.Sort.ModifyFn(account.ValidColumn), f.Paging.ModifyFn())
    return query
}
```

约束：

- `parseFilter` 输入是业务 filter，不是协议请求对象
- 只改 query 条件，不读取远程依赖
- 排序、分页通常在这里统一收口

## 关联配置如何定义

仓库里真实做法不是写裸字符串，而是先定义 `baseenum.Relation`，再通过 `filter.WithRelation(...)` 传入 Repo。

基础类型：

```go
type Relation string
```

业务 relation 常量示例：

```go
const (
    AccountCollectRelation       baseenum.Relation = "account_collect_relation"
    AccountCheckUserRelation     baseenum.Relation = "account_check_user_relation"
    AccountFlowPageRelation      baseenum.Relation = "account_flow_page_relation"
    AccountFlowPageFieldRelation baseenum.Relation = "account_flow_page_field_relation"
)
```

调用侧示例：

```go
account, err := srv.accountUseCase.GetAccountByID(
    ctx,
    req.Id,
    filter.WithRelation(openenum.AccountCollectRelation),
)
```

```go
pages, err := srv.accountFlowPageUseCase.ListAccountFlowPage(
    ctx,
    req.AccountId,
    filter.WithRelation(openenum.AccountFlowPageFieldRelation),
)
```

## `queryConfig` 如何使用

`queryConfig` 根据 `opts ...filter.Option` 解析 relation 配置，并决定是否在 Ent 查询阶段预加载本地 edges。

典型模板：

```go
func (r *accountRepo) queryConfig(query *ent.AccountQuery, f *openbiz.AccountFilter, opts ...filter.Option) *ent.AccountQuery {
    cfg := filter.NewConfig(opts...)
    if _, ok := cfg.Relations[openenum.AccountCollectRelation]; ok {
        query = query.WithAccountCollect()
    }
    if _, ok := cfg.Relations[openenum.AccountFlowPageRelation]; ok {
        if f != nil && f.PageStatus > 0 {
            query = query.WithAccountFlowPage(func(q *ent.AccountFlowPageQuery) {
                q.Where(accountflowpage.StatusEQ(f.PageStatus))
            })
        } else {
            query = query.WithAccountFlowPage()
        }
    }
    return query
}
```

约束：

- 本地 relation 在 `queryConfig` 收口
- relation 是否加载由 `filter.WithRelation(...)` 决定
- 需要 relation 局部过滤时，可结合 filter 参数一起配置

## `queryRelation` 如何使用

`queryRelation` 负责把 ent edges 转成业务对象字段，不做远程补查。

典型模板：

```go
func (r *accountRepo) queryRelation(info *openbiz.Account, edges ent.AccountEdges) *openbiz.Account {
    if info == nil {
        return nil
    }
    info.PageList = []*openbiz.AccountFlowPage{}
    if len(edges.AccountCollect) > 0 {
        info.CollectInfo = collectConvert(edges.AccountCollect[0])
    }
    if len(edges.AccountFlowPage) > 0 {
        info.PageList = make([]*openbiz.AccountFlowPage, 0, len(edges.AccountFlowPage))
        for _, page := range edges.AccountFlowPage {
            if page == nil {
                continue
            }
            info.PageList = append(info.PageList, &openbiz.AccountFlowPage{
                PageCode:   page.PageCode,
                Status:     page.Status,
                ReasonText: page.ReasonText,
            })
        }
    }
    return info
}
```

## `serviceRelation` 如何使用

`serviceRelation` 负责补远程 relation。只有命中 relation 枚举时，才去 depend repo 批量补数据。

典型模板：

```go
func (r *accountRepo) serviceRelation(ctx context.Context, data interface{}, opts ...filter.Option) error {
    cfg := filter.NewConfig(opts...)
    if _, ok := cfg.Relations[openenum.AccountCheckUserRelation]; !ok || r.adminUserRepo == nil {
        return nil
    }

    list := helper.SliceNormalize[*openbiz.Account](data)
    if len(list) == 0 {
        return nil
    }

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

    idList := make([]uint32, 0, len(idSet))
    for id := range idSet {
        idList = append(idList, id)
    }

    userMap, err := r.adminUserRepo.MapAdminUser(ctx, &adminbiz.AdminUserFilter{IDList: idList})
    if err != nil {
        return err
    }

    for _, item := range list {
        if item == nil {
            continue
        }
        item.FirstCheckUserInfo = userMap[item.FirstCheckUserID]
        item.SecondCheckUserInfo = userMap[item.SecondCheckUserID]
    }
    return nil
}
```

约束：

- `serviceRelation` 只补远程 relation，不做主查询
- 必须批量收集、批量查询、批量回填
- `data interface{}` 是为了兼容单对象和列表，仓库里通常配合 `helper.SliceNormalize` 使用

## 完整查询模板

### 单对象查询

```go
func (r *accountRepo) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*openbiz.Account, error) {
    query := r.data.Db.Account(ctx).Query()
    f := &openbiz.AccountFilter{IDList: []uint32{id}}
    query = r.parseFilter(query, f)
    query = r.queryConfig(query, f, opts...)

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
func (r *accountRepo) PageListAccount(ctx context.Context, f *openbiz.AccountFilter, opts ...filter.Option) ([]*openbiz.Account, int, error) {
    query := r.data.Db.Account(ctx).Query()
    query = r.parseFilter(query, f)
    query = r.queryConfig(query, f, opts...)

    count, err := query.Clone().Count(ctx)
    if err != nil {
        return nil, 0, err
    }

    list, err := query.All(ctx)
    if err != nil {
        return nil, 0, err
    }

    result := make([]*openbiz.Account, 0, len(list))
    for _, info := range list {
        result = append(result, r.queryRelation(accountConvert(info), info.Edges))
    }

    if err := r.serviceRelation(ctx, result, opts...); err != nil {
        return nil, 0, err
    }
    return result, count, nil
}
```

## 常见坑

- 在 UseCase 里手写 relation 查询
- 在 Repo 文档里只写调用，不定义 relation 常量和用法
- 把 `queryRelation` 和 `serviceRelation` 混成同一层职责
- 单个查询路径引入 N+1
- 把普通参数防御式校验写成 Repo 默认策略

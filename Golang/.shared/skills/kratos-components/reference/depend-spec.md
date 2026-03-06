# Depend Reference

## 这个主题解决什么问题

说明跨服务依赖如何封装为 depend 层，并在 Repo 或组件侧完成批量查询与回填。

## 适用场景

- 新增 InnerRPC 依赖
- 统一上游服务调用方式
- 把散落调用收口到一层封装

## 设计意图

Depend 参考用于解释跨服务依赖为什么要被统一封装，而不是散落在 Repo 或 UseCase 的任意位置。

- 统一封装后，调用协议、超时、批量化和返回映射更容易复用。
- 在做 relation 补全时，可以先判断是否已有统一依赖入口，而不是直接拼 RPC 调用。
- 依赖层稳定后，后续做替换、降级或观测也会更直接。

## 实施提示

- 先整理调用方真正需要的远程数据形态，再设计依赖接口。
- 优先围绕批量获取和结果映射设计方法签名。
- 如果多个 Repo 在重复同一远程查询和 map 组装，通常适合抽到 depend 层。

## 推荐结构

- depend 层提供 `MapXxx`、`ListXxx` 一类批量接口
- Repo 的 `serviceRelation(...)` 使用 depend 进行批量装配

## 标准模板

```go
func (r *adminUserRepo) MapAdminUser(ctx context.Context, f *biz.AdminUserFilter) (map[uint32]*biz.AdminUser, error) {
    list, err := r.ListAdminUser(ctx, f)
    if err != nil {
        return nil, err
    }
    res := make(map[uint32]*biz.AdminUser, len(list))
    for _, item := range list {
        res[item.ID] = item
    }
    return res, nil
}
```

## 代码示例参考

```go
func (r *accountRepo) serviceRelation(ctx context.Context, data interface{}, opts ...filter.Option) error {
    list := data.([]*biz.Account)
    ids := make([]uint32, 0, len(list))
    for _, item := range list {
        if item != nil && item.FirstCheckUserID > 0 {
            ids = append(ids, item.FirstCheckUserID)
        }
    }
    userMap, err := r.adminUserRepo.MapAdminUser(ctx, &adminbiz.AdminUserFilter{IDList: ids})
    if err != nil {
        return err
    }
    for _, item := range list {
        item.FirstCheckUser = userMap[item.FirstCheckUserID]
    }
    return nil
}
```

## 项目通用 depend 接口示例

```go
type AdminUserRepo interface {
    GetAdminUser(ctx context.Context, id uint32, opts ...filter.Option) (*AdminUser, error)
    MapAdminUser(ctx context.Context, filter *AdminUserFilter, opts ...filter.Option) (map[uint32]*AdminUser, error)
}
```

## InnerRPC 封装示例

```go
type adminUserRepo struct {
    adminUserClient v1.AdminUserServiceClient
}

func (r *adminUserRepo) GetAdminUser(ctx context.Context, id uint32, opts ...filter.Option) (*adminbiz.AdminUser, error) {
    info, err := r.adminUserClient.GetAdminUser(ctx, &v1.GetAdminUserRequest{
        Id:           id,
        FilterConfig: proto.BuildFilterConfig(opts...),
    })
    if err != nil {
        return nil, err
    }
    return AdminUserConvert(info), nil
}

func (r *adminUserRepo) MapAdminUser(ctx context.Context, filter *adminbiz.AdminUserFilter, opts ...filter.Option) (map[uint32]*adminbiz.AdminUser, error) {
    reply, err := r.adminUserClient.MapAdminUser(ctx, &v1.MapAdminUserRequest{
        IdList:       filter.IDList,
        FilterConfig: proto.BuildFilterConfig(opts...),
    })
    if err != nil {
        return nil, err
    }

    result := make(map[uint32]*adminbiz.AdminUser, len(reply.Map))
    for id, item := range reply.Map {
        result[id] = AdminUserConvert(item)
    }
    return result, nil
}
```

## Repo 批量回填示例

```go
func (r *accountRepo) serviceRelation(ctx context.Context, data interface{}, opts ...filter.Option) error {
    cfg := filter.NewConfig(opts...)
    if _, ok := cfg.Relations[openenum.AccountCheckUserRelation]; !ok {
        return nil
    }

    list := data.([]*openbiz.Account)
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
        item.FirstCheckUserInfo = userMap[item.FirstCheckUserID]
        item.SecondCheckUserInfo = userMap[item.SecondCheckUserID]
    }
    return nil
}
```

## 常见坑

- 每个业务模块都自己封装一次同一个上游调用
- 只提供单对象查询接口，导致上层无法批量装配
- 上游返回值没有统一转换为稳定领域模型

## 相关 Rule

- `../rules/depend-rule.md`

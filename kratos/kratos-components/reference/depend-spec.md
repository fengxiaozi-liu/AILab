# Depend Spec

## 批量 vs 逐条调用

| 条件 | 做法 |
|------|------|
| 需要在列表中补全远程字段 | 批量收集 ID → 一次 Map 查询 → 回填 |
| 单个对象查询 | 允许单条 RPC，但必须设置超时 |
| 循环中逐条 RPC | 禁止 |

```go
// ✅ 批量收集 ID → Map 查询 → 回填
func (r *accountRepo) serviceRelation(ctx context.Context, data interface{}, opts ...filter.Option) error {
    list := data.([]*biz.Account)
    ids := make([]uint32, 0, len(list))
    for _, item := range list {
        if item.FirstCheckUserID > 0 {
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

// ❌ 循环中逐条 RPC（N+1）
for _, item := range list {
    user, err := r.adminUserRepo.GetAdminUser(ctx, item.FirstCheckUserID)
    item.FirstCheckUser = user
}
```

---

## Depend 接口设计

```go
// ✅ 成对提供 Get + Map，MapXxx 供批量回填使用
type AdminUserRepo interface {
    GetAdminUser(ctx context.Context, id uint32, opts ...filter.Option) (*AdminUser, error)
    MapAdminUser(ctx context.Context, filter *AdminUserFilter, opts ...filter.Option) (map[uint32]*AdminUser, error)
}

// ✅ InnerRPC 封装：超时由 context deadline 控制，不散落在各调用点
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
```

---

## 重试策略

| 场景 | 做法 |
|------|------|
| 幂等操作（查询、读取） | 允许重试，采用退避 |
| 非幂等操作（写入、扣减） | 禁止自动重试 |

---

## 组合场景

Repo serviceRelation 完整示例（带 relation 开关 + 批量回填）：

```go
func (r *accountRepo) serviceRelation(ctx context.Context, data interface{}, opts ...filter.Option) error {
    cfg := filter.NewConfig(opts...)
    if _, ok := cfg.Relations[openenum.AccountCheckUserRelation]; !ok {
        return nil
    }
    list := data.([]*biz.Account)
    ids := make([]uint32, 0, len(list))
    for _, item := range list {
        if item.FirstCheckUserID > 0 {
            ids = append(ids, item.FirstCheckUserID)
        }
        if item.SecondCheckUserID > 0 {
            ids = append(ids, item.SecondCheckUserID)
        }
    }
    if len(ids) == 0 {
        return nil
    }
    userMap, err := r.adminUserRepo.MapAdminUser(ctx, &adminbiz.AdminUserFilter{IDList: ids})
    if err != nil {
        return err
    }
    for _, item := range list {
        item.FirstCheckUser = userMap[item.FirstCheckUserID]
        item.SecondCheckUser = userMap[item.SecondCheckUserID]
    }
    return nil
}
```

---

## 常见错误模式

```go
// ❌ 不经 Depend 封装，直连上游 gRPC Client
reply, err := r.rawGRPCClient.GetAdminUser(ctx, &pb.GetAdminUserRequest{Id: id})

// ❌ 循环逐条调用
for _, item := range list {
    user, _ := r.adminUserRepo.GetAdminUser(ctx, item.UserID)
    item.User = user
}

// ❌ 非幂等写操作自动重试
for i := 0; i < 3; i++ {
    err = r.chargeClient.Deduct(ctx, req)
    if err == nil { break }
}
```

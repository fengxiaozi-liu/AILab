# Layer Spec

## 各层职责决策

| 逻辑类型 | 落在哪层 |
|---------|---------|
| 协议字段转换、显式入参构造 | Service |
| 业务编排、权限、状态流转、事务边界 | UseCase |
| 查询过滤、relation 装配、数据访问 | Repo/Data |
| gRPC/HTTP 注册、中间件配置 | Server |

```go
// ✅ Service 只做参数转换，不做业务判断
func (s *AccountService) GetOpenStatus(ctx context.Context, req *v1.GetOpenStatusRequest) (*v1.GetOpenStatusReply, error) {
    userCode := strconv.FormatUint(uint64(metadata.GetViewerID(ctx)), 10)
    status, err := s.uc.GetOpenStatus(ctx, userCode)
    if err != nil { return nil, err }
    return &v1.GetOpenStatusReply{OpenStatus: uint32(status)}, nil
}

// ❌ Service 中实现业务判断
func (s *AccountService) GetOpenStatus(ctx context.Context, req *v1.GetOpenStatusRequest) (*v1.GetOpenStatusReply, error) {
    account, err := s.uc.GetAccount(ctx, req.Id)
    if account.Status != openenum.AccountStatusOpened {  // ❌ 业务判断不应在 Service
        return nil, errors.New(400, "NOT_OPEN", "未开户")
    }
    return &v1.GetOpenStatusReply{}, nil
}
```

---

## 禁止越层调用

```go
// ❌ Service 直接调用 Repo（越层）
type AccountService struct {
    repo biz.AccountRepo  // ❌ Service 不应持有 Repo
}

// ❌ UseCase 直接写 DB
func (u *AccountUseCase) GetAccount(ctx context.Context, id uint32) (*Account, error) {
    return u.data.Db.Account(ctx).Query().Where(...).First(ctx)  // ❌ UseCase 不直接访问 DB
}

// ✅ 标准调用链
// Request → Service → UseCase → Repo → DB/Depend
```

---

## relation 补查收口

```go
// ✅ relation 统一在 Repo 中装配
func (r *accountRepo) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*biz.Account, error) {
    query := r.data.Db.Account(ctx).Query().Where(entaccount.IDEQ(id))
    query = r.queryConfig(query, opts...)
    info, err := query.First(ctx)
    ...
    res := r.queryRelation(accountConvert(info), info.Edges)
    r.serviceRelation(ctx, res, opts...)
    return res, nil
}

// ❌ Service 级别补查 reviewer
func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    account, _ := s.uc.GetAccount(ctx, req.Id)
    reviewer, _ := s.adminUserRepo.GetAdminUser(ctx, account.ReviewerID)  // ❌ Service 补查 relation
    account.Reviewer = reviewer
    return convertToReply(account), nil
}
```

---

## 组合场景

```go
// 完整调用链：Service → UseCase → Repo
func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    info, err := s.uc.GetAccountDetail(ctx, req.Id)
    if err != nil { return nil, err }
    return &v1.GetAccountReply{Info: convertToProto(info)}, nil
}

func (u *AccountUseCase) GetAccountDetail(ctx context.Context, id uint32) (*Account, error) {
    opts := []filter.Option{
        filter.WithRelation(openenum.AccountCollectRelation),
        filter.WithRelation(openenum.AccountCheckUserRelation),
    }
    return u.accountRepo.GetAccount(ctx, id, opts...)
}

func (r *accountRepo) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*biz.Account, error) {
    query := r.data.Db.Account(ctx).Query().Where(entaccount.IDEQ(id))
    query = r.queryConfig(query, opts...)
    info, err := query.First(ctx)
    if err != nil {
        if ent.IsNotFound(err) { return nil, openerror.ErrorAccountNotFound(ctx) }
        return nil, err
    }
    res := r.queryRelation(accountConvert(info), info.Edges)
    if err := r.serviceRelation(ctx, res, opts...); err != nil { return nil, err }
    return res, nil
}
```

---

## 常见错误模式

```go
// ❌ Service 持有 Repo
type AccountService struct { repo biz.AccountRepo }

// ❌ Repo 处理状态流转
func (r *accountRepo) ApproveAccount(ctx context.Context, id uint32) error {
    account.Status = openenum.AccountStatusApproved  // 状态流转不应在 Repo
    return r.data.Db.Account(ctx).UpdateOneID(id).SetStatus(account.Status).Exec(ctx)
}

// ❌ UseCase 在 Service 之后再补查 relation
func (s *AccountService) Review(ctx context.Context, req) {
    _ = s.uc.Review(ctx, req.Id)
    account, _ := s.uc.GetAccount(ctx, req.Id)  // ❌ 重复查询
}
```

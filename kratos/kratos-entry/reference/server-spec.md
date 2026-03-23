# Server Spec

## Service 职责范围

| 允许在 Service 层做 | 禁止在 Service 层做 |
|--------------------|--------------------|
| proto 字段转换（`req.Id` → `id uint32`） | 业务逻辑、状态判断 |
| 元数据提取（`metadata.GetViewerID(ctx)`） | 调用 Repo 或直接访问 DB |
| 调用 UseCase 方法 | 发布 EventBus |
| 将 biz 对象转换为 proto reply | 控制事务边界 |

---

## Service 实现（协议适配器）

```go
// ✅ Service 只做 proto ↔ biz 转换，业务逻辑在 UseCase
type AccountService struct {
    v1.UnimplementedAccountServer
    uc *biz.AccountUseCase
}

func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    info, err := s.uc.GetAccount(ctx, req.Id)
    if err != nil { return nil, err }
    return &v1.GetAccountReply{Info: convertToProto(info)}, nil
}

// ❌ Service 做业务判断
func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    info, err := s.uc.GetAccount(ctx, req.Id)
    if err != nil { return nil, err }
    if info.Status != openenum.AccountStatusOpened {  // ❌ 业务判断不在 Service
        return nil, errors.New(400, "NOT_OPEN", "账户未开户")
    }
    return &v1.GetAccountReply{Info: convertToProto(info)}, nil
}
```

---

## Server 注册（统一入口）

```go
// ✅ 每个聚合根的注册在独立入口函数中
func NewHTTPServer(cfg *conf.Server, lis net.Listener, svc *service.AccountService) *http.Server {
    srv := http.NewServer(
        http.Listener(lis),
        http.Middleware(middleware.Server()...),
    )
    v1.RegisterAccountHTTPServer(srv, svc)
    return srv
}

func NewGRPCServer(cfg *conf.Server, lis net.Listener, svc *service.AccountService) *grpc.Server {
    srv := grpc.NewServer(
        grpc.Listener(lis),
        grpc.Middleware(middleware.Server()...),
    )
    v1.RegisterAccountServer(srv, svc)
    return srv
}

// ❌ 在多处散落注册
func someOtherFile(srv *http.Server, svc *service.AccountService) {
    v1.RegisterAccountHTTPServer(srv, svc)  // ❌ 注册散落在多处
}
```

---

## Service 构造函数

```go
// ✅ 构造函数只接收 UseCase 依赖
func NewAccountService(uc *biz.AccountUseCase) *AccountService {
    return &AccountService{uc: uc}
}

// ❌ Service 构造时注入 Repo
func NewAccountService(uc *biz.AccountUseCase, repo biz.AccountRepo) *AccountService {
    return &AccountService{uc: uc, repo: repo}  // ❌ Service 不应持有 Repo
}
```

---

## 组合场景

```go
// 完整：proto → Service → UseCase → proto reply
func (s *AccountService) PageListAccount(
    ctx context.Context, req *v1.PageListAccountRequest,
) (*v1.PageListAccountReply, error) {
    pg := page.New(req.Paging)
    sort := protohelper.ParseSort(req.Sort)
    filterOpts := protohelper.ParseFilterConfig(req.Filter)

    list, err := s.uc.PageListAccount(ctx, pg, sort, filterOpts...)
    if err != nil { return nil, err }

    items := make([]*v1.AccountItem, 0, len(list))
    for _, a := range list {
        items = append(items, convertToProtoItem(a))  // ✅ 转换在 Service
    }
    return &v1.PageListAccountReply{
        List:   items,
        Paging: protohelper.BuildPaging(pg),
    }, nil
}
```

---

## 常见错误模式

```go
// ❌ Service 直接查询 Repo
func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    info, _ := s.repo.GetAccount(ctx, req.Id)  // ❌ 越层
    return &v1.GetAccountReply{Info: convertToProto(info)}, nil
}

// ❌ Server 注册散落在多处业务文件
// biz/account_usecase.go 包含 v1.RegisterAccountHTTPServer(...)  ❌

// ❌ Service 方法过大（含大量 if/else 业务判断）
// 超过 30 行的 Service 方法，通常意味着业务逻辑未下沉到 UseCase
```

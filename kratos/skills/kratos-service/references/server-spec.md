# Server Spec

## Service 与业务边界

| 条件 | 做法 |
|------|------|
| service 层处理请求 | 做参数转换、调用 usecase、返回协议对象 |
| 出现状态编排、规则判断、跨 repo 组合 | 回到 usecase/repo 边界治理 |

```go
// ✅ service 只做协议适配
func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    account, err := s.uc.GetAccount(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return &v1.GetAccountReply{Info: convertAccount(account)}, nil
}
```

```go
// ❌ service 承担业务编排
func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    if req.Id == 0 {
        return nil, errors.New("invalid id")
    }
    account, err := s.accountRepo.Get(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    if account.Status == "init" {
        if err := s.auditRepo.Save(ctx, account.ID); err != nil {
            return nil, err
        }
    }
    return convert(account), nil
}
```

```go
// ⚠️ 轻量参数归一化允许留在 service，业务规则不允许
req.Page = max(req.Page, 1)
```

---

## 服务注册入口

| 条件 | 做法 |
|------|------|
| 新增 service | 在统一的 server 注册入口补齐 |
| 调整路由/协议 | 同步检查 HTTP 和 gRPC 注册 |

```go
// ✅ 统一入口注册
func NewGRPCServer(c *conf.Server, account *service.AccountService) *grpc.Server {
    srv := grpc.NewServer()
    v1.RegisterAccountServiceServer(srv, account)
    return srv
}
```

```go
// ❌ 散落注册
func init() {
    globalServer.Register(AccountServiceDesc, &AccountService{})
}
```

```go
// ⚠️ 只有项目明确拆分多入口时，才允许按 server 类型分别维护，但每类入口仍要统一
func NewHTTPServer(...) *http.Server { ... }
func NewGRPCServer(...) *grpc.Server { ... }
```

---

## 新增接口联动

| 条件 | 做法 |
|------|------|
| 新增 RPC | 同步检查 proto、service 实现、server 注册、provider/wire、codegen |
| 只改 service 内部映射 | 至少确认协议和注册不受影响 |

```text
// ✅ 新增接口联动顺序
api/*.proto
-> internal/service/account.go
-> internal/server/http.go / grpc.go
-> internal/service/service.go 或 ProviderSet
-> make generate / wire / go build ./...
```

```text
// ❌ 只补 service 方法，不补注册和生成
api/account.proto
internal/service/account.go
```

---

## 组合场景

```go
type AccountService struct {
    v1.UnimplementedAccountServiceServer
    uc *biz.AccountUseCase
}

func NewAccountService(uc *biz.AccountUseCase) *AccountService {
    return &AccountService{uc: uc}
}

func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    account, err := s.uc.GetAccount(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return &v1.GetAccountReply{Info: convertAccount(account)}, nil
}
```

```go
func NewGRPCServer(account *service.AccountService) *grpc.Server {
    srv := grpc.NewServer()
    v1.RegisterAccountServiceServer(srv, account)
    return srv
}
```

这个组合场景同时满足：

- service 只做协议适配
- provider 构造清晰
- 注册集中在统一入口
- 新增接口时容易定位联动点

---

## 常见错误模式

```go
// ❌ service 中直接调 repo
s.accountRepo.Get(...)
```

```go
// ❌ service 中维护状态机
switch order.Status { ... }
```

```go
// ❌ 注册散落在 init 或业务文件
func init() { ... }
```

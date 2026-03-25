# Gateway Spec

## 代理与业务边界

| 条件 | 做法 |
|------|------|
| gateway 请求转发 | 做参数映射、调用上游、响应适配 |
| 出现业务状态流转或复杂规则判断 | 回到下游服务或 usecase |

```go
// ✅ gateway 只做代理与转换
func (s *GatewayService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    resp, err := s.accountClient.GetAccount(ctx, &openv1.GetAccountRequest{Id: req.Id})
    if err != nil {
        return nil, err
    }
    return &v1.GetAccountReply{Info: convertReply(resp.Info)}, nil
}
```

```go
// ❌ gateway 承担业务状态管理
func (s *GatewayService) OpenAccount(ctx context.Context, req *v1.OpenAccountRequest) (*v1.OpenAccountReply, error) {
    if err := s.accountRepo.SaveDraft(ctx, req.Id); err != nil {
        return nil, err
    }
    ...
}
```

```go
// ⚠️ 聚合多个下游可以接受，但仍只做编排式聚合，不维护领域状态
func (s *GatewayService) GetOverview(ctx context.Context, req *v1.GetOverviewRequest) (*v1.GetOverviewReply, error) {
    account, err := s.accountClient.GetAccount(ctx, &openv1.GetAccountRequest{Id: req.Id})
    if err != nil {
        return nil, err
    }
    user, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: account.Info.UserId})
    if err != nil {
        return nil, err
    }
    return convertOverview(account, user), nil
}
```

---

## 参数映射与错误转换

| 条件 | 做法 |
|------|------|
| 协议差异在接入层可消化 | 在 gateway 做字段映射和错误转换 |
| 需要改变业务语义 | 回到上游服务改契约 |

```go
// ✅ 参数映射留在 gateway
resp, err := s.client.List(ctx, &openv1.ListRequest{
    PageNo:   req.Page,
    PageSize: req.PageSize,
})
```

```go
// ❌ 为适配前端强行改业务语义
if req.Status == "all" {
    req.Status = ""
    req.IncludeDeleted = true
}
```

---

## 路由与落位

| 条件 | 做法 |
|------|------|
| gateway service 提供对外接口 | 放在统一 gateway service/provider 下 |
| 代理链路变化 | 同步检查 HTTP 映射、下游 client 和 OpenAPI 展示 |

```text
// ✅
internal/service/gateway_account.go
api/admin/v1/gateway_account.proto
```

```text
// ❌
internal/server/gateway_logic.go
internal/service/tmp_proxy.go
```

---

## 组合场景

```go
type GatewayService struct {
    accountClient openv1.AccountServiceClient
    userClient    userv1.UserServiceClient
}

func (s *GatewayService) GetOverview(ctx context.Context, req *v1.GetOverviewRequest) (*v1.GetOverviewReply, error) {
    account, err := s.accountClient.GetAccount(ctx, &openv1.GetAccountRequest{Id: req.Id})
    if err != nil {
        return nil, err
    }
    user, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: account.Info.UserId})
    if err != nil {
        return nil, err
    }
    return &v1.GetOverviewReply{
        AccountName: account.Info.Name,
        UserName:    user.Info.Name,
    }, nil
}
```

这个组合场景同时满足：

- gateway 只做代理和聚合适配
- 没有直接访问 repo 或维护状态
- 参数映射和响应映射都留在接入层

---

## 常见错误模式

```go
// ❌ gateway 直接访问 repo
s.accountRepo.Get(...)
```

```go
// ❌ gateway 维护业务状态
s.stateMachine.Transfer(...)
```

```go
// ❌ 路由和代理逻辑散落在 server 文件
internal/server/http_gateway.go
```

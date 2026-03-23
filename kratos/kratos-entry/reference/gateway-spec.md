# Gateway Spec

## Gateway 职责范围

| 允许在 Gateway 层做 | 禁止在 Gateway 层做 |
|--------------------|---------------------|
| HTTP 路径参数提取 | 查询 DB / 调用 Repo |
| 协议转换（HTTP → gRPC req） | 承接业务流程（状态流转）|
| 调用上游 gRPC client | 聚合多个上游并做业务决策 |
| 响应适配（gRPC reply → HTTP resp） | 发布 EventBus |

---

## 标准代理模式

```go
// ✅ Gateway 只做参数提取 + 代理调用 + 响应适配
func (s *GatewayService) GetAccount(ctx http.Context) error {
    req := &v1.GetAccountRequest{
        Id: cast.ToUint32(ctx.Vars("id")),
    }
    reply, err := s.accountClient.GetAccount(ctx, req)
    if err != nil { return err }
    return ctx.Result(200, reply)  // ✅ 透传 reply，不重组
}

// ❌ Gateway 内做业务判断
func (s *GatewayService) GetAccount(ctx http.Context) error {
    req := &v1.GetAccountRequest{Id: cast.ToUint32(ctx.Vars("id"))}
    reply, err := s.accountClient.GetAccount(ctx, req)
    if err != nil { return err }
    if reply.Info.Status != "opened" {  // ❌ 业务判断不在 Gateway
        return errors.New(400, "NOT_OPEN", "未开户")
    }
    return ctx.Result(200, reply)
}
```

---

## 跨协议参数映射

```go
// ✅ 入参从 HTTP path/query/body 构造 gRPC req
func (s *GatewayService) PageListAccount(ctx http.Context) error {
    var req v1.PageListAccountRequest
    if err := ctx.BindQuery(&req); err != nil { return err }
    req.TenantId = cast.ToUint32(ctx.Header("X-Tenant-ID"))

    reply, err := s.accountClient.PageListAccount(ctx, &req)
    if err != nil { return err }
    return ctx.Result(200, reply)
}

// ❌ Gateway 内自行拼分页逻辑
func (s *GatewayService) PageListAccount(ctx http.Context) error {
    page := cast.ToInt(ctx.Query("page"))
    pageSize := cast.ToInt(ctx.Query("page_size"))
    offset := (page - 1) * pageSize  // ❌ 分页计算应在 proto helper 或上游处理
    ...
}
```

---

## Gateway 构造函数

```go
// ✅ 只注入上游 gRPC client
func NewGatewayService(accountClient v1.AccountClient) *GatewayService {
    return &GatewayService{accountClient: accountClient}
}

// ❌ Gateway 注入 Repo 或 UseCase
func NewGatewayService(accountClient v1.AccountClient, repo biz.AccountRepo) *GatewayService {
    return &GatewayService{accountClient: accountClient, repo: repo}  // ❌
}
```

---

## 组合场景

```go
// 完整：HTTP 请求 → Gateway 代理 → gRPC 上游 → HTTP 响应
func (s *GatewayService) RegisterRoutes(r *http.Router) {
    r.GET("/api/v1/accounts/{id}", s.GetAccount)
    r.POST("/api/v1/accounts/list", s.PageListAccount)
}

func (s *GatewayService) GetAccount(ctx http.Context) error {
    req := &v1.GetAccountRequest{Id: cast.ToUint32(ctx.Vars("id"))}
    reply, err := s.accountClient.GetAccount(ctx, req)  // ✅ 代理调用
    if err != nil { return err }
    return ctx.Result(200, reply)  // ✅ 透传响应
}

func (s *GatewayService) PageListAccount(ctx http.Context) error {
    var req v1.PageListAccountRequest
    if err := ctx.Bind(&req); err != nil { return err }
    reply, err := s.accountClient.PageListAccount(ctx, &req)
    if err != nil { return err }
    return ctx.Result(200, reply)
}
```

---

## 常见错误模式

```go
// ❌ Gateway 实现完整业务流程
func (s *GatewayService) OpenAccount(ctx http.Context) error {
    // 查数据库检查是否已开户 ❌
    account, _ := s.accountRepo.GetAccount(ctx, id)
    if account != nil { return errors.New(400, "EXISTS", "已开户") }
    _, err := s.accountClient.CreateAccount(ctx, &v1.CreateAccountRequest{...})
    return err
}

// ❌ Gateway 聚合多个上游并做决策
func (s *GatewayService) GetAccountDetail(ctx http.Context) error {
    account, _ := s.accountClient.GetAccount(ctx, req)
    store, _ := s.storeClient.GetStore(ctx, req)
    if account.Status == store.Status { ... }  // ❌ 跨服务业务决策
}
```

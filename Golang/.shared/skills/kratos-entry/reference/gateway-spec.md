# Gateway Reference

## 这个主题解决什么问题

说明 Gateway/Openapi 场景下如何做代理转发、参数映射和响应适配。

## 适用场景

- 新增或修改 Gateway 路由
- 设计代理请求和响应映射
- 处理协议适配

## 设计意图

Gateway 参考解释的是代理层如何负责协议转换和转发，而不是承接完整业务流程。

- Gateway 面向外部入口、路由整形和参数映射，它与业务核心的变化节奏不同。
- 理解这一点后，更容易把网关实现写成清晰代理层，而不是第二套业务服务。
- 代理层结构稳定时，排查外部请求问题会更直接。

## 实施提示

- 先说明请求从哪里来、要转发到哪里、参数如何映射。
- 再决定是直接代理、聚合多个下游，还是只做协议转换。
- 如果 gateway 方法开始维护业务状态，通常意味着已经越过代理边界。

## 推荐结构

- Gateway 只做代理与适配
- 参数映射、鉴权信息、错误转换在接入层完成

## 典型实现方式

```text
HTTP Request
-> 参数转换
-> 调用上游服务
-> 响应适配
-> 对外返回
```

## 标准模板

```go
func (s *GatewayService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    resp, err := s.accountClient.GetAccount(ctx, convert(req))
    if err != nil {
        return nil, err
    }
    return convertReply(resp), nil
}
```

## 代码示例参考

```go
func (s *GatewayService) GetAccount(ctx http.Context) error {
    req := &v1.GetAccountRequest{Id: cast.ToUint32(ctx.Param("id"))}
    reply, err := s.accountClient.GetAccount(ctx, req)
    if err != nil {
        return err
    }
    return ctx.Result(200, reply)
}
```

## 项目通用 HTTP 注册示例

```go
func NewHTTPServer(
    accountService *openservice.AccountService,
    accountAdminService *adminservice.AccountService,
    c *conf.Server,
) *http.Server {
    srv := http.NewServer(...)

    openv1.RegisterAccountServiceHTTPServer(srv, accountService)
    adminv1.RegisterAccountServiceHTTPServer(srv, accountAdminService)

    return srv
}
```

## 内部 Gateway 协议示例

```proto
service AuthService {
  rpc CreateAuth (Auth) returns (CreateAuthReply) {};
  rpc ParseAuth (ParseAuthRequest) returns (Auth) {};
  rpc RefreshAuth (RefreshAuthRequest) returns (RefreshAuthReply) {};
}

message Auth {
  string subject = 1;
  string audience = 2;
  string viewer_id = 3;
  bool sso_mod = 4;
  map<string, string> extra = 5;
}
```

## 参数透传示例

```go
func (r *authRepo) ParseAuth(ctx context.Context, token string, subject string) (*gatewaybiz.Auth, error) {
    reply, err := r.authClient.ParseAuth(ctx, &gatewayv1.ParseAuthRequest{
        Token:   token,
        Subject: subject,
    })
    if err != nil {
        return nil, err
    }
    return authConvert(reply), nil
}
```

## 常见坑

- 请求和响应转换散落在多个文件
- 同一代理逻辑的鉴权、路由、转换不在一个闭环里
- Gateway 直接开始补业务字段

## 相关 Rule

- `../rules/gateway-rule.md`
- `../rules/entry-security-rule.md`

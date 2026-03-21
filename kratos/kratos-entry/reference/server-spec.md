# Server Reference

## 这个主题解决什么问题

说明 Kratos 项目中 gRPC/HTTP 服务实现和注册通常如何落位，以及新增接口时需要联动哪些位置。

## 适用场景

- 新增 gRPC/HTTP 接口
- 调整服务注册
- 检查 Service 与 Server 的连接关系

## 设计意图

Server 参考主要解释服务实现层为什么更像协议适配器，以及注册代码为什么要保持统一入口。

- 服务实现越薄，业务变化越容易收敛到 UseCase，协议变化也越容易单独处理。
- 注册路径统一后，新增接口时更容易知道该改哪里，而不是全仓库搜索散落路由。
- 对外接口和内部业务分层稳定后，鉴权、观测和排障都会更顺畅。

## 实施提示

- 先从 proto/service 契约反推 service 方法签名。
- 再把参数转换、身份提取和返回映射留在 service 层。
- 如果一个 service 方法出现大量状态判断和数据拼装，通常应回看 UseCase/Repo 边界。

## 推荐结构

- `internal/service/`：服务实现
- `internal/server/`：HTTP / gRPC server 初始化与注册

## 典型实现方式

```text
proto 定义
-> service 实现
-> server 注册
-> wire/provider 检查
```

## 标准模板

```go
func NewHTTPServer(cfg *conf.Server, svc *service.AccountService) *http.Server {
    srv := http.NewServer(...)
    v1.RegisterAccountHTTPServer(srv, svc)
    return srv
}
```

## 代码示例参考

```go
type AccountService struct {
    v1.UnimplementedAccountServer
    uc *biz.AccountUseCase
}

func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    info, err := s.uc.GetAccount(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return &v1.GetAccountReply{Info: convertToProto(info)}, nil
}
```

## 项目通用注册示例

```go
func NewGRPCServer(
    accountService *openservice.AccountService,
    accountAdminService *adminservice.AccountService,
    accountInnerService *innerservice.AccountService,
    c *conf.Server,
    _ log.Logger,
) *grpc.Server {
    srv := grpc.NewServer(...)

    openv1.RegisterAccountServiceServer(srv, accountService)
    adminv1.RegisterAccountServiceServer(srv, accountAdminService)
    innerv1.RegisterAccountServiceServer(srv, accountInnerService)

    return srv
}
```

## Service 层参数转换示例

```go
func (s *AccountService) GetAccountStatus(ctx context.Context, req *v1.GetAccountStatusRequest) (*v1.GetAccountStatusReply, error) {
    userCode := strconv.FormatUint(uint64(metadata.GetViewerID(ctx)), 10)
    status, err := s.uc.GetAccountStatus(ctx, userCode)
    if err != nil {
        return nil, err
    }
    return &v1.GetAccountStatusReply{OpenStatus: uint32(status)}, nil
}
```

## 常见坑

- 只写了 service 实现，忘记 server 注册
- HTTP 和 gRPC 注册不一致
- 新接口依赖已变，但 provider 未同步

## 相关 Rule

- `../rules/server-rule.md`
- `../rules/codegen-rule.md`

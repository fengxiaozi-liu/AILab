# Layer Spec

## business service 目录结构

| 目录 | 作用 | 不做什么 |
|------|------|------|
| `internal/server` | 注册 gRPC/HTTP 服务、中间件、启动装配 | 业务编排 |
| `internal/service` | 协议适配、参数转换、响应映射、gateway 转发 | 查询组织、状态流转 |
| `internal/biz` | UseCase、领域对象、事务边界、权限和状态流转 | 协议 DTO、DB 查询细节 |
| `internal/data` | Repo、查询组织、relation 装配、依赖访问 | 业务状态编排 |
| `internal/listener` | 本地事件、生命周期 hook、回调入口 | 业务流程和事务边界 |
| `internal/consumer` | MQ consumer、消息解析、最小入参组装 | 业务流程和事务边界 |
| `internal/crontab` 或等价目录 | scheduler、job 注册、调度基础设施 | 业务流程和事务边界 |

```text
// ✅ business service 典型目录
internal/
├── server/
├── service/
├── biz/
├── data/
├── listener/
├── consumer/
└── crontab/
```

---

## 完整项目层级

| 层 | 做法 | 不做什么 |
|------|------|------|
| `server` | 注册 gRPC/HTTP 服务、中间件、启动装配 | 业务编排 |
| `service` | 协议适配、参数转换、响应映射、gateway 转发 | 查询组织、状态流转 |
| `biz/usecase` | 业务编排、权限、事务、状态流转 | 协议 DTO、DB 查询细节 |
| `data/repo` | 查询组织、relation 装配、依赖访问 | 业务状态编排 |
| `listener/consumer/cron` | 组件触发、解析、最小入参组装、调用 UseCase | 业务流程和事务边界 |

```text
// ✅ 业务项目 的完整层级
Request
-> server
-> service
-> biz/usecase
-> data/repo
-> DB / depend

Event / MQ / Cron
-> listener / consumer / cron
-> biz/usecase
-> data/repo
```

```text
// ❌ 越层
Request -> service -> repo
Request -> service -> usecase -> dto build -> repo
listener -> repo -> state transition
```

---

## 层级判断

| 条件 | 做法 |
|------|------|
| 注册 server、挂中间件、启动装配 | 归 `server` |
| 协议适配、HTTP/gRPC/gateway 参数转换 | 归 `service` |
| 业务编排、权限、事务、状态流转 | 归 `biz/usecase` |
| 查询、filter、relation 装配、远程依赖补全 | 归 `data/repo` |
| 事件监听、消息消费、定时触发 | 归 `listener/consumer/cron`，再调 UseCase |

```go
// ✅ service
func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    account, err := s.uc.GetAccount(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return &v1.GetAccountReply{Info: convertAccount(account)}, nil
}
```

```go
// ✅ usecase
func (u *AccountUseCase) Submit(ctx context.Context, req *SubmitAccountRequest) error {
    return u.tm.InTx(ctx, func(ctx context.Context) error { ... })
}
```

```go
// ✅ repo
func (r *accountRepo) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*biz.Account, error) {
    ...
}
```

---

## relation 收口

| 条件 | 做法 |
|------|------|
| 需要 relation | UseCase 声明需求，Repo 实现装配 |
| 上层发现字段为空 | 不在 Service/UseCase 临时补查 |

```go
// ✅
account, err := u.accountRepo.GetAccount(ctx, id, filter.WithRelation(openenum.AccountCollectRelation))
```

```go
// ❌
if account.CollectInfo == nil {
    account.CollectInfo, err = u.collectRepo.Get(...)
}
```

---

## 组合场景

```text
{business}Service
|- internal/server   -> gRPC/HTTP register
|- internal/service  -> protocol adapter
|- internal/biz      -> usecase orchestration
|- internal/data     -> repo + relation assembly
|- internal/listener -> local event / lifecycle hook
|- internal/consumer -> mq consumer
|- internal/crontab  -> scheduler / cron infra
```

这个组合场景同时满足：

- `{business}Service` 是完整项目，不是 gateway/base 裁剪形态
- 请求链路和事件链路都有明确层级
- listener/consumer 不替代 usecase

---

## 常见错误模式

```go
// ❌ service 直接调 repo
s.accountRepo.Get(ctx, req.Id)
```

```go
// ❌ usecase 拼协议 reply
return &v1.GetAccountReply{}, nil
```

```go
// ❌ listener 直接做状态流转
if account.Status == ...
```

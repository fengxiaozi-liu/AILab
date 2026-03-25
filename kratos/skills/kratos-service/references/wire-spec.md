# Wire Spec

## ProviderSet 组织

| 条件 | 做法 |
|------|------|
| 新增构造函数 | 补到所在层的 `ProviderSet` |
| 新增模块注入 | 沿现有装配入口增量扩展 |

```go
// ✅ 每层维护自己的 ProviderSet
var ProviderSet = wire.NewSet(
    NewAccountRepo,
    NewAccountUseCase,
    NewAccountService,
)
```

```go
// ❌ 另起散落注入入口
var AccountModuleProviderSet = wire.NewSet(NewAccountUseCase)
```

---

## 构造函数依赖表达

| 条件 | 做法 |
|------|------|
| 构造对象依赖已知 | 在构造函数参数中显式声明 |
| 依赖尚未稳定 | 先回到 proto/usecase/service 设计，不抢改 wire |

```go
// ✅ 显式依赖
func NewAccountUseCase(repo biz.AccountRepo, logger log.Logger) *AccountUseCase {
    return &AccountUseCase{
        repo: repo,
        log:  log.NewHelper(logger),
    }
}
```

```go
// ❌ 隐式拿全局或配置拼装
func NewAccountUseCase() *AccountUseCase {
    return &AccountUseCase{
        repo: globalRepo,
    }
}
```

```go
// ⚠️ 配置对象优先注入父配置，再在构造函数读取子字段
func NewHTTPServer(c *conf.Server, svc *service.AccountService) *http.Server {
    httpc := c.Http
    _ = httpc
    ...
}
```

---

## 装配顺序

| 条件 | 做法 |
|------|------|
| 新增能力尚在设计期 | 先稳定 domain、proto、service，再改 wire |
| 只是补 provider 接线 | 直接更新 ProviderSet 和 wire 入口 |

```text
// ✅ 顺序
domain / usecase
-> proto / service
-> repo / provider
-> ProviderSet
-> wire
```

```text
// ❌ 反过来用 wire 反推设计
wire
-> provider
-> service
-> proto
```

---

## Wire 入口

| 条件 | 做法 |
|------|------|
| 应用启动装配 | 在统一 `initApp` 或等价入口维护 |
| 新增层级模块 | 挂到现有 `wire.Build(...)` 中 |

```go
// ✅ 统一入口
func initApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
    panic(wire.Build(
        server.ProviderSet,
        service.ProviderSet,
        biz.ProviderSet,
        data.ProviderSet,
        newApp,
    ))
}
```

```go
// ❌ 在业务文件中单独写 wire.Build
func initAccountApp() {
    panic(wire.Build(service.ProviderSet))
}
```

---

## 组合场景

```go
// internal/service/service.go
var ProviderSet = wire.NewSet(NewAccountService)

func NewAccountService(uc *biz.AccountUseCase) *AccountService {
    return &AccountService{uc: uc}
}
```

```go
// cmd/server/wire.go
func initApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
    panic(wire.Build(
        server.ProviderSet,
        service.ProviderSet,
        biz.ProviderSet,
        data.ProviderSet,
        newApp,
    ))
}
```

这个组合场景同时满足：

- ProviderSet 归属清晰
- 构造函数显式声明依赖
- wire 入口统一
- 改动顺序可控

---

## 常见错误模式

```go
// ❌ provider set 散落
var XXXSet = wire.NewSet(...)
```

```go
// ❌ 构造函数偷拿全局依赖
repo: globalRepo
```

```go
// ❌ 用 wire 反推上游设计
panic(wire.Build(...))
```

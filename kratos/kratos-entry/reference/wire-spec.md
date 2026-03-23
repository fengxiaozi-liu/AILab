# Wire Spec

## Provider 组织决策

| 层 | ProviderSet 位置 | 包含内容 |
|----|----------------|---------|
| data | `internal/data/data.go` | `NewData`, `NewXxxRepo` |
| biz | `internal/biz/biz.go` | `NewXxxUseCase` |
| service | `internal/service/service.go` | `NewXxxService` |

---

## ProviderSet 模板

```go
// ✅ 每层维护自己的 ProviderSet（以 data 层为例）
var ProviderSet = wire.NewSet(
    NewData,
    NewAccountRepo,
    NewStoreRepo,
    NewAdminUserRepo,
)

// ✅ 构造函数显式声明所有依赖
func NewAccountUseCase(
    repo biz.AccountRepo,
    tx biz.Tx,
    eventBus eventbus.EventBus,
    logger log.Logger,
) *biz.AccountUseCase {
    return &biz.AccountUseCase{
        repo:     repo,
        tx:       tx,
        eventBus: eventBus,
        log:      log.NewHelper(logger),
    }
}

// ❌ 构造函数使用 struct 初始化绕过 wire 依赖声明
func NewAccountUseCase() *biz.AccountUseCase {
    return &biz.AccountUseCase{}  // ❌ 依赖未声明，wire 无法注入
}
```

---

## 新增 provider 步骤

```text
// ✅ 增量扩展步骤（不新建注入入口）
1. 在对应层的构造函数文件中添加 NewXxxRepo / NewXxxUseCase / NewXxxService
2. 在该层的 ProviderSet 中追加构造函数
3. 在 cmd/server/wire.go 中引用顶层 ProviderSet
4. 运行 cd cmd/server && wire 重新生成

// ❌ 为新聚合根单独新建 ProviderSet 文件
// cmd/server/account_wire.go  ❌
```

---

## 配置注入

```go
// ✅ 配置对象按父配置整体注入，构造函数内读取子字段
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
    db, err := openDB(c.Database)  // ✅ 构造函数内读取子字段
    if err != nil { return nil, nil, err }
    ...
}

// ❌ 每个构造函数单独接收细碎配置字段
func NewData(host string, port int, dbName string) (*Data, func(), error) {  // ❌ 松散配置
    ...
}
```

---

## 联动顺序

```text
// 正确顺序：领域模型先，Wire 最后
1. 聚合根建模（biz/account.go）
2. proto/service 契约（api/.../account.proto）
3. Repo/UseCase/Service 实现
4. 更新 ProviderSet
5. wire 生成（cd cmd/server && wire）
```

---

## 组合场景

```go
// 完整：三层 ProviderSet + wire 入口
// internal/data/data.go
var ProviderSet = wire.NewSet(NewData, NewAccountRepo)

// internal/biz/biz.go
var ProviderSet = wire.NewSet(NewAccountUseCase)

// internal/service/service.go
var ProviderSet = wire.NewSet(NewAccountService)

// cmd/server/wire.go
func initApp(c *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
    panic(wire.Build(
        data.ProviderSet,
        biz.ProviderSet,
        service.ProviderSet,
        server.NewHTTPServer,
        server.NewGRPCServer,
        newApp,
    ))
}
```

---

## 常见错误模式

```go
// ❌ 构造函数未声明依赖（wire 无法自动注入）
func NewAccountUseCase() *biz.AccountUseCase { return &biz.AccountUseCase{} }

// ❌ 在 proto/service 未稳定时改动 wire
// 导致 wire 生成后被 proto 变更覆盖，需重复操作

// ❌ 同一个功能在多个 ProviderSet 重复声明
var DataProviderSet = wire.NewSet(NewAccountRepo)
var BizProviderSet  = wire.NewSet(NewAccountRepo, NewAccountUseCase)  // ❌ NewAccountRepo 重复
```

# Util Spec

## 纳入 util 的判断标准

| 特征 | 是否可纳入 util |
|------|----------------|
| 无状态，纯函数，可独立测试 | ✅ 可以 |
| 已在 2+ 模块复用 | ✅ 可以 |
| 仅服务于单一业务流程 | ❌ 不可以，放在该业务模块 |
| 访问 Repo / Config / Business Context | ❌ 不可以 |

---

## 并发工具：Parallel

```go
// ✅ util.Parallel 并发执行多个任务，统一收集错误
results, err := util.Parallel(ctx,
    func(ctx context.Context) (interface{}, error) { return getAccount(ctx) },
    func(ctx context.Context) (interface{}, error) { return getStore(ctx) },
)
if err != nil { return err }

// ❌ 手写 goroutine + channel（每处手写，不一致且易泄漏）
errCh := make(chan error, 2)
go func() { _, err := getAccount(ctx); errCh <- err }()
go func() { _, err := getStore(ctx); errCh <- err }()
// ❌ 未统一 cancel 传播
```

---

## 工具函数命名要求

```go
// ✅ 名字直接表达能力，无歧义
util.Parallel(...)       // 并发执行
util.Chain(...)          // 链式执行
util.DecodeBase64(...)   // base64 解码
util.Decimal(...)        // 精度计算

// ❌ 模糊入口（不知道做什么）
util.Process(...)        // ❌ 语义不清
util.Handle(...)         // ❌
util.Do(...)             // ❌
```

---

## 文件拆分

```text
// ✅ 按能力主题拆分，各自配测试
internal/pkg/util/
├── parallel.go       // 并发工具
├── parallel_test.go
├── base64.go         // 编码工具
├── base64_test.go
├── decimal.go        // 数值工具
└── decimal_test.go

// ❌ 所有工具堆在 util.go
internal/pkg/util/util.go  // ❌
```

---

## 禁止在 util 中访问业务依赖

```go
// ❌ util 函数访问 config
func GetMaxRetry() int {
    return config.GetGlobalConfig().RetryTimes  // ❌ util 不依赖业务 config
}

// ❌ util 函数访问 repo
func FillUserInfo(ctx context.Context, ids []uint32) ([]*User, error) {
    return userRepo.GetUsers(ctx, ids)  // ❌ 不是工具函数，是业务函数
}

// ❌ 为单一流程创建的 util，未被复用
// util/open_account_helper.go  ← 只被 account usecase 使用 ❌
```

---

## 组合场景

```go
// 完整：UseCase 使用 util.Parallel 并发补查关联数据
func (u *AccountUseCase) GetAccountWithRelations(ctx context.Context, id uint32) (*Account, error) {
    account, err := u.accountRepo.GetAccount(ctx, id)
    if err != nil { return nil, err }

    _, err = util.Parallel(ctx,
        func(ctx context.Context) (interface{}, error) {
            return nil, u.adminUserDepend.FillCheckUser(ctx, account)
        },
        func(ctx context.Context) (interface{}, error) {
            return nil, u.storageDepend.FillStorage(ctx, account)
        },
    )
    if err != nil { return nil, err }
    return account, nil
}
```

---

## 常见错误模式

```go
// ❌ 一次性业务函数放入 util
// util/account_convert.go  ← 只被 AccountService 调用 ❌

// ❌ 无测试的工具函数
// util/decimal.go 无对应 decimal_test.go  ❌

// ❌ 函数名模糊
util.Run(ctx, tasks...)   // ❌ 和 util.Parallel 重复，语义不清
util.Execute(...)          // ❌
```

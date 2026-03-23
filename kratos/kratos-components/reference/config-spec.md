# Config Spec

## 配置项规范

| 条件 | 做法 |
|------|------|
| 新增配置项 | MUST 提供默认值或 `Normalize()` 校验 |
| 区分环境差异 | 用不同 yaml 文件或环境变量覆盖，不硬编码 |
| 配置影响业务分支 | 补启动校验或最小验证样例 |

```go
// ✅ 结构体字段有 json/yaml tag，提供 Normalize 默认值
type AccountConfig struct {
    SyncCron string `json:"sync_cron" yaml:"sync_cron"`
    PageSize  int32  `json:"page_size" yaml:"page_size"`
}

func (c *AccountConfig) Normalize() {
    if c.PageSize <= 0 {
        c.PageSize = 100
    }
    if c.SyncCron == "" {
        c.SyncCron = "0 */5 * * * *"
    }
}

// ❌ 无默认值，调用方拿到零值崩溃
type AccountConfig struct {
    PageSize int32
}
// 未执行任何 Normalize，PageSize=0 → 分页查询全量数据
```

---

## Depend 超时配置

```go
// ✅ 超时来自配置，不硬编码
type Depend struct {
    Timeout string `json:"timeout" yaml:"timeout"`  // e.g. "3s"
}

func NewAdminUserRepo(conf *conf.Depend, ...) *adminUserRepo {
    timeout, _ := time.ParseDuration(conf.AdminUser.Timeout)
    // 用 timeout 构造带截止时间的 context
}

// ❌ 超时硬编码在代码里
ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
```

---

## 配置落位

```go
// ✅ 构造函数中读取所需子字段，不把整个 conf 传递下去
func NewAccountUseCase(repo biz.AccountRepo, conf *conf.Business) *AccountUseCase {
    return &AccountUseCase{
        repo:     repo,
        pageSize: conf.Account.PageSize,
    }
}

// ❌ 把整个 conf 注入业务层
func NewAccountUseCase(repo biz.AccountRepo, conf *conf.Bootstrap) *AccountUseCase {
    return &AccountUseCase{repo: repo, conf: conf}
}
```

---

## 组合场景

```go
// 完整：配置结构 + Normalize + 构造函数使用
type AccountConfig struct {
    SyncCron string `json:"sync_cron" yaml:"sync_cron"`
    PageSize  int32  `json:"page_size" yaml:"page_size"`
    Timeout   string `json:"timeout" yaml:"timeout"`
}

func (c *AccountConfig) Normalize() {
    if c.PageSize <= 0 {
        c.PageSize = 100
    }
    if c.SyncCron == "" {
        c.SyncCron = "0 */5 * * * *"
    }
    if c.Timeout == "" {
        c.Timeout = "3s"
    }
}

func NewAccountRepo(data *Data, conf *conf.Data, logger log.Logger) *accountRepo {
    conf.Account.Normalize()
    timeout, _ := time.ParseDuration(conf.Account.Timeout)
    return &accountRepo{
        data:    data,
        timeout: timeout,
        log:     log.NewHelper(logger),
    }
}
```

---

## 常见错误模式

```go
// ❌ 新增配置项无默认值说明
type Config struct {
    NewFlag string  // 没有 Normalize，调用方拿到空字符串
}

// ❌ 配置字段名只服务于实现细节
type Config struct {
    V1EndpointForLegacyQuery string  // 难以理解、难以回滚

// ❌ 多个无关模块共用同一配置字段
type Config struct {
    PageSize int32  // 被 account、store、order 等模块随意共用
}
```

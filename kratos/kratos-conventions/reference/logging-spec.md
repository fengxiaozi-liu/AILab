# Logging Spec

## Logger 初始化与使用方式

| 场景 | 做法 |
|------|------|
| 长生命周期结构体 | 注入 `log.Logger`，内部创建 `log.NewHelper` |
| 需要模块标识 | `log.With(logger, "module", "...")` |
| 请求上下文日志 | `helper.WithContext(ctx)` |

```go
// ✅ 注入 log.Logger，内部转换 Helper
type AccountListener struct {
    log *log.Helper
}

func NewAccountListener(logger log.Logger) *AccountListener {
    return &AccountListener{
        log: log.NewHelper(log.With(logger, "module", "listener/account")),
    }
}

// ❌ 直接构造底层 zap logger，绕过 Kratos 日志体系
func NewAccountListener() *AccountListener {
    zapLogger, _ := zap.NewProduction()
    return &AccountListener{log: zapLogger}
}
```

---

## 日志输出规范

```go
// ✅ 使用 WithContext 透传请求上下文
r.log.WithContext(ctx).Info("schema preload: start")
r.log.WithContext(ctx).Errorf("schema preload: download failed: %v", err)

// ✅ 错误日志包含动作语义 + 业务标识 + err
l.log.WithContext(ctx).Errorf("unmarshal account event failed: account_id=%d err=%v", accountID, err)

// ❌ 无上下文，无业务标识
log.Printf("error: %v", err)

// ❌ 只记日志，吞掉 err
l.log.Error("update failed")
return nil  // 错误被忽略
```

---

## 敏感信息

```go
// ✅ 脱敏后记录
l.log.Infof("user login: user_id=%d ip=%s", userID, maskIP(ip))

// ❌ 打印敏感字段
l.log.Infof("user login: token=%s password=%s cert=%s", token, password, certNo)
```

---

## 高频路径日志控制

```go
// ✅ SQL 调试日志仅在 debug 模式启用
if conf.Debug {
    drv = dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
        log.Context(ctx).Info(i...)
    })
}

// ❌ 在批量处理热点路径中无控制输出
for _, item := range list {
    log.Infof("processing item: %v", item)  // 千万条记录时引发日志风暴
}
```

---

## 生命周期日志

```go
// ✅ 资源关闭日志简短稳定
helper := log.NewHelper(logger)
helper.Info("closing the data resources")
defer func() { helper.Info("data resources closed") }()

// ✅ 预加载日志带关键节点
r.log.WithContext(ctx).Info("schema preload: start")
r.log.WithContext(ctx).Info("schema preload: done")
```

---

## 组合场景

```go
// 完整：初始化 + 模块标识 + 上下文日志 + 错误日志
type SchemaRepo struct {
    log *log.Helper
}

func NewSchemaRepo(logger log.Logger) *SchemaRepo {
    return &SchemaRepo{
        log: log.NewHelper(log.With(logger, "module", "repo/schema")),
    }
}

func (r *SchemaRepo) PreloadAll(ctx context.Context) error {
    r.log.WithContext(ctx).Info("schema preload: start")
    if err := r.downloadManifest(ctx); err != nil {
        r.log.WithContext(ctx).Errorf("schema preload: download manifest failed: %v", err)
        return err
    }
    r.log.WithContext(ctx).Info("schema preload: done")
    return nil
}
```

---

## 常见错误模式

```go
// ❌ 同一模块混用 zap 和 log.Helper，字段不一致
// ❌ 打印 token / password / PII
log.Infof("request: token=%s body=%v", token, req)

// ❌ 只记日志不返回错误（吞错）
if err != nil {
    log.Error(err)
    return nil
}

// ❌ 不使用 WithContext，日志无请求追踪信息
r.log.Info("processing account")
```

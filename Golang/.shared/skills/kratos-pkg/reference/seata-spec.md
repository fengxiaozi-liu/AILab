# Seata 参考

## 这个主题解决什么问题

统一全局事务模板、传播级别和执行回调封装，避免各模块重复实现事务样板代码。

## 适用场景

- 新增或调整 Seata 全局事务入口
- 调整传播级别、事务名、超时、锁重试等事务配置

## 推荐结构或实现方式

- 通过显式 `GtxConfig` 和 callback 执行业务逻辑。
- 事务包装只负责 begin、commit、rollback 和传播控制。

## 标准模板

```go
type GtxConfig struct {
    Timeout           time.Duration
    Name              string
    Propagation       tm.Propagation
    LockRetryInternal time.Duration
    LockRetryTimes    int16
}

func WithGlobalTx(ctx context.Context, gc *GtxConfig, business CallbackWithCtx) error {
    // begin / commit / rollback
    return nil
}
```

## Good Example

- 通过 `WithGlobalTx` 统一承接 callback，把事务名、超时、传播级别放到配置中，而不是散落在业务函数里。

## 常见坑

- 在事务 helper 中加入业务补偿和业务状态判断
- 每个模块都复制一份全局事务模板

## 相关 rule / 相关 reference

- `../rules/seata-rule.md`

# Context Pattern

## 这个主题解决什么问题

说明 Go 中 `context.Context` 如何在调用链上传递、取消和设置超时。

## 适用场景

- 远程调用
- IO 操作
- goroutine 协调

## 设计意图

Context 模式主要解释请求级状态如何沿调用链传递，而不是把 context 当作通用参数容器。

- 理解 context 的生命周期后，更容易正确处理超时、取消和链路透传。
- context 使用方式清晰时，服务间调用和并发场景也更容易保持一致。
- 把它看成请求边界对象，比把它当成随手塞值的 map 更能保持可读性。

## 实施提示

- 先判断一个值是不是只在当前请求生命周期内有效。
- 再决定它更适合作为显式参数还是 context 元数据。
- 如果 goroutine 会跨越原始请求边界，先确认新的 context 生命周期如何定义。

## 推荐模式

- `ctx` 作为函数第一个参数
- 调用链向下透传
- 需要超时时用 `context.WithTimeout`

## 标准模板

```go
func Do(ctx context.Context, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    return work(ctx)
}
```

## 常见坑

- 忽略上游 `ctx`
- 忘记 `cancel()`
- 用 `Background()` 打断调用链

## 相关 Rule

- `../rules/context-rule.md`

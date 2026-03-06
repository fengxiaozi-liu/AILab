# Concurrency Pattern

## 这个主题解决什么问题

说明 Go 并发中 goroutine、channel、共享状态的常见组织方式。

## 适用场景

- 后台 worker
- 并发 fan-out / fan-in
- channel 协调

## 设计意图

并发模式参考不是为了鼓励多开 goroutine，而是帮助判断什么时候并发真的能带来吞吐或延迟收益。

- 并发结构选对后，代码会更容易看出任务拆分、同步点和错误汇总位置。
- 如果只看到 goroutine 而看不到边界，后续排查泄漏、阻塞和竞态会很困难。
- 先理解工作拆分模型，再选择 channel、mutex 或 errgroup，会比直接套模板更稳。

## 实施提示

- 先明确是并发扇出、流水线还是后台消费者模型。
- 再确定谁负责收集错误、谁负责关闭 channel、谁负责停止信号。
- 如果并发只是为了少写几次顺序调用，先评估复杂度是否真的值得。

## 推荐模式

- 启动 goroutine 前先定义退出机制
- 尽量减少共享状态
- 共享状态需要同步保护

## 标准模板

```go
func Worker(ctx context.Context, in <-chan int, out chan<- int) {
    for {
        select {
        case <-ctx.Done():
            return
        case v, ok := <-in:
            if !ok {
                return
            }
            out <- v * 2
        }
    }
}
```

## 常见坑

- goroutine 没有退出路径
- channel 所有权不清
- 锁粒度过大或共享状态过多

## 相关 Rule

- `../rules/concurrency-rule.md`

# Concurrency（语言层）

## Scope

只覆盖 Go 并发原语与生命周期：goroutine、channel、mutex、atomic、select，以及 “如何退出”。

### 1) goroutine 启动点要可解释（谁负责生命周期）

```go
go r.reader(client)
go r.writer(client)
```

### 2) 共享状态：`sync.RWMutex` + `defer Unlock`

```go
mu.Lock()
defer mu.Unlock()
```

### 3) channel：使用缓冲、select 避免阻塞（通用原则）

```go
ch := make(chan T, 1024)
select {
case ch <- v:
default:
    // 丢弃/降级/计数（按业务决定）
}
```

## Rule hints

- 启动 goroutine 前先定义退出路径（`ctx.Done()` / `close(done)` / 关闭输入 channel）
- channel 所有权：通常“创建/写入方负责 close”，接收方不要 close（除非明确所有权）
- 共享内存必须加锁或用原子；优先减少共享

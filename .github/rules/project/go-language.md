---
name: Go-Language-Rule
description: Go语言 层面的规则约束
---

# Go Language Rules（语言层规则）

> 本文件只描述 Go 语言层面的规则：用来约束 AI 写出的 Go 代码“长得像 Go、可维护、可控”。
> 每个主题包含：Principles / Specification / Prohibit / Demo（Good & Bad）。

## Errors

### Principles

- 错误要可追踪（保留根因）、可判定（可用 `errors.Is/As` 判断）、可读（上下文清晰）
- 在“项目已有统一错误体系”的前提下，语言层不应随意发明对外错误语义

### Specification

- 发生错误时使用早返回：`if err != nil { return ..., err }`
- 需要补充上下文时使用 `%w` 包装，保留错误链
- 判断特定错误使用 `errors.Is/As`（而不是字符串比较）
- 当项目存在框架级错误体系（例如 `internal/error/...` 提供的错误构造函数）时：在**对外/跨层边界**返回错误，必须返回框架级错误（而不是原生错误）

### Prohibit

- 禁止吞错（只 log 不 return / return nil err）
- 禁止用字符串内容判断错误类型（如 `err.Error() == "xxx"`）
- 当项目存在框架级错误体系时：禁止在对外/跨层边界直接返回原生错误（`errors.New(...)` / `fmt.Errorf(...)` / `return err` 作为最终错误）

### Demo

Good:
```go
func ReadUser(id int) (string, error) {
	v, err := loadFromDB(id)
	if err != nil {
		// Good: 底层 err 可用于日志/排障，但对外返回使用项目统一的框架级错误
		// 例如：return "", baseerror.ErrorFailed(ctx)
		return "", baseerror.ErrorFailed(ctx)
	}
	return v, nil
}
```

Bad:
```go
func ReadUser(id int) (string, error) {
	v, err := loadFromDB(id)
	if err != nil {
		// BAD: 在已有框架/项目错误体系时仍返回原生错误作为最终错误
		return "", fmt.Errorf("db failed: %w", err)
	}
	return v, nil
}
```

Good（项目错误示例）:
```go
func ValidateName(ctx context.Context, name string) error {
	if name == "" {
		// Good: 返回项目统一错误（见 internal/error/...）
		return baseerror.ErrorBadRequest(ctx)
	}
	return nil
}
```

Bad（项目反例）:
```go
func ValidateName(ctx context.Context, name string) error {
	if name == "" {
		// BAD: 明明有项目统一错误体系，仍返回原生错误
		return errors.New("name is required")
	}
	return nil
}
```

## Context

### Principles

- `context` 用于请求级取消/超时/元信息传递；应沿调用链向下传递
- 派生的 `context` 必须及时释放（`cancel()`）

### Specification

- `ctx context.Context` 作为函数的第一个参数（跨层/跨包尤其如此）
- 使用 `WithTimeout/WithCancel` 后必须 `defer cancel()`

### Prohibit

- 禁止在请求路径中随意使用 `context.Background()` 替代入参 `ctx`
- 禁止把 `context.Context` 存到 struct 字段长期持有（除非你明确生命周期并能证明安全）

### Demo

Good:
```go
func Do(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return work(ctx)
}
```

Bad:
```go
func Do(ctx context.Context) error {
	// BAD: 忽略调用方的取消/超时
	return work(context.Background())
}
```

## Concurrency

### Principles

- goroutine 必须可退出（防泄漏）；共享数据必须同步（防竞态）
- channel 必须有清晰的所有权（谁创建、谁关闭）

### Specification

- 启动 goroutine 前定义退出机制（`ctx.Done()` / 关闭输入 channel / done 信号）
- 共享状态用 `sync.Mutex/RWMutex` 或 `atomic`；加锁后 `defer Unlock()`
- 发送到 channel 的阻塞风险要可控（必要时用缓冲 + `select`）

### Prohibit

- 禁止启动无退出条件的 goroutine（如 `for {}` 无 `ctx.Done()`/无输入关闭）
- 禁止在未明确所有权时由接收方 `close(channel)`

### Demo

Good:
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

Bad:
```go
func Worker(in <-chan int, out chan<- int) {
	// BAD: goroutine 泄漏；无退出条件
	go func() {
		for {
			v := <-in
			out <- v * 2
		}
	}()
}
```

## Defer & Cleanup

### Principles

- 资源与锁必须成对释放；用 `defer` 保证异常路径也能释放

### Specification

- 资源获取/加锁成功后第一时间 `defer Close/Unlock`
- 临时修改状态后用 `defer` 恢复（避免遗漏恢复路径）

### Prohibit

- 禁止手写多分支释放导致遗漏（如某条 return 路径忘记 Unlock/Close）
- 禁止在高频循环里滥用 `defer`（除非你明确需要延迟释放）

### Demo

Good:
```go
mu.Lock()
defer mu.Unlock()

f, err := os.Open(name)
if err != nil {
	return err
}
defer f.Close()
```

Bad:
```go
mu.Lock()
if err := do(); err != nil {
	// BAD: 忘记 Unlock，造成死锁
	return err
}
mu.Unlock()
```

## Generics

### Principles

- 泛型用于“可复用的通用结构/算法”，避免为单点业务增加复杂度
- 约束（constraints）应尽量收敛，避免 `any` 污染调用链

### Specification

- 泛型优先用于通用工具层；业务域类型谨慎使用
- 能用具体类型清晰表达时，不要强行泛型化

### Prohibit

- 禁止为了“看起来通用”而引入复杂约束/多类型参数，导致可读性下降

### Demo

Good:
```go
type Handle[T, V any] func(ctx context.Context, req T, resp V) (bool, error)
```

Bad:
```go
// BAD: 为单点业务硬套泛型，类型参数无实际复用价值
type UserThing[A, B, C any] struct {
	A A
	B B
	C C
}
```

## Performance

### Principles

- 先写清晰正确，再做低风险优化；避免明显的可预见浪费（反复分配、无意义拷贝）

### Specification

- 能预估容量时对 slice/map 预分配
- 热路径避免无必要的临时对象创建（必要时使用复用手段）

### Prohibit

- 禁止在热循环中无脑 `append` 导致大量扩容（能预估就预估）

### Demo

Good:
```go
out := make([]int, 0, len(in))
for _, v := range in {
	out = append(out, v)
}
```

Bad:
```go
var out []int
for _, v := range in {
	// BAD: 频繁扩容（可预估却不预分配）
	out = append(out, v)
}
```

## Style

### Principles

- 可读性优先：命名清晰、控制流简单、文件组织符合 Go 直觉

### Specification

- 早返回减少缩进层级
- import 仅在冲突/提升可读性时使用别名
- 方法接收者风格保持一致（同一类型尽量统一用指针或值）

### Prohibit

- 禁止过度嵌套（多层 `if/for`）导致可读性差（可用早返回/拆函数）

### Demo

Good:
```go
func Parse(s string) (int, error) {
	if s == "" {
		return 0, errors.New("empty")
	}
	return strconv.Atoi(s)
}
```

Bad:
```go
func Parse(s string) (int, error) {
	// BAD: 不必要嵌套
	if s != "" {
		n, err := strconv.Atoi(s)
		if err == nil {
			return n, nil
		} else {
			return 0, err
		}
	} else {
		return 0, errors.New("empty")
	}
}
```

# Design Patterns（语言层映射）

## Scope

只覆盖“设计模式在 Go 里的语言级落地写法”，避免引入框架概念；以轻量组合与函数类型为主。

## Out of scope（必须交接）

- 框架级中间件/插件体系的具体规范（交给 Kratos/项目技能）
- 业务域建模方法论

## Patterns（demo）

### 1) Builder → Functional Options

```go
type Server struct {
	addr string
}

type Option func(*Server)

func WithAddr(addr string) Option {
	return func(s *Server) { s.addr = addr }
}

func NewServer(opts ...Option) *Server {
	s := &Server{addr: ":8080"}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
```

### 2) Decorator → Function Wrapper

```go
type Handler func(ctx context.Context, req any) (any, error)

func WithLogging(next Handler) Handler {
	return func(ctx context.Context, req any) (any, error) {
		start := time.Now()
		res, err := next(ctx, req)
		_ = start
		return res, err
	}
}
```

### 3) Strategy → Function Type

```go
type Comparator func(a, b int) bool

func Sort(xs []int, less Comparator) {
	sort.Slice(xs, func(i, j int) bool { return less(xs[i], xs[j]) })
}
```

### 4) Singleton → `sync.Once`

```go
var (
	once sync.Once
	inst *Thing
)

func GetThing() *Thing {
	once.Do(func() { inst = &Thing{} })
	return inst
}
```

### 5) Observer → Channel

```go
type Event struct{ Name string }

func Publisher(ctx context.Context, out chan<- Event) {
	defer close(out)
	out <- Event{Name: "created"}
}

func Subscriber(ctx context.Context, in <-chan Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-in:
			if !ok {
				return
			}
		}
	}
}
```


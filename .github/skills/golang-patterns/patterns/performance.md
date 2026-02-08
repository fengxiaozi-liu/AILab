# Performance（语言层）

## Scope

只覆盖 Go 语言层面的“性能意识”与可复用写法：减少分配、预分配、避免不必要的拷贝、`sync.Pool` 的适用边界。

## Out of scope（必须交接）

- DB/缓存/网络调用优化、链路观测、框架级连接池配置等
- 任何依赖具体框架的性能调优手段

## Patterns（demo）

### 1) slice 预分配容量（减少扩容与复制）

```go
items := make([]int, 0, len(input))
for _, v := range input {
	items = append(items, v)
}
```

### 2) map 预估容量（减少 rehash）

```go
m := make(map[string]int, len(keys))
for _, k := range keys {
	m[k]++
}
```

### 3) 避免在热路径频繁创建临时对象：`sync.Pool`（谨慎）

```go
var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

func Use() []byte {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	_, _ = buf.WriteString("hello")
	out := append([]byte(nil), buf.Bytes()...) // 复制，避免引用池内对象
	return out
}
```

使用边界（语言层）：

- `sync.Pool` 适合“高频创建/销毁、生命周期短、可复用”的临时对象（如 buffer）
- 放进 `Pool` 的对象不应被外部长期持有（否则复用会导致数据错乱）


# defer & Resource Cleanup（语言层）

## Scope

只覆盖 Go 语言层面的 `defer` 使用：资源释放、锁释放、临时状态恢复。

### 1) 临时修改状态后，立即 `defer` 恢复

```go
old := cfg.Name
cfg.Name = ""
defer func() { cfg.Name = old }()
```

### 2) 锁释放用 `defer Unlock()`

```go
mu.Lock()
defer mu.Unlock()
```

## Rule hints

- `defer` 放在“成功获得资源/完成状态变更”之后的第一时间
- 避免在循环里滥用 `defer` 导致延迟释放（除非你明确需要）

# Go Errors（语言层）

## Scope

只覆盖 Go 语言层面的错误表达与传播：`error`、包装（wrap）、判定（`errors.Is/As`）、就地返回（early return）。

## Out of scope（必须交接）

当你需要决定“对外错误长什么样”（错误码/Reason/国际化/统一格式化/HTTP 映射），必须交接到 Kratos/项目技能；本文件不提供替代规范。

## Patterns（demo）

### 1) 早返回（不要吞错）

```go
x, err := do()
if err != nil {
    return nil, err
}
```

### 2) 包装错误链（保留根因）

```go
if err := do(); err != nil {
    return fmt.Errorf("xxx失败: %w", err)
}
```

### 3) 用 `errors.Is` 判定“特定错误”

```go
if err != nil && !errors.Is(err, ErrNotFound) {
    return err
}
```

## Rule hints（语言层 hard edges）

- 不要用字符串比较判断错误类型；用 `errors.Is/As`
- 包装错误时用 `%w` 保留链路（调用方才可判定根因）

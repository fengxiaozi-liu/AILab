# Style & Organization（语言层）

## Scope

只覆盖 Go 语言层面的代码风格：注释、控制流、可读性等；目标是让代码“长得像 Go”且可读。

> 命名规范（文件/变量/接口/结构体/函数/错误/枚举等）已单独沉淀到命名规范章节，不在本文件展开。

## Out of scope（必须交接）

- “分层目录怎么定、哪些包放哪一层”这类项目结构约定（交给项目/框架技能）
- 对外 API/协议、配置体系等

## Patterns（demo）

### 1) 控制流：早返回减少嵌套

```go
func Do(x int) error {
	if x <= 0 {
		return errors.New("invalid")
	}
	// 正常路径
	return nil
}
```

### 2) 注释：导出符号按 GoDoc 习惯书写

```go
// UserRepo provides access to users.
type UserRepo interface {
	Get(ctx context.Context, id UserID) (*User, error)
}
```

### 3) 控制流：减少“无意义 else”

```go
func Parse(s string) (int, error) {
	if s == "" {
		return 0, errors.New("empty")
	}
	return strconv.Atoi(s)
}
```


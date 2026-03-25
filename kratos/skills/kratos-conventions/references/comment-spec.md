# Comment Spec

## 注释的目标

| 条件 | 做法 |
|------|------|
| 业务角色、边界、生命周期职责 | 写注释 |
| 代码已显而易见 | 不重复解释 |

```go
// ✅ 表达职责
// AccountUseCase 负责编排账户相关业务流程。
type AccountUseCase struct{}
```

```go
// ❌ 重复代码行为
// 设置 ID 为 req.Id
id := req.Id
```

---

## 注释落点

| 条件 | 做法 |
|------|------|
| Service / UseCase / Repo / Provider | 保留简短职责注释 |
| Ent schema / field | 保留业务语义注释 |
| 长流程关键分段 | 用少量章节注释 |

```go
// ✅ ProviderSet 服务提供集合
var ProviderSet = wire.NewSet(NewAccountService)
```

```go
// ✅ Ent 字段注释，使用中文注释
field.String("status").Comment("账户当前审核状态")
```

```go
// ✅ 流程分段注释
// 先校验提交状态，再进入审核流转。
if err := u.validateSubmit(ctx, account); err != nil {
    return err
}
```

---

## 生成物边界

| 条件 | 做法 |
|------|------|
| proto / ent 生成物注释 | 通过源文件生成 |
| 直接改生成物 | 不允许 |

```text
// ❌ 直接修改 *.pb.go 注释
```

---

## 组合场景

```go
// AccountRepo 定义账户聚合的数据访问边界。
type AccountRepo interface{}

// AccountUseCase 负责编排账户相关业务流程。
type AccountUseCase struct{}

// 先校验提交状态，再进入审核流转。
if err := u.validateSubmit(ctx, account); err != nil {
    return err
}
```

这个组合场景同时满足：

- 对象级注释说明职责
- 流程注释只出现在关键转折点
- 没有用注释补救糟糕命名

---

## 常见错误模式

```go
// ❌ 注释比代码还过时
```

```go
// ❌ 注释解释显而易见的赋值
```

```text
// ❌ 直接手改生成物注释
```

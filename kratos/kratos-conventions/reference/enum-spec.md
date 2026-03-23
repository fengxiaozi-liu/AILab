# Enum Spec

## 枚举定义

```go
// ✅ typed 枚举，值从 1 开始（0 保留为零值/无效）
type AccountOpenStatus uint32

const (
    AccountOpenStatusInit AccountOpenStatus = 1 + iota
    AccountOpenStatusFilling
    AccountOpenStatusFirstChecking
    AccountOpenStatusSecondChecking
    AccountOpenStatusOpened
)

// ❌ 裸数字，switch 分支无法检测遗漏
const (
    StatusInit     = 1
    StatusFilling  = 2
)

// ❌ 裸字符串替代枚举
if account.Status == "pending" { ... }
```

---

## switch 完整性

```go
// ✅ 包含 default，能捕获未处理值
func (u *AccountUseCase) canReview(status openenum.AccountStatus) bool {
    switch status {
    case openenum.AccountStatusPending, openenum.AccountStatusChecking:
        return true
    case openenum.AccountStatusRejected, openenum.AccountStatusApproved:
        return false
    default:
        return false  // ✅ 明确返回，不静默忽略
    }
}

// ❌ 只处理部分已知状态，漏掉新增枚举值
switch status {
case openenum.AccountStatusPending:
    return true
case openenum.AccountStatusApproved:
    return false
}
// 没有 default：新增枚举值后静默返回零值
```

---

## 枚举与 Localize

```go
// ✅ 枚举自带 Localize 方法，展示文案统一来自枚举
func (s AccountOpenStatus) Localize(localizer *i18n.Localizer) string {
    switch s {
    case AccountOpenStatusInit:
        return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
            DefaultMessage: &i18n.Message{ID: "ENUM_OPEN_STATUS_INIT", Other: "初始化"},
        })
    case AccountOpenStatusFilling:
        return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
            DefaultMessage: &i18n.Message{ID: "ENUM_OPEN_STATUS_FILLING", Other: "填写中"},
        })
    default:
        return ""
    }
}
```

---

## 状态流转

```go
// ✅ 状态流转基于枚举，错误语义明确
switch account.OpenStatus {
case openenum.AccountOpenStatusFirstChecking:
    account.OpenStatus = openenum.AccountOpenStatusSecondChecking
case openenum.AccountOpenStatusSecondChecking:
    account.OpenStatus = openenum.AccountOpenStatusOpened
default:
    return openerror.ErrorReviewStatusInvalid(ctx)
}

// ❌ 用数字做状态流转，意义不明
if account.Status == 3 {
    account.Status = 4
}
```

---

## 枚举变更联动

| 变更类型 | 必须同步检查 |
|----------|-------------|
| 新增枚举值 | proto 枚举定义、DB 列允许值、所有 switch default 覆盖 |
| 重命名枚举值 | proto、DB 中已存储的值、所有引用点 |
| 删除枚举值 | DB 历史数据处理方案、协议兼容性 |

---

## 常见错误模式

```go
// ❌ proto 枚举值和 Go 常量不一致（proto 用 0 开始，Go 用 1+iota）
// 导致序列化后值错位

// ❌ switch 无 default，新增枚举值后静默忽略
switch status {
case openenum.AccountStatusPending: ...
case openenum.AccountStatusApproved: ...
}

// ❌ 用裸字符串
if account.Type == "personal" { ... }
```

# Enum Reference

## 这个主题解决什么问题

说明业务状态、类型和值域如何使用枚举表达，以及枚举在 proto、Biz、DB 之间如何对齐。

## 适用场景

- 新增业务状态
- 调整状态机分支
- 统一不同层的状态表达

## 设计意图

Enum 参考的重点是解释状态和值域如何表达业务阶段，而不是只列一组常量。

- 枚举值会直接影响分支、序列化、筛选和展示文案。
- 先理解状态从何而来、会流向哪里，再补常量和 switch，更不容易漏分支。
- 枚举命名稳定后，错误语义、i18n 和协议字段也更容易统一。

## 实施提示

- 先画出状态集合和流转方向，再决定常量命名。
- 同步检查数据库字段、proto 枚举和业务分支是否共用同一套值域。
- 如果一个状态只在临时流程里出现，先确认它是否应成为长期枚举。

## 推荐结构

- proto、Biz、数据库映射使用同一语义名
- 状态流转优先围绕枚举展开

## 典型写法

```go
switch account.OpenStatus {
case openenum.AccountOpenStatusInit:
    ...
case openenum.AccountOpenStatusFilling:
    ...
default:
    ...
}
```

```go
func (u *AccountUseCase) canReview(status openenum.AccountStatus) bool {
    switch status {
    case openenum.AccountStatusPending, openenum.AccountStatusChecking:
        return true
    case openenum.AccountStatusRejected, openenum.AccountStatusApproved:
        return false
    default:
        return false
    }
}
```

## 项目通用枚举定义示例

```go
type AccountOpenStatus uint32

const (
    AccountOpenStatusInit AccountOpenStatus = 1 + iota
    AccountOpenStatusFilling
    AccountOpenStatusFirstChecking
    AccountOpenStatusSecondChecking
    AccountOpenStatusOpened
)

func (s AccountOpenStatus) Localize(localizer *i18n.Localizer) string {
    switch s {
    case AccountOpenStatusInit:
        return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
            DefaultMessage: &i18n.Message{ID: "ENUM_OPEN_STATUS_INIT", Other: "初始化状态"},
        })
    case AccountOpenStatusFilling:
        return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
            DefaultMessage: &i18n.Message{ID: "ENUM_OPEN_STATUS_FILLING", Other: "信息填充中"},
        })
    default:
        return ""
    }
}
```

## 状态流转示例

```go
switch account.OpenStatus {
case openenum.AccountOpenStatusFirstChecking:
    account.OpenStatus = openenum.AccountOpenStatusSecondChecking
case openenum.AccountOpenStatusSecondChecking:
    account.OpenStatus = openenum.AccountOpenStatusOpened
default:
    return openerror.ErrorReviewStatusInvalid(ctx)
}
```

## 常见坑

- 状态语义写成裸数字或裸字符串
- proto 和 Biz 中的枚举名称不一致
- 状态新增后漏改 switch 分支

## 相关 Rule

- `../rules/enum-rule.md`

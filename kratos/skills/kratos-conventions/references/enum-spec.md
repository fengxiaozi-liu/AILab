# Enum Spec

## 值域是否类型化

| 条件 | 做法 |
|------|------|
| 会参与 `switch`、状态流转、筛选、序列化 | 优先建 enum 或类型化常量 |
| 只是协议片段、Header、默认值 | 用具名常量，不强行提升为业务 enum |

```go
// ✅ 类型化状态
type ReviewStatus string

const (
    ReviewStatusPassed ReviewStatus = "reviewed-pass"
    ReviewStatusFailed ReviewStatus = "reviewed-fail"
)
```

```go
// ✅ 协议常量
const HeaderWebhookSignature = "X-Payload-Digest"
```

```go
// ❌ 裸业务字符串
if review.Answer == "GREEN" {
    return "reviewed-pass"
}
```

---

## 跨层联动

| 条件 | 做法 |
|------|------|
| 新增或修改 enum | 同步检查 DB、proto、业务分支、i18n |
| 只改一个调用点 | 视为未完成 |

```go
// ✅ switch 覆盖新增状态
switch answer {
case openenum.ReviewAnswerGreen:
    return openenum.ReviewStatusPassed
case openenum.ReviewAnswerRed:
    return openenum.ReviewStatusFailed
default:
    return openenum.ReviewStatusInvalid
}
```

```go
// ❌ 只改常量定义，不补分支
const ReviewStatusPaused ReviewStatus = "paused"
```

---

## 值域分层

| 条件 | 做法 |
|------|------|
| 状态、原因、结果 | 用业务 enum |
| URI、Header、分页默认值 | 用具名常量 |

```go
// ✅ 原因型常量
type ApplicantReason string

const ApplicantReasonAccountFlowCompleted ApplicantReason = "account-flow-completed"
```

```go
// ❌ 混放
const (
    ReviewStatusPassed      = "pass"
    HeaderWebhookSignature  = "X-Payload-Digest"
)
```

---

## 组合场景

```go
type ReviewAnswer string

const (
    ReviewAnswerGreen ReviewAnswer = "GREEN"
    ReviewAnswerRed   ReviewAnswer = "RED"
)

type ReviewStatus string

const (
    ReviewStatusPassed  ReviewStatus = "reviewed-pass"
    ReviewStatusFailed  ReviewStatus = "reviewed-fail"
    ReviewStatusInvalid ReviewStatus = "invalid-callback"
)
```

这个组合场景同时满足：

- 稳定业务值域被类型化
- 状态与输入值分层表达
- 分支和联动更容易补齐

---

## 常见错误模式

```go
// ❌ "GREEN" / "RED" 散落
if answer == "GREEN" { ... }
```

```go
// ❌ enum 与协议常量混在一起
```

```go
// ❌ 新状态未补 switch / i18n / proto
```

# Proto Helper Spec

## 协议辅助边界

| 条件 | 做法 |
|------|------|
| 稳定的 paging/sort/filter/time range 转换 | 可以放 `internal/pkg/proto` |
| 单一业务 reply 拼装 | 不下沉到公共 proto helper |

```go
// ✅
func ToPaging(in *base.Paging) *filter.Paging {
    ...
}
```

```go
// ❌
func BuildAccountReply(info *biz.Account) *v1.GetAccountReply { ... }
```

---

## 复用优先

| 条件 | 做法 |
|------|------|
| 新增 helper 前 | 先查是否已有 paging/排序/过滤/时间范围转换 |
| 只在单一 proto 中使用 | 留在对应 service/proto 邻近位置 |

---

## 常见错误模式

```text
// ❌ 把业务专用协议拼装沉到 internal/pkg/proto
```

# Util Spec

## 是否应该下沉

| 条件 | 做法 |
|------|------|
| 无业务语义、无状态、会被多处复用 | 可以下沉到 `util` |
| 只服务单一业务流程 | 留在业务层 |

```go
// ✅
func JoinBaseURL(baseURL, uri string) string {
    return strings.TrimRight(baseURL, "/") + uri
}
```

```go
// ❌ 业务专用逻辑
func BuildAccountReviewURL(account *biz.Account) string { ... }
```

---

## 文件拆分

| 条件 | 做法 |
|------|------|
| 能力主题明确 | 按主题拆文件，如 `url.go`、`base64.go` |
| 能力堆在一个大文件 | 不推荐 |

```text
// ✅
internal/pkg/util/url.go
internal/pkg/util/base64.go
```

```text
// ❌
internal/pkg/util/helper.go
```

---

## 最小创建流程

| 步骤 | 做法 |
|------|------|
| 1. 确认边界 | 先确认能力脱离业务语义后仍成立 |
| 2. 选择目录 | 复用现有子目录，或新建能表达边界的子目录 |
| 3. 创建最小文件 | 只放当前主题能力，例如 `url.go` |
| 4. 暴露最小 API | 先提供一个稳定函数，不预先造接口和工厂 |
| 5. 回接调用方 | 至少替换当前调用点，并确认复用价值 |

```text
// ✅ 最小创建流程
确认 JoinBaseURL 无业务语义
-> 放到 internal/pkg/util/url.go
-> 暴露 func JoinBaseURL(baseURL, uri string) string
-> 替换业务层重复实现
```

```text
// ❌ 过度创建流程
新建 internal/pkg/common/helper/
-> 先建 interface / manager / factory / options
-> 只有一个调用点
```

---

## 常见错误模式

```text
// ❌ util 里访问 repo/config/业务 context
```

```text
// ❌ 单一业务流程硬抽 util
```

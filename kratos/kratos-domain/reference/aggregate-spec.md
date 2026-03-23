# Aggregate Spec

## 对象复用 vs 新建

| 条件 | 做法 |
|------|------|
| 已有对象语义稳定，语义未变化 | 直接复用，不新建 |
| 需要裁剪/脱敏/跨协议投影 | 新建专用结构体 |
| 仅因传递场景（EventBus/ContextDTO）变化 | 禁止新建近义壳结构 |
| 跨 proto 边界投影 | 新建 proto message，不共用 biz 对象 |

```go
// ✅ EventBus 传递直接用聚合对象
u.eventBus.Publish(ctx, &eventbus.Event{Payload: account})
u.eventBus.Publish(ctx, &eventbus.Event{Payload: store})

// ❌ 仅为 EventBus 传递包装近义结构
type AccountAfterOpenEventPayload struct{ Account *Account }
u.eventBus.Publish(ctx, &eventbus.Event{Payload: &AccountAfterOpenEventPayload{Account: account}})

// ❌ 仅为 context 传递包装 DTO
type AccountFlowPageContextDTO struct{ Store *AccountFlowPageStore }
```

---

## 聚合根命名

```go
// ✅ 领域对象名表达业务概念
type Account struct { ... }
type AccountCollect struct { ... }
type AccountFlowPage struct { ... }

// ❌ 用动作词驱动命名
type ReviewResult struct { ... }
type OpenFlow struct { ... }
```

---

## 聚合根字段排序

```go
// ✅ 字段顺序：普通字段 → 时间字段 → Info 字段 → List 字段
type Account struct {
    ID           uint32                 `json:"id"`
    UserCode     string                 `json:"user_code"`
    OpenStatus   openenum.AccountOpenStatus `json:"open_status"`
    CreateTime   uint32                 `json:"create_time"`
    UpdateTime   uint32                 `json:"update_time"`
    // Info 字段
    FirstCheckUserInfo  *adminbiz.AdminUser `json:"first_check_user_info"`
    AccountCollectInfo  *AccountCollect     `json:"account_collect_info"`
    // List 字段
    AccountFlowPageList []*AccountFlowPage  `json:"account_flow_page_list"`
}

// ❌ 字段顺序混乱，Info/List 散落在普通字段中间
```

---

## 聚合根驱动的文件展开

```text
// ✅ 以聚合根为中心平级展开
Account 聚合根
├── internal/biz/account.go           (领域对象 + UseCase 接口)
├── internal/biz/account_usecase.go   (UseCase 实现)
├── internal/data/account.go          (Repo 实现)
├── api/admin/account/v1/account.proto
└── internal/service/account_service.go

// ❌ 按接口动作驱动展开（每个接口单独建文件）
├── internal/biz/get_account.go
├── internal/biz/review_account.go
├── internal/data/get_account_repo.go
```

---

## 组合场景

```go
// 完整，从聚合根定义到对象复用传递：

// 1. 定义聚合根，带 json tag 和字段排序
type Account struct {
    ID          uint32                    `json:"id"`
    OpenStatus  openenum.AccountOpenStatus `json:"open_status"`
    CreateTime  uint32                    `json:"create_time"`
    UpdateTime  uint32                    `json:"update_time"`
    CollectInfo *AccountCollect           `json:"collect_info"`
}

// 2. EventBus 直接传递，不包装
func (u *AccountUseCase) passReview(ctx context.Context, account *Account) error {
    return u.eventBus.Publish(ctx, &eventbus.Event{
        Topic:   openenum.LocalEventAccountAfterOpen.Value(),
        Payload: account,  // ✅ 直接用聚合对象
    })
}

// 3. Depend 直接传递，不包装
func (u *AccountUseCase) FillAdminUsers(ctx context.Context, list []*Account) error {
    return u.adminUserDepend.FillUsers(ctx, list)  // ✅ 传 list，不另建 DTO
}
```

---

## 常见错误模式

```go
// ❌ 近义壳结构
type AccountVO struct { Account *Account }
type EventPayload struct { Data interface{} }
type ContextDTO struct { Account *Account; Store *Store }

// ❌ 动作词命名聚合根
type ReviewFlow struct { ... }
type OpenProcess struct { ... }

// ❌ 缺少 json tag，参与事件投递但字段无法序列化
type Account struct {
    ID         uint32
    OpenStatus openenum.AccountOpenStatus
}
```

# Proto Helper Spec

## ParseXxx / BuildXxx 职责

| Helper 类型 | 职责 | 禁止 |
|-------------|------|------|
| `ParseXxx` | proto → 内部 filter/biz 对象 | 访问 DB / 调用 Repo |
| `BuildXxx` | 内部对象 → proto reply 字段 | 携带业务逻辑 |
| `ParseFilterConfig` | proto FilterConfig → filter.Option 链 | 业务校验 |

---

## TransField 双向转换

```go
// ✅ ParseTransField：proto → map
func ParseTransField(t *businessv1.TransField) map[baseenum.Language]string {
    if t == nil { return map[baseenum.Language]string{} }
    return map[baseenum.Language]string{
        baseenum.LanguageZH: t.Zh_CN,
        baseenum.LanguageTC: t.Tc,
        baseenum.LanguageEN: t.En,
    }
}

// ✅ BuildTransField：map → proto
func BuildTransField(m map[baseenum.Language]string) *businessv1.TransField {
    if m == nil { return &businessv1.TransField{} }
    return &businessv1.TransField{
        Zh_CN: m[baseenum.LanguageZH],
        Tc:    m[baseenum.LanguageTC],
        En:    m[baseenum.LanguageEN],
    }
}

// ❌ 在 Service 中手写 TransField 转换
reply.Name = &businessv1.TransField{
    Zh_CN: account.Name[baseenum.LanguageZH],  // ❌ 应用 BuildTransField
    Tc:    account.Name[baseenum.LanguageTC],
    En:    account.Name[baseenum.LanguageEN],
}
```

---

## Paging / Sort / TimeRange 解析

```go
// ✅ ParseSort：proto → filter.Sort（含驼峰→下划线转换）
func ParseSort(s *businessv1.Sort) filter.Sort {
    res := filter.Sort{Order: "desc", Field: "id"}
    if s == nil { return res }
    if s.Field != "" { res.Field = helper.HumpToSnake(s.Field) }
    if s.Order != "" { res.Order = s.Order }
    return res
}

// ✅ 在 Service 中用 ParseXxx 解析入参
func (s *AccountService) PageListAccount(ctx context.Context, req *v1.PageListAccountRequest) (*v1.PageListAccountReply, error) {
    pg := page.New(req.Paging)    // ✅ proto helper 转换
    sort := protohelper.ParseSort(req.Sort)  // ✅
    ...
}

// ❌ 在 Service 中直接算分页 offset
offset := (req.Paging.Page - 1) * req.Paging.PageSize  // ❌ 应用 page.New 统一转换
```

---

## 禁止在 proto-helper 中做业务组装

```go
// ❌ proto-helper 中访问 Repo
func BuildAccountReply(ctx context.Context, account *biz.Account) (*v1.AccountReply, error) {
    user, err := adminUserRepo.GetUser(ctx, account.CheckAdminUserID)  // ❌ helper 不访问 DB
    ...
}

// ❌ proto-helper 中做业务校验
func ParsePaging(p *businessv1.Paging) (*page.Page, error) {
    if p.PageSize > 100 {
        return nil, errors.New(400, "PAGE_SIZE_TOO_LARGE", "分页过大")  // ❌ 校验不属于 helper
    }
    ...
}
```

---

## 组合场景

```go
// 完整：Service 使用多个 ParseXxx / BuildXxx
func (s *AccountService) PageListAccount(
    ctx context.Context, req *v1.PageListAccountRequest,
) (*v1.PageListAccountReply, error) {
    pg := page.New(req.Paging)
    sort := protohelper.ParseSort(req.Sort)
    filter := protohelper.ParseFilterConfig(req.Filter)

    accounts, err := s.uc.PageListAccount(ctx, pg, sort, filter...)
    if err != nil { return nil, err }

    items := make([]*v1.AccountItem, 0, len(accounts))
    for _, a := range accounts {
        items = append(items, &v1.AccountItem{
            Id:   a.ID,
            Name: protohelper.BuildTransField(a.Name),  // ✅
        })
    }
    return &v1.PageListAccountReply{
        List:   items,
        Paging: protohelper.BuildPaging(pg),  // ✅
    }, nil
}
```

---

## 常见错误模式

```go
// ❌ Service 散落写 paging offset 计算
offset := (req.Page - 1) * req.PageSize

// ❌ Service 散落写 TransField 字段拼装
reply.Name = &businessv1.TransField{ Zh_CN: name["zh"] }

// ❌ BuildXxx 函数中携带业务逻辑
func BuildAccountReply(account *biz.Account) *v1.AccountReply {
    if account.Status == openenum.AccountStatusOpened {  // ❌ 业务判断不属于 helper
        ...
    }
}
```

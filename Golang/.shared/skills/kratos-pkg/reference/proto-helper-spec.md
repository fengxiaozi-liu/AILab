# Proto Helper 参考

## 这个主题解决什么问题

统一公共 proto 结构与内部 filter、枚举、多语言字段之间的双向转换，避免在 service、repo、gateway 中重复写 `Paging/Sort/TimeRange/TransField` 解析代码。

## 适用场景

- 解析或构造 `TransField`
- 解析或构造 `Paging`、`Sort`、`TimeRange`
- 解析或构造 `GroupBy`、`Compare`、`FilterConfig`

## 推荐结构或实现方式

- 使用 `ParseXxx` / `BuildXxx` 成对命名。
- helper 只依赖 `base/business proto`、`filter`、基础枚举和通用转换函数。

## 标准模板

```go
func ParseTransField(t *businessv1.TransField) map[baseenum.Language]string {
    if t == nil {
        return map[baseenum.Language]string{}
    }
    return map[baseenum.Language]string{
        baseenum.LanguageZH: t.Zh_CN,
        baseenum.LanguageTC: t.Tc,
        baseenum.LanguageEN: t.En,
    }
}

func BuildTransField(m map[baseenum.Language]string) *businessv1.TransField {
    if m == nil {
        return &businessv1.TransField{}
    }
    return &businessv1.TransField{
        Zh_CN: m[baseenum.LanguageZH],
        Tc:    m[baseenum.LanguageTC],
        En:    m[baseenum.LanguageEN],
    }
}
```

```go
func ParseSort(s *businessv1.Sort) filter.Sort {
    res := filter.Sort{Order: "desc", Field: "id"}
    if s == nil {
        return res
    }
    if s.Field != "" {
        res.Field = helper.HumpToSnake(s.Field)
    }
    if s.Order != "" {
        res.Order = s.Order
    }
    return res
}
```

## Good Example

- `ParseFilterConfig` 递归解析 relation 配置，把协议层配置统一转成内部 filter option。
- `BuildTimeRange`、`ParseTimeRange` 处理日期边界，不把这类逻辑散落在 service。

## 常见坑

- 在 service/repo 中重复写 `Paging/Sort/TimeRange` 解析代码
- 把业务 reply 组装逻辑写进 proto helper

## 相关 rule / 相关 reference

- `../rules/proto-helper-rule.md`
- `../../kratos-entry/reference/proto-spec.md`

# Schema Spec

## 通用 vs 业务耦合决策

| 特征 | 处理方式 |
|------|---------|
| 入参/出参均为通用类型（interface{}/reflect） | 放入 `internal/pkg/schema` |
| 仅服务于单一聚合根 | 不放入 schema，放在对应 Repo/biz 中 |
| 含业务校验逻辑 | 禁止放入 schema，放在 UseCase |
| 含业务配置判断 | 禁止放入 schema，放在 UseCase/Service |

---

## Extract 通用调用

```go
// ✅ schema.Extract 对任意结构体提取字段定义
fields, err := schema.Extract(SomeStruct{})
if err != nil {
    return err
}

// ✅ extract.go 保持无业务状态
// 输入：任意结构体值
// 输出：字段名 + 类型 + tag 信息

// ❌ schema helper 绑定特定聚合根
func ExtractAccountFields(account *biz.Account) ([]Field, error) {  // ❌ 业务绑定
    ...
}
```

---

## 文件拆分

```text
// ✅ 按能力主题拆分文件
internal/pkg/schema/
├── extract.go       // 字段提取
├── extract_test.go  // 对应测试
├── mapping.go       // 字段映射
└── mapping_test.go

// ❌ 所有功能堆在一个文件
internal/pkg/schema/
└── schema.go  // ❌ 职责不清晰
```

---

## 禁止在 schema 中混入业务逻辑

```go
// ❌ schema helper 中做业务校验
func Extract(v interface{}) ([]Field, error) {
    account, ok := v.(*biz.Account)
    if ok && account.Status == "" {  // ❌ 业务校验不属于 schema
        return nil, errors.New("status is required")
    }
    ...
}

// ❌ schema helper 中读取业务配置
func GetColumnNames(v interface{}) []string {
    cfg := config.GetTableConfig()  // ❌ schema 不依赖业务配置
    ...
}
```

---

## 组合场景

```go
// 完整：基础设施层调用 schema.Extract 做通用字段处理
// schema/extract.go
type Field struct {
    Name string
    Type reflect.Type
    Tag  reflect.StructTag
}

func Extract(v interface{}) ([]Field, error) {
    t := reflect.TypeOf(v)
    if t.Kind() == reflect.Ptr { t = t.Elem() }
    if t.Kind() != reflect.Struct {
        return nil, fmt.Errorf("schema: expected struct, got %s", t.Kind())
    }
    fields := make([]Field, 0, t.NumField())
    for i := 0; i < t.NumField(); i++ {
        f := t.Field(i)
        fields = append(fields, Field{Name: f.Name, Type: f.Type, Tag: f.Tag})
    }
    return fields, nil
}

// 调用方（基础设施层）
fields, err := schema.Extract(biz.Account{})
if err != nil { return err }
columns := make([]string, 0, len(fields))
for _, f := range fields {
    if col := f.Tag.Get("db"); col != "" {
        columns = append(columns, col)  // ✅ 通用列提取，不含业务逻辑
    }
}
```

---

## 常见错误模式

```go
// ❌ schema helper 绑定单一业务聚合根
func ExtractAccountSchema(account *biz.Account) []Column { ... }

// ❌ schema helper 中包含业务校验或配置判断
func Extract(v interface{}) ([]Field, error) {
    if isBlacklisted(v) { return nil, errors.New("blacklisted") }  // ❌
}

// ❌ schema 缺少测试，仅作为一次性工具
// internal/pkg/schema/extract.go 无对应 extract_test.go  ❌
```

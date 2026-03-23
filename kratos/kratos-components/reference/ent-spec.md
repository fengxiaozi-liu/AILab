# Ent Spec

## Schema 结构

```go
// ✅ 字段、edge、index、annotation、表注释齐全
func (Account) Fields() []ent.Field {
    return []ent.Field{
        field.Uint32("id"),
        field.String("user_code").MaxLen(64).Default("").Comment("user code"),
        field.Uint8("open_status").Default(0).Comment("open status"),
        field.Uint32("create_time").DefaultFunc(func() uint32 { return uint32(time.Now().Unix()) }),
        field.Uint32("update_time").DefaultFunc(func() uint32 { return uint32(time.Now().Unix()) }),
    }
}

func (Account) Annotations() []schema.Annotation {
    return []schema.Annotation{schema.Comment("account table")}
}

// ❌ 字段无注释，无 Default，表无注释
func (Account) Fields() []ent.Field {
    return []ent.Field{
        field.Uint32("id"),
        field.String("user_code"),
        field.Uint8("open_status"),
    }
}
```

---

## Upsert 冲突写法

| 数据库 | 正确写法 | 禁止写法 |
|--------|----------|----------|
| MySQL | `OnConflict().UpdateXxx()` | `OnConflictColumns(...)` |
| PostgreSQL | `OnConflictColumns(...).UpdateXxx()` | `OnConflict()` 不指定列 |

```go
// ✅ MySQL Upsert
err := r.data.Db.AccountFlowPage(ctx).Create().
    SetAccountID(info.AccountID).
    SetPageCode(info.PageCode).
    SetStatus(info.Status).
    OnConflict().
    UpdateStatus().
    UpdateUpdateTime().
    Exec(ctx)

// ❌ MySQL 中错误使用 OnConflictColumns
err := r.data.Db.AccountFlowPage(ctx).Create().
    SetAccountID(info.AccountID).
    OnConflictColumns(accountflowpage.FieldAccountID, accountflowpage.FieldPageCode).
    UpdateStatus().
    Exec(ctx)
```

---

## Edge 设计

```go
// ✅ edge 命名与查询侧 relation 对应，评估是否需要 eager-load
func (Account) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("account_collect", AccountCollect.Type),
        edge.To("account_flow_pages", AccountFlowPage.Type),
    }
}

// ⚠️ edge 只服务当前接口、不考虑 Repo relation 装配需求 → 需要回看 queryConfig 是否能正确 WithXxx
```

---

## Schema 变更联动

| 变更类型 | 必须联动 |
|----------|---------|
| 新增 field | 评估 Repo parseFilter、queryRelation 是否受影响 |
| 新增 edge | 评估 queryConfig WithXxx 和 serviceRelation 批量回填 |
| 新增 index | 评估查询排序和分页是否依赖新索引 |

---

## 组合场景

新增 AccountFlowPage schema，包含 field、edge、index、Upsert：

```go
func (AccountFlowPage) Fields() []ent.Field {
    return []ent.Field{
        field.Uint32("id"),
        field.Uint32("account_id").Default(0).Comment("关联Account"),
        field.String("page_code").MaxLen(64).Default("").Comment("页面编码"),
        field.Uint8("status").Default(0).Comment("页面状态"),
        field.Uint32("create_time").DefaultFunc(...),
        field.Uint32("update_time").DefaultFunc(...),
    }
}

func (AccountFlowPage) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("account", Account.Type).Ref("account_flow_pages").Field("account_id").Unique(),
    }
}

func (AccountFlowPage) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("account_id", "page_code").Unique(),
    }
}

func (AccountFlowPage) Annotations() []schema.Annotation {
    return []schema.Annotation{schema.Comment("account flow page table")}
}
```

---

## 常见错误模式

```go
// ❌ 手改 Ent 生成产物（ent/account.go）
// ❌ MySQL 中使用 OnConflictColumns
// ❌ 字段无 Default、无 Comment
// ❌ schema 改后未评估 Repo relation 装配影响
```

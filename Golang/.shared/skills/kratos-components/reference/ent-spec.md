# Ent Reference

## 这个主题解决什么问题

说明 Ent schema 中字段、edge、index、annotation、表注释等通常如何组织。

## 适用场景

- 新增 Ent schema
- 修改字段、索引、关系
- 调整表注释和命名

## 设计意图

Ent 参考主要解释 schema 如何承接领域模型和查询需求，而不是把它当作孤立的数据库定义文件。

- 字段、edge、index、annotation 一起决定后续 Repo 的查询形态和 relation 装配成本。
- 先理解 schema 的业务语义，再补字段和关系，通常比先写表结构更稳定。
- schema 组织清楚后，codegen、Repo 和测试的联动也更容易保持一致。

## 实施提示

- 先从领域对象提炼字段和关系，再映射到 Ent schema。
- 设计 edge 时顺手考虑查询侧是否需要 eager-load、列表装配或反向关联。
- 看到只服务单次展示的字段时，先判断它是否更适合留在协议层或派生视图里。

## 推荐结构

- `Fields`
- `Edges`
- `Indexes`
- `Mixin`
- `Annotations`

## 标准模板

```go
func (Account) Fields() []ent.Field {
    return []ent.Field{
        field.Uint32("id"),
        field.String("name"),
    }
}
```

## 代码示例参考

```go
func (Account) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("collects", AccountCollect.Type),
    }
}

func (Account) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("name"),
    }
}
```

## Upsert 示例

```go
func (r *accountFlowPageRepo) UpsertAccountFlowPage(ctx context.Context, info *biz.AccountFlowPage) error {
    nowTime := uint32(time.Now().Unix())

    return r.data.Db.AccountFlowPage(ctx).
        Create().
        SetAccountID(info.AccountID).
        SetPageCode(info.PageCode).
        SetStatus(info.Status).
        SetCreateTime(nowTime).
        SetUpdateTime(nowTime).
        OnConflict().
        UpdateStatus().
        UpdateTime().
        Exec(ctx)
}
```

如果仓库使用的是 PostgreSQL，再改用 `OnConflictColumns(accountflowpage.FieldAccountID, accountflowpage.FieldPageCode)` 一类按冲突列显式声明的写法。

## Good Example

- 索引名称、表注释、edge 含义清晰
- schema 改动后同步评估 relation 装配影响

## 注释说明

Ent schema 有比较稳定的注释写法：

- 表级注释使用 `schema.Comment("...")`
- 字段级注释使用 `.Comment("...")`
- 时间、状态、关联字段也会补充注释，方便数据库和生成物阅读

示例：

```go
field.Uint32("account_id").Default(0).Comment("关联Account")
field.String("page_code").MaxLen(64).Default("").Comment("页面编码")
field.Uint8("status").Default(0).Comment("页面状态")
```

## 常见坑

- 只改 schema，不回看 Repo relation
- 字段含义不清，注释缺失
- edge 设计只服务当前接口，不考虑长期边界

## 相关 Rule

- `../rules/ent-rule.md`
- `../../kratos-conventions/rules/comment-rule.md`

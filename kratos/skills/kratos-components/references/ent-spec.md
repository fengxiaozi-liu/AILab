# Ent Spec

## schema 与生成物

| 条件 | 做法 |
|------|------|
| 修改 schema / field / edge / index | 改 `internal/data/ent/schema` 并重新生成 |
| 生成物变化 | 通过源 schema 触发，不手改生成物 |

```text
// ✅
internal/data/ent/schema/account.go
go generate ./internal/data/ent
```

```text
// ❌
直接修改 ent 生成代码
```

---

## 更新与 upsert

| 条件 | 做法 |
|------|------|
| 单条更新 | 优先 `Update().Where(...IDEQ(id))` |
| 明确幂等写入 | 才使用 `OnConflictColumns(...)` |

```go
// ✅ 单条更新
func (r *accountRepo) UpdateStatus(ctx context.Context, id uint32, status uint8) error {
    return r.data.Db.Account(ctx).
        Update().
        Where(entaccount.IDEQ(id)).
        SetStatus(status).
        Exec(ctx)
}
```

```go
// ✅ 明确幂等键的 upsert
func (r *accountKycRecordRepo) UpsertByIdempotencyKey(ctx context.Context, info *biz.AccountKycRecord) error {
    return r.data.Db.AccountKycRecord(ctx).
        Create().
        SetIdempotencyKey(info.IdempotencyKey).
        OnConflictColumns(entaccountkycrecord.FieldIdempotencyKey).
        UpdateReviewStatus().
        Exec(ctx)
}
```

```go
// ❌ 默认用 UpdateOneID
r.data.Db.Account(ctx).UpdateOneID(id).SetStatus(status).Exec(ctx)
```

```go
// ❌ 无唯一索引也使用 OnConflict
OnConflictColumns(entaccountkycrecord.FieldStatus)
```

---

## 组合场景

```text
修改 schema
-> go generate ./internal/data/ent
-> 若 provider 有联动，再执行 wire ./cmd/server
-> go test ./... -run Repo
```

这个组合场景同时满足：

- schema 是事实源
- 更新与 upsert 边界清晰
- 生成和验证闭环

---

## 常见错误模式

```text
// ❌ 手改 ent 生成物
```

```go
// ❌ OnConflict 无唯一约束
```

```go
// ❌ 把 upsert 当成普通更新模板
```

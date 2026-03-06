# Ent Rule

## Principles

- Ent schema 变更必须可生成、可回归、可追踪。

## Specification

- schema 固定放在约定位置。
- 明确字段、edge、index、annotation、表注释。
- schema 变化后联动评估 Repo relation 和 codegen 影响。
- Upsert 冲突处理按数据库类型区分：MySQL 优先使用 `OnConflict()`，PostgreSQL 才使用 `OnConflictColumns(...)`。

## Prohibit

- 禁止手改 Ent 生成产物。
- 禁止在 MySQL 场景下默认使用 `OnConflictColumns(...)` 作为冲突写法。

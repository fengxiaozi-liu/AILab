# Repo Rule

## Principles

- Repo 统一负责查询、relation 装配和数据访问边界。

## Specification

- Repo 命名按实体或聚合边界命名。
- 查询实现按 `parseFilter` / `queryConfig` / `queryRelation` / `serviceRelation` 组织。
- 分页必须稳定排序。
- 远程 relation 必须批量收集、批量查询、批量回填。
- 单条更新优先使用 `Update().Where(xxx.IDEQ(id))` 一类写法，避免 `UpdateOneID(id)` 带来的额外查询或不必要 ORM 行为。

## Prohibit

- 禁止 N+1。
- 禁止 Service/UseCase 补查 relation。
- 禁止把“必须存在”的 not found 传染成 `nil, nil`。
- 禁止在普通 Repo 更新路径中把 `UpdateOneID(id)` 作为默认写法。

# Repo Rule

## Principles

- Repo 负责查询组装、relation 装配和数据访问边界。

## Specification

- Repo 按实体或聚合边界命名。
- 查询按 `parseFilter`、`queryConfig`、`queryRelation`、`serviceRelation` 组织。
- 分页必须保持稳定排序。
- 远程 relation 必须批量收集、批量查询、批量回填。
- 单条更新优先使用 `Update().Where(xxx.IDEQ(id))` 这类写法，避免 `UpdateOneID(id)` 带来的额外查询或不必要的 ORM 行为。
- 普通查询路径不要把前置参数校验作为默认策略，应由查询结果和错误语义决定返回。
- 返回值语义必须稳定，不能让调用方依赖 `if obj != nil` 一类补判来猜测状态。

## Prohibit

- 禁止引入 N+1 查询。
- 禁止把 relation 查询补回到 Service 或 UseCase。
- 禁止把必须存在的 not found 传染成 `nil, nil`。
- 禁止在普通 Repo 更新路径中默认使用 `UpdateOneID(id)`。
- 禁止使用 `nil, nil` 或其他模糊返回，把状态猜测压给调用方。

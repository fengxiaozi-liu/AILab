# Naming Rule

## Principles

- 同一领域概念在各层保持同名同义。

## Specification

- Repo、UseCase、实体、聚合命名按领域对象命名。
- 命名应稳定、可搜索、可表达边界。
- 聚合根 DTO、实体 DTO 参与协议对齐、JSON 编解码或事件投递时，字段应补全 `json` tag，并使用 `snake_case`。
- 聚合根语境已经明确时，主对象 ID 字段或入参优先命名为 `id`，不要重复追加聚合根前缀。

## Prohibit

- 禁止动作词驱动领域对象命名。
- 禁止同义多名。
- 禁止在聚合根语境已经明确的文件、消息或方法中，把主对象 ID 重复命名为 `account_id`、`accountId` 一类前缀形式。

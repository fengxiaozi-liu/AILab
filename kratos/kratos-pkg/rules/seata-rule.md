# Seata Rule

## Principles

- `seata` 封装负责统一事务入口和传播控制，不负责业务补偿编排。

## Specification

- 事务封装通过显式配置和 callback 承接业务执行。
- 传播级别、事务名、超时、重试等由配置表达，不在业务逻辑中散落硬编码。

## Prohibit

- 禁止把具体业务分支、补偿逻辑、领域状态判断写入 seata helper。
- 禁止在多个模块重复实现相同的全局事务模板。

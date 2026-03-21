# Routing Rule

## Principles

- 任何 Kratos 代码修改前先路由，再实施。
- 项目类型判定优先看职责与语义，路径只作证据提示。

## Specification

- 先判定 `WorkDomain`：BaseService / 业务 / 网关。
- 根据本次改动的职责变化加载最小聚合 skill 集合。
- 无法直接判定时，优先读取 `.env.*` 中的 `SERVER_NAME`，再用目录和关键字兜底。

## Prohibit

- 禁止跳过 `kratos-patterns` 直接选择子 skill。
- 禁止仅凭文件路径决定业务边界。

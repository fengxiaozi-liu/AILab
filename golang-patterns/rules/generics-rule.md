# Generics Rule

## Principles

- 泛型用于真正可复用的通用结构和算法。

## Specification

- 优先在工具层使用泛型。
- 约束尽量收敛，避免无边界 `any` 污染。

## Prohibit

- 禁止为单点业务过度泛型化。

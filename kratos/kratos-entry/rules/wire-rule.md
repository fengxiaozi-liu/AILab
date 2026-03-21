# Wire Rule

## Principles

- 依赖注入关系清晰，provider 组织稳定。

## Specification

- 新增依赖或服务注册时检查 wire provider。
- wire 变更后进入生成与构建验证。

## Prohibit

- 禁止手改 `wire_gen.go`。

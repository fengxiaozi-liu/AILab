# Server Rule

## Principles

- 服务注册按约定位置落位，不散落。

## Specification

- gRPC/HTTP 服务实现与注册保持一致。
- 新增接口时同步检查 provider 和注册链路。

## Prohibit

- 禁止把注册逻辑散落在非约定位置。

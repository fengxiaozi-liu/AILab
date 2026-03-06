# Gateway Rule

## Principles

- Gateway 只做代理和协议适配，不承载业务编排。

## Specification

- 参数映射、响应适配、代理边界显式定义。
- 网关侧错误转换与限流策略保持统一。

## Prohibit

- 禁止把核心业务逻辑下沉到网关。

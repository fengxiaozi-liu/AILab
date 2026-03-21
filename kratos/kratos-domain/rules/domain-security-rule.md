# Domain Security Rule

## Principles

- 外部输入默认不可信，权限前置，日志脱敏。

## Specification

- 业务输入要做长度、格式、范围、枚举校验。
- 认证与授权在业务执行前完成。
- 错误和日志不得暴露 token、password、PII。

## Prohibit

- 禁止完整打印请求体或敏感标识。

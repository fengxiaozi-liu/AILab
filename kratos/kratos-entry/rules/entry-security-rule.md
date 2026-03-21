# Entry Security Rule

## Principles

- 接入层先鉴权、先校验、先限流。

## Specification

- 对外接口必须做认证和授权。
- 输入做长度、格式、范围校验。
- 上游调用和代理必须设置超时。

## Prohibit

- 禁止对外暴露底层错误细节。
- 禁止无保护的公开入口。

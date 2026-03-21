# Components Security Rule

## Principles

- 组件接入必须控制资源生命周期和外部依赖风险。

## Specification

- 连接、goroutine、句柄必须可释放。
- 外部依赖调用必须有超时和失败边界。
- 配置和依赖日志必须脱敏。

## Prohibit

- 禁止无退出条件 goroutine。
- 禁止拼接命令或不受控路径执行外部行为。

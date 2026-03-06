# UseCase Rule

## Principles

- UseCase 负责业务编排、事务边界、权限和状态流转。

## Specification

- 事务边界只放在 UseCase。
- relation 通过 opts 显式声明并透传 Repo。
- 权限校验在真正执行业务前完成。

## Prohibit

- 禁止 UseCase 直接写 DB 查询细节。
- 禁止 UseCase 手写 relation 装配。
- 禁止从隐式上下文偷读业务身份数据代替显式入参。

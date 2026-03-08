# UseCase Rule

## Principles

- UseCase 负责业务编排、事务边界、权限和状态流转。

## Specification

- 事务边界只放在 `UseCase`。
- 通过显式 opts 将 relation 加载需求透传给 Repo。
- 在执行业务动作前完成权限校验。
- 默认只处理 `err` 分支，不要仅为了防御就在上层重复做 Repo 已负责的普通参数校验。

## Prohibit

- 禁止在 UseCase 中直接写 DB 查询细节。
- 禁止在 UseCase 中手工装配 relation。
- 禁止依赖隐式上下文代替显式业务入参。
- 禁止对已经由 Repo 或其他下层负责的参数增加重复的防御式校验。

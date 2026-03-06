# Layer Rule

## Principles

- Service 只做协议适配，UseCase 只做编排，Repo 只做装配与访问。

## Specification

- 保持 biz/data/service/server 单一职责。
- relation 补查统一收口，不在 Service 和 UseCase 零散实现。

## Prohibit

- 禁止越层调用。
- 禁止循环依赖。

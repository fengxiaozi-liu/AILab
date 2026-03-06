# Util Rule

## Principles

- `util` 只沉淀真正通用、无状态、稳定复用的辅助函数。

## Specification

- 仅当函数跨模块重复出现、语义稳定且与业务无关时，才下沉到 `util`。
- 工具函数应保持短小、纯净、可测试。

## Prohibit

- 禁止把只服务单一业务流程的逻辑放入 `util`。
- 禁止新增 `common`、`helper`、`misc` 这类语义不清的入口。
- 禁止在 `util` 中访问 repo、config、外部服务或业务 context。

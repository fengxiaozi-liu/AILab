# Testing Rule

## Principles

- 业务变更必须配套验证，不只覆盖 happy path。

## Specification

- 一个实现文件只对应一个测试文件。
- 新增分支、边界条件、关键回归点必须补测试。

## Prohibit

- 禁止为同一实现新增多个分散测试文件，除非存在充分理由。
- 禁止只改逻辑不补验证。

# Schema Rule

## Principles

- `schema` 只负责通用结构提取和描述，不绑定具体业务流程。

## Specification

- `schema` 下的能力应围绕字段提取、结构分析、通用反射辅助等基础能力组织。
- 新增 schema helper 时，优先保持输入输出通用，不绑定具体 proto、具体聚合。

## Prohibit

- 禁止把业务校验或业务配置逻辑放入 `schema`。
- 禁止为一次性场景创建不可复用的 schema helper。

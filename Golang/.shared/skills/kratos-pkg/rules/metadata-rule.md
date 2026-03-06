# Metadata Rule

## Principles

- Kratos metadata 的读写必须集中封装，保持 key、方向和默认值一致。
- metadata 表达跨边界透传的基础信息，不表达业务聚合语义。

## Specification

- 所有 metadata 读写统一通过 `internal/pkg/metadata` 暴露的 `GetXxx/SetXxx` 完成。
- metadata key 使用稳定命名，不在调用方拼接或散落定义。
- 同一字段要明确来自 client、server 或两侧合并读取的优先级。

## Prohibit

- 禁止在 service、repo、depend、中间件中直接硬编码 metadata key。
- 禁止为同一语义创建多个近义 key。
- 禁止把业务字段、业务 DTO 直接透传为 metadata。

## Checklist

- 是否提供了成对的 `GetXxx/SetXxx`。
- 是否明确了 client/server 读取来源。
- 是否避免了直接硬编码 `x-md-*` key。

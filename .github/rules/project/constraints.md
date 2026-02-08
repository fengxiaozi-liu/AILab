# Constraints（全局约束）

> 本文件描述全局约束类规则：用于控制变更节奏、契约优先与风险边界，不承载技能教程与长流程。
> 结构：Principles / Specification / Prohibit / Demo（Good & Bad）。

## Principles

- 先确认再实施：未获授权不进行写操作
- 契约优先：跨边界/对外协议先定义契约再实现
- 变更可控：不做无关重构；引依赖需说明理由与风险

## Specification

- 写入任何文件前，必须先陈述修改计划并获得用户明确授权
- 对外可见数据（HTTP/gRPC/WebSocket/Event）必须先在 `.proto` 定义（Contract First）
- 变更 `.proto` / `wire` / `ent schema` 后必须走生成链路；生成物文件禁止手改
- 外部依赖调用必须设置超时；重试仅限幂等操作，并采用指数退避 + 抖动
- DB 查询避免全表扫描；分页必须稳定排序；避免 N+1

## Prohibit

- 禁止未需求明确时做大范围重构/改公共 API 行为
- 禁止用 Go struct 作为长期对外协议替代 `.proto`
- 禁止在请求热路径做全量扫描/高复杂度维护操作（例如遍历全表 key 的清理）

## Demo

Good:
```go
// Contract First: 先改 *.proto → 生成 → 再落地转换与实现
// Change Control: 仅实现需求相关改动，不做顺手重构
```

Bad:
```go
// BAD: 直接改 *.pb.go / wire_gen.go 让它编译通过
// BAD: 未经确认做大范围重构
```

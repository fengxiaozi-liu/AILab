# Coding Conventions（工程治理规则）

> 本文件描述工程治理类规则：用于约束产出质量与一致性，不承载技能教程。
> 结构：Principles / Specification / Prohibit / Demo（Good & Bad）。

## Principles

- 生成物不可手改：只改源文件并重新生成
- 测试不分散：同一实现文件只对应一个测试文件
- 禁止魔法值：业务语义常量/枚举集中治理
- 错误与日志可治理：跨边界错误必须统一；日志可检索、可定位、可脱敏且尽量简短

## Specification

- 测试文件组织：一个实现文件只对应一个 `*_test.go`，补用例必须追加到该文件
- 生成物治理：以下文件/目录禁止手改
  - `*.pb.go` / `*_grpc.pb.go` / `*_http.pb.go` / `*.validate.go`
  - `cmd/server/wire_gen.go`
  - `internal/data/ent/**`（除 schema）
- 业务语义字面量统一放 `internal/enum/**`，业务代码只引用枚举/常量
- 统一错误体系（跨边界/对外/跨层返回）
  - 必须返回 `internal/error/**` 定义的项目统一错误（不存在则先新增错误定义）
  - 对外返回与内部日志解耦：对外不直接透出底层原始错误细节
- 错误日志（Error 级别）必须包含
  - 方法名 tag（如 `[Package.Struct.Method]`）
  - 关键业务标识（脱敏后，如 ak/id）
  - 操作语义（如 `upstream_failed` / `db_query_failed`）
  - `err`（原始错误，仅内部日志；敏感信息需脱敏）
- 日志尽可能简短
  - 一行表达清楚“发生了什么 + 关键上下文 + err”，避免长段落解释
  - 优先结构化字段（tag/ak/id/cost_ms/err），避免重复信息
  - 高频路径避免日志风暴（可聚合/采样；避免在循环内逐条打印）

## Prohibit

- 禁止为同一实现文件新增多个测试文件（除 build tag/平台差异等充分理由）
- 禁止直接硬编码有业务语义的 path/reason/code/topic/header/key 前缀等
- 禁止跨边界直接使用 `fmt.Errorf()` / `errors.New()` 构造并返回“原生错误”（必须走 `internal/error/**`）
- 禁止在错误日志中打印敏感信息（token/secret/password/PII）或完整请求/响应体
- 禁止在 `.proto` / `config.yaml` 中新增/埋入日志内容（不要把解释性文字当日志输出）

## Demo

Good:
```go
// 生成物：修改 proto → 走生成链路 → 在业务层做 DTO 转换
// 错误：对外返回统一错误；内部日志记录脱敏上下文 + err
if err != nil {
    logger.Errorf("[service.User.Get] db_query_failed user_id=%d err=%v", userID, err)
    return nil, baseerror.ErrorFailed(ctx)
}
```

Bad:
```go
// BAD: 直接手改生成物文件
// BAD: 对外直接返回原生错误
if err != nil {
    return nil, fmt.Errorf("db query failed: %w", err)
}
```

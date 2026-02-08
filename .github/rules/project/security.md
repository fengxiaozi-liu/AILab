# Security（安全规则）

> 本文件描述安全类规则：用于约束输入处理、权限校验与敏感信息边界，不承载技能教程。
> 结构：Principles / Specification / Prohibit / Demo（Good & Bad）。

## Principles

- 外部输入不可信：默认拒绝，先校验再使用
- 最小权限：认证/授权必须前置
- 敏感信息不出域：日志/错误/输出都必须脱敏

## Specification

- 输入校验：长度/格式/范围/枚举/白名单；字符串需 trim；解析失败默认拒绝
- 认证与授权：对外接口必须鉴权；权限校验在业务逻辑之前
- 限流与配额：对外接口必须有保护；配额校验在资源消耗之前
- 资源治理：外部依赖必须超时；goroutine 必须可退出；连接必须释放

## Prohibit

- 禁止在日志/错误中输出密钥、token、密码、PII（身份证/完整手机号/银行卡）
- 禁止拼接 SQL 字符串；禁止用用户输入直接构造文件路径；禁止用字符串拼接执行系统命令

## Demo

Good:
```go
// 先校验再使用；对外错误返回统一错误；内部日志只记录脱敏字段
if !isValid(userInput) {
    return baseerror.ErrorBadRequest(ctx)
}
logger.Infof("[service.X.Y] invalid_input user_id=%s", mask(userID))
```

Bad:
```go
// BAD: 记录敏感信息/完整请求体；对外透出底层错误细节
logger.Errorf("req=%+v token=%s err=%v", req, token, err)
return fmt.Errorf("upstream failed: %w", err)
```

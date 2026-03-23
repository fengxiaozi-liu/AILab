---
name: kratos-conventions
description: 用于 Kratos 项目横切规范，包括错误语义、防御式编程、枚举、i18n、日志与注释一致性。
  适用于修改错误处理、not found 行为、nil nil 返回、吞错、重复参数校验、枚举/状态流转、i18n、日志或注释的场景。
  触发关键词：error、not found、nil nil、吞错、重复校验、校验职责、enum、状态、i18n、locale、logging、comment、一致性。
  DO NOT USE FOR：业务逻辑编排（→ kratos-domain）、接入层协议适配（→ kratos-entry）、基础组件设计（→ kratos-components）。
---

# Kratos Conventions

## 输入

- 必需：变更目标描述（规范类别 / 问题描述）
- 可选：涉及文件路径（service、usecase、repo、data、proto、ent schema）

缺少必需输入时，MUST 先向用户提问，不得猜测继续。

## 工作流

1. 识别本次变更类型：`error` / `defensive` / `enum` / `i18n` / `logging` / `comment`
2. 输出开始前结构化状态（见强制输出）
3. 按需加载对应参考文件（见参考文件清单）
4. IF 涉及 defensive / error → 同时加载 `reference/error-spec.md`
5. 执行规范整改或审查
6. 按提交前检查项逐项核对
7. 输出完成后结构化状态（见强制输出）

## 约束

### MUST
- 对外或跨层返回 MUST 使用项目统一错误，不直接外抛 `errors.New` / `fmt.Errorf`
- "必须存在"的查询 MUST 查不到时返回明确错误，不返回 `nil, nil`
- 错误 MUST 可观测、可定位
- 枚举变更 MUST 同步检查 proto、DB 列、switch 分支逻辑和状态流转完整性
- i18n key MUST 按语义稳定命名，不散落原始文案到业务流程
- 日志 MUST 使用 Kratos `log.Logger` / `log.Helper`，通过 `WithContext(ctx)` 透传上下文
- 错误日志 MUST 包含动作语义、关键业务标识和错误对象
- 注释 MUST 表达业务语义和边界，不重复代码字面意思

### MUST NOT
- MUST NOT 把"必须存在"的查询结果返回为 `nil, nil`
- MUST NOT 吞错，不使用 `_, _ = call()` 或只记日志后伪装成功
- MUST NOT 用裸字符串替代业务枚举
- MUST NOT 在 UseCase / Service 重复做 Repo 已负责的普通参数校验
- MUST NOT 手改 `active.*.toml` I18n 产物
- MUST NOT 打印 token、secret、password、PII、证件号或完整请求体
- MUST NOT 在高频循环、批处理热点路径无控制地输出日志
- MUST NOT 在 proto 文件中添加说明性注释、分隔线或装饰性注释
- MUST NOT 在生成物上手改注释

### SHOULD
- 涉及 proto / Ent schema / 对外契约改动 SHOULD 联动 `kratos-entry` 或 `kratos-components` 验证
- 注释 SHOULD 默认使用中文，只有外部接口名或英文语义更稳定时保留英文

## 强制输出

开始前输出：

```json
{
  "conventionScope": "error | defensive | enum | i18n | logging | comment",
  "semanticChange": "本次语义变更摘要"
}
```

完成后输出：

```json
{
  "notFoundConverted": true,
  "nilNilEliminated": true,
  "noSwallowedError": true,
  "noRedundantValidation": true,
  "enumBranchComplete": true,
  "i18nKeyStable": true,
  "noSensitiveLog": true,
  "commentSemanticsValid": true
}
```

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `reference/error-spec.md` | 修改错误码、错误返回语义、not found 处理、防御式编程 | 涉及 error / defensive 时 |
| `reference/enum-spec.md` | 修改枚举值、状态流转、switch 分支 | 涉及 enum 变更时 |
| `reference/i18n-spec.md` | 修改 i18n key、多语言文案 | 涉及 i18n 变更时 |
| `reference/logging-spec.md` | 修改日志结构、字段、脱敏策略 | 涉及 logging 变更时 |
| `reference/comment-spec.md` | 修改 proto、Ent、Service、Repo、UseCase 中的注释 | 涉及 comment 整改时 |

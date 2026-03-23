---
name: kratos-pkg
description: 用于 Kratos 项目 internal/pkg 通用基础能力的设计与实现，包括 context、metadata、middleware、proto helper、schema、seata、util 的边界、复用与实现规范。
  适用于新增或修改 internal/pkg 下公共能力、上下文透传、metadata 注入提取、中间件接入、proto 与 filter 转换、schema 提取、分布式事务封装或稳定工具沉淀的场景。
  触发关键词：internal/pkg、context、metadata、middleware、trans field、time range、paging、sort、filter、schema、seata、util。
  DO NOT USE FOR：业务逻辑编排（→ kratos-domain）、基础设施组件（→ kratos-components）、接入层协议（→ kratos-entry）。
---

# Kratos Pkg

## 输入

- 必需：变更目标描述（pkg 子域 / 能力描述）
- 必需：说明为什么该能力属于 `internal/pkg`（脱离业务后是否仍有复用价值）
- 可选：涉及的文件路径或现有近义实现

缺少必需输入时，MUST 先向用户提问，不得猜测继续。

## 工作流

1. 判断需求是否真正属于 `internal/pkg` 公共基础能力
2. 输出开始前结构化状态（见强制输出）
3. 识别本次变更子域：`context` / `metadata` / `middleware` / `proto` / `schema` / `seata` / `util`
4. 按需加载对应参考文件（见参考文件清单）
5. IF 能力被多层调用 → 评估 context / metadata / middleware / proto helper 之间的联动
6. 执行变更
7. 回看是否引入了近义重复结构、业务泄漏或过度封装
8. 输出完成后结构化状态（见强制输出）

## 约束

### MUST
- `internal/pkg` MUST 只承载跨层、稳定、可复用的基础能力，能力脱离具体业务服务后仍然成立
- `context` MUST 通过统一的 typed helper（`WithTenant`、`WithLocalize`、`WithViewer` 等）读写，不使用裸字符串 key
- 新增 context 字段时 MUST 同步评估是否需要联动 `metadata` 与 `middleware`
- `proto helper` MUST 只负责稳定的协议辅助转换，不负责业务 reply 拼装或聚合装配
- 同一能力被多个模块复用且边界明确、接口稳定时，MUST 才沉淀到 `internal/pkg`

### MUST NOT
- MUST NOT 把聚合对象、业务状态流转、业务校验、业务规则判断放入 `internal/pkg`
- MUST NOT 仅为单个调用点创建 `util`、`helper`、`common` 壳目录或壳函数
- MUST NOT 在 `internal/pkg` 中拼装业务 reply、聚合关系或下沉 repo/usecase 逻辑
- MUST NOT 在业务代码中直接使用裸字符串 key 读写 context
- MUST NOT 把大对象、聚合关系、临时业务结果直接塞进 context
- MUST NOT 同一语义同时由 context struct、metadata key、middleware 重复维护而不做收口

### SHOULD
- `util` SHOULD 只沉淀真正跨模块复用的无状态工具
- `internal/pkg` 中的类型、函数和目录命名 SHOULD 围绕基础能力本身，不围绕具体业务动作

## 强制输出

开始前输出：

```json
{
  "pkgScope": "context | metadata | middleware | proto | schema | seata | util",
  "foundationChange": "本次新增或调整的基础能力摘要",
  "boundaryCheck": "为什么这部分应该放在 internal/pkg"
}
```

完成后输出：

```json
{
  "noBusinessLeakage": true,
  "noRedundantImplementation": true,
  "contextMetadataMiddlewareConsistent": true
}
```

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `reference/context-spec.md` | 新增或修改上下文透传、viewer、租户注入 | 涉及 context 变更时 |
| `reference/metadata-spec.md` | 新增或修改 metadata 注入与提取 | 涉及 metadata 变更时 |
| `reference/middleware-spec.md` | 新增或调整通用中间件 | 涉及 middleware 变更时 |
| `reference/proto-helper-spec.md` | 新增或修改 TransField、Paging、Sort、TimeRange、FilterConfig 转换 | 涉及 proto helper 变更时 |
| `reference/schema-spec.md` | 新增或调整 schema 提取 | 涉及 schema 时 |
| `reference/seata-spec.md` | 新增或调整分布式事务封装 | 涉及 seata 时 |
| `reference/util-spec.md` | 新增稳定工具沉淀 | 涉及 util 时 |

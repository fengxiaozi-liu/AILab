---
name: kratos-pkg
description: |
  用于 Kratos 项目 internal/pkg 通用基础能力的设计与实现，包括 context、metadata、middleware、proto helper、schema、seata、util 的边界、复用与实现规范。适用于新增或修改 internal/pkg 下的公共能力、上下文透传、metadata 注入与提取、中间件接入、proto 与 filter 转换、schema 提取、分布式事务封装或稳定工具沉淀的场景。触发关键词包括 internal/pkg、context、metadata、middleware、trans field、time range、paging、sort、filter、schema、seata、util。
---

# Kratos Pkg

## 必读规则

- `./rules/pkg-boundary-rule.md`
- `./rules/context-rule.md`
- `./rules/metadata-rule.md`
- `./rules/middleware-rule.md`
- `./rules/proto-helper-rule.md`

## 按需参考

- Context：`./reference/context-spec.md`
- Metadata：`./reference/metadata-spec.md`
- Middleware：`./reference/middleware-spec.md`
- Proto Helper：`./reference/proto-helper-spec.md`
- Schema：`./reference/schema-spec.md`
- Seata：`./reference/seata-spec.md`
- Util：`./reference/util-spec.md`

## 读取顺序

先读取 `./rules/*.md` 明确 `internal/pkg` 的边界和禁止项，再按当前任务读取对应 `./reference/*.md` 获取推荐实现方式、模板与示例。

## 何时使用

- 新增或修改 `internal/pkg/context`、`internal/pkg/metadata`
- 新增或调整通用中间件、Context 注入、错误格式化、租户或用户透传
- 新增或修改 `internal/pkg/proto` 下的 `TransField`、`Paging`、`Sort`、`TimeRange`、`FilterConfig` 转换
- 新增或调整 `schema`、`seata`、`util` 等基础能力
- 判断某个通用能力应放入 `internal/pkg` 还是业务层时

## 核心约束

1. `internal/pkg` 只承载跨层、稳定、可复用的基础能力，不承载业务聚合、业务流程和具体领域规则。
2. `context`、`metadata`、`middleware` 必须协同设计，避免同一语义被重复存放在多个入口。
3. `proto helper` 负责稳定的协议辅助转换，不负责业务 `reply` 拼装或聚合装配。
4. `util` 只沉淀真正跨模块复用的无状态工具，避免把零散业务逻辑塞进工具包。
5. `schema`、`seata` 等基础封装要保持通用性，调用方通过显式参数和配置接入。

## 实施流程

1. 先判断需求是否真的属于 `internal/pkg` 公共基础能力。
2. 选择对应主题的 rule，确认边界、命名和禁止项。
3. 读取对应 reference，复用已有模式补充实现。
4. 若能力会被多层调用，评估 `context`、`metadata`、`middleware`、`proto helper` 之间的联动。
5. 回看是否引入了近义重复结构、业务泄漏或过度封装。

## 强制输出

开始前输出：

- `PkgScope:` 本次涉及的 pkg 子域，例如 `context / metadata / middleware / proto / seata / util`
- `FoundationChange:` 本次新增或调整的基础能力摘要
- `BoundaryCheck:` 为什么这部分应该放在 `internal/pkg`

提交前输出：

- 是否避免了业务逻辑泄漏到 `internal/pkg`（Yes/No）
- 是否复用了现有 `context / metadata / filter / proto helper` 能力而不是重复造轮子（Yes/No）
- 是否保持了 middleware、metadata、context 的一致性（Yes/No）

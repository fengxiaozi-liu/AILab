---
name: kratos-components
description: |
  用于 Kratos 组件与基础设施能力的使用规范，包括 Ent、EventBus、Crontab、Depend、Config。适用于新增或修改 schema、索引、edge、事件、监听器、定时任务、跨服务依赖封装、配置项、默认值或校验逻辑的场景。触发关键词包括 ent schema、eventbus、listener、cron、crontab、depend、InnerRPC、config、幂等、重试、批量装配。
---

# Kratos Components

## 必读规则

- `./rules/ent-rule.md`
- `./rules/depend-rule.md`
- `./rules/config-rule.md`

## 按需参考

- Ent：`./reference/ent-spec.md`
- EventBus：`./reference/eventbus-spec.md`
- Crontab：`./reference/crontab-spec.md`
- Depend：`./reference/depend-spec.md`
- Config：`./reference/config-spec.md`

## 读取顺序

先读 `./rules/*.md` 明确组件边界与治理要求，再按当前任务读取必要的 `./reference/*.md`。

按需读取以下参考文档：

- Ent：`./reference/ent-spec.md`
- EventBus：`./reference/eventbus-spec.md`
- Crontab：`./reference/crontab-spec.md`
- Depend：`./reference/depend-spec.md`
- Config：`./reference/config-spec.md`

## 何时使用

- 新增或修改 Ent schema、edge、index、annotation
- 引入或调整 EventBus 事件、发布和 Listener
- 新增或修改定时任务、补偿任务、调度策略
- 调整 InnerRPC 或跨服务依赖封装
- 新增或变更配置项、默认值、校验、环境差异

## 核心约束

1. 组件能力按职责收口，不把业务编排塞进基础设施层。
2. Depend 和 EventBus 都必须避免散落直连与逐条远程调用，优先批量收集和统一封装。
3. Crontab 必须显式考虑幂等、重试、并发控制和可观测性。
4. Ent schema 变更要同步评估 Repo relation 装配和生成物影响。
5. Config 必须可回滚、可校验、可观测，并同步更新示例或文档位置。

## 实施流程

1. 判断变更落在哪类组件职责内。
2. 读取对应组件 reference，仅加载必要文档。
3. 设计组件接入点、边界和失败处理。
4. 若涉及模型或注入变更，联动 `kratos-entry` 做生成验证。
5. 补充组件相关验证或回归检查。

## 强制输出

开始前输出：

- `ComponentScope:` 本次涉及的组件类别
- `InfraChange:` schema/event/cron/depend/config 的变更摘要
- `RiskControl:` 幂等、重试、批量化、校验等控制点

提交前输出：

- 是否避免散落依赖和 N+1 式远程调用（Yes/No）
- 是否评估组件失败场景与可观测性（Yes/No）
- 是否同步处理配置默认值/校验或 schema 影响（Yes/No）

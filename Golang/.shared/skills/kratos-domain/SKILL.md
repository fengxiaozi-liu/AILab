---
name: kratos-domain
description: |
  用于 Kratos 业务核心分层实现，包括聚合建模、分层职责、命名、UseCase 编排、Repo/Data 设计与业务测试。适用于修改业务逻辑、聚合边界、Repo/UseCase 行为、relation 装配、事务边界、命名或分层职责的场景。触发关键词包括 aggregate、entity、repo、usecase、biz、data、relation、transaction、N+1、测试。
---

# Kratos Domain

## 必读规则

- `./rules/aggregate-rule.md`
- `./rules/layer-rule.md`
- `./rules/repo-rule.md`
- `./rules/usecase-rule.md`
- `./rules/testing-rule.md`

## 按需参考

- 聚合建模：`./reference/aggregate-spec.md`
- 分层落位：`./reference/layer-spec.md`
- 命名方式：`./reference/naming-spec.md`
- Repo 实现：`./reference/repo-spec.md`
- UseCase 编排：`./reference/usecase-spec.md`
- 测试组织：`./reference/testing-spec.md`

## 读取顺序

先读 `./rules/*.md` 明确边界与约束，再按当前任务加载必要的 `./reference/*.md`，不要一次性全部加载。

按需读取以下参考文档，不要一次性全部加载：

- 聚合与边界：`./reference/aggregate-spec.md`
- 分层职责：`./reference/layer-spec.md`
- 命名规范：`./reference/naming-spec.md`
- Repo/Data：`./reference/repo-spec.md`
- UseCase：`./reference/usecase-spec.md`
- 测试规范：`./reference/testing-spec.md`

## 何时使用

- 新增或重构业务领域能力
- 调整聚合根、实体边界、领域命名
- 变更 `internal/biz` 或 `internal/data` 中的 UseCase/Repo 协同
- 出现越层调用、relation 散落、N+1、事务边界不清
- 业务行为变更后需要补测试

## 核心约束

1. 先定聚合与层边界，再定命名、Repo、UseCase 和测试，不要反过来用代码细节倒推领域模型。
2. Service 只做协议适配；UseCase 只做编排、事务和业务决策；Repo 统一负责查询、relation 装配和数据访问。
3. 同一领域概念在 aggregate/usecase/repo/data/service 间保持同一术语，避免一词多名。
4. relation 装配统一收口到 Repo；UseCase 不手写 relation 查询细节，Service 不补查 relation。
5. 业务行为变化必须补齐相应测试，覆盖新增分支和关键边界，而不是只覆盖 happy path。

## 实施流程

1. 定义聚合根、实体和关系边界。
2. 校验命名是否稳定、能表达领域含义。
3. 明确层职责和依赖方向，避免越层调用。
4. 设计 UseCase 的事务边界、权限、状态流转和 opts 策略。
5. 设计 Repo 接口与实现，按四段流程组织：`parseFilter` / `queryConfig` / `queryRelation` / `serviceRelation`。
6. 为业务变更补测试，覆盖边界条件和关键分支。

## 强制输出

开始前输出：

- `DomainScope:` 本次涉及的聚合根、UseCase、Repo
- `LayerPlan:` 本次涉及的层和依赖方向
- `NamingChange:` 是否新增或重命名关键领域符号
- `TestPlan:` 准备补哪些测试

提交前输出：

- 聚合边界是否清晰且无动作词驱动命名（Yes/No）
- 是否不存在越层调用和 Service/UseCase relation 补查（Yes/No）
- Repo 是否按四段流程组织并避免 N+1（Yes/No）
- 事务边界是否只放在 UseCase（Yes/No）
- 是否补齐与变更匹配的测试（Yes/No）

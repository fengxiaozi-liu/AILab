---
name: kratos-domain
description: 用于 Kratos 业务核心分层实现，包括聚合建模、分层职责、命名、UseCase 编排、Repo/Data 设计与业务测试。
  适用于修改业务逻辑、聚合边界、Repo/UseCase 行为、relation 装配、事务边界、命名或分层职责的场景。
  触发关键词：aggregate、entity、repo、usecase、biz、data、relation、transaction、N+1、测试。
  DO NOT USE FOR：接入层协议适配（→ kratos-entry）、基础组件设计（→ kratos-components）、横切规范（→ kratos-conventions）。
---

# Kratos Domain

## 输入

- 必需：变更目标描述（聚合 / UseCase / Repo / 测试）
- 可选：涉及的聚合根名称、层名称（biz / data / service）
- 可选：`specs/<feature>/tasks.md`

缺少必需输入时，MUST 先向用户提问，不得猜测继续。

## 工作流

1. 定义聚合根、实体和关系边界
2. 输出开始前结构化状态（见强制输出）
3. 按需加载对应参考文件（见参考文件清单）
4. 校验命名是否稳定、能表达领域含义
5. 明确层职责和依赖方向，避免越层调用
6. IF 涉及 Repo → 按 `parseFilter` / `queryConfig` / `queryRelation` / `serviceRelation` 四段流程组织
   IF 涉及 UseCase → 设置事务边界、权限、状态流转和 opts 策略
7. 为业务变更补测试，覆盖边界条件和关键分支
8. 输出完成后结构化状态（见强制输出）

## 约束

### MUST
- 先定聚合与层边界，再定命名、Repo、UseCase 和测试，不要用代码细节反推领域模型
- Service MUST 只做协议适配；UseCase MUST 只做编排、事务和业务决策；Repo MUST 统一负责查询、relation 装配和数据访问
- 事务边界 MUST 只放在 UseCase
- relation 装配 MUST 统一收口到 Repo，UseCase 和 Service 不手写 relation 查询细节
- Repo 查询 MUST 按 `parseFilter` / `queryConfig` / `queryRelation` / `serviceRelation` 四段流程组织
- 远程 relation MUST 批量收集、批量查询、批量回填，禁止 N+1
- "必须存在"的单对象查询 MUST 查不到时返回明确错误，不返回 `nil, nil`
- Repo 单条更新 MUST 优先使用 `Update().Where(xxx.IDEQ(id))`，避免 `UpdateOneID(id)` 额外查询
- 聚合命名 MUST 表达领域对象，不使用动作词驱动命名
- 同一领域概念 MUST 在各层保持同一术语，不同义多名
- 聚合根 DTO 参与协议对齐或事件投递时字段 MUST 补全 `json` tag，使用 `snake_case`
- 聚合根语境已明确时，主对象 ID MUST 优先命名为 `id`，不追加聚合根前缀
- 业务变更 MUST 补齐相应测试，覆盖新增分支和关键边界
- 对象跨层传递时 MUST 优先复用已有稳定的聚合对象或应用层对象，不新建近义传递结构

### MUST NOT
- MUST NOT 越层调用（如 Service 直接调用 Repo）
- MUST NOT 循环依赖
- MUST NOT 在 UseCase 中直接写 DB 查询细节
- MUST NOT 在 Service / UseCase 中补查 relation
- MUST NOT 引入 N+1 查询
- MUST NOT 把"必须存在"的 not found 传染成 `nil, nil`
- MUST NOT 在单个 Repo 方法中跨多个表完成整个聚合的更新或删除
- MUST NOT 把整个聚合的更新/删除编排封装进单个 Repo 方法
- MUST NOT 使用动作词驱动聚合根或实体命名
- MUST NOT 仅因传递场景变化就新建 `EventPayload`、`ContextDTO`、`XxxData`、`XxxVO` 近义壳结构

### SHOULD
- 测试 SHOULD 覆盖边界条件和关键分支，不只覆盖 happy path
- 一个实现文件 SHOULD 只对应一个测试文件

## 强制输出

开始前输出：

```json
{
  "domainScope": "涉及的聚合根、UseCase、Repo",
  "layerPlan": "涉及的层和依赖方向",
  "namingChange": "是否新增或重命名关键领域符号",
  "testPlan": "准备补哪些测试"
}
```

完成后输出：

```json
{
  "aggregateBoundaryClean": true,
  "noLayerViolation": true,
  "repoFourStageOrganized": true,
  "transactionInUseCaseOnly": true,
  "noN1Query": true,
  "testAdded": true
}
```

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `reference/aggregate-spec.md` | 聚合根、实体边界建模、对象复用决策 | 涉及聚合建模或对象传递时 |
| `reference/layer-spec.md` | 分层职责确认、越层检查 | 涉及分层或职责边界变更时 |
| `reference/naming-spec.md` | 领域命名、ID 字段命名、tag 规范 | 新增或重命名领域符号时 |
| `reference/repo-spec.md` | Repo 四段流程、relation 装配、N+1 防范 | 涉及 Repo/Data 实现时 |
| `reference/usecase-spec.md` | UseCase 事务边界、opts 策略、权限编排 | 涉及 UseCase 编排时 |
| `reference/testing-spec.md` | 测试组织、边界覆盖、测试文件结构 | 需要补测试时 |

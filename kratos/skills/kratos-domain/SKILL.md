---
name: kratos-domain
description: Kratos 业务核心实现与分层收敛。用于 aggregate、usecase、repo/data、transaction boundary、relation 收口、biz test、越层调用修复。不要用于 proto/gateway、wire/codegen、服务注册或纯代码审查。
---

# Kratos Domain

## 何时使用

- 新增或重构业务能力。
- 调整聚合根、实体归属、关系边界。
- 调整 `internal/biz` / `internal/data` 的 Repo、UseCase 协作。
- 修复越层调用、relation 散落、事务边界不清。
- 业务变更后补领域测试和失败路径测试。

## 职责边界

- 本技能负责聚合建模、UseCase 编排、Repo/Data 设计和业务测试。
- 本技能不负责 proto、gateway、wire、codegen、服务注册，也不负责纯风格审查。

## 输入

- 必需：当前业务改动目标与范围，例如新增能力、重构 UseCase、调整 Repo、补测试。
- 可选：涉及的领域模块、文件路径、聚合名、UseCase 名或 Repo 名。
- 可选：需要重点判断的问题，例如聚合边界、事务边界、relation 收口、命名一致性。

缺少必需输入时，MUST 先从工作区和任务上下文补齐；仍无法判断业务边界时，再向用户提问，不得直接套模板。

## 工作流

### 收集证据与补齐输入

- 优先使用：改动文件路径、模块归属、现有聚合/用例/Repo 命名、事务边界位置、relation 组装位置、现有测试位置。
- 输入不足时：先在仓库中检索现有实现与边界约定；仍不足再向用户追问（不要猜）。

### 判定领域主题（只选一个主类型）

- `aggregate`：改聚合根、实体归属、关系边界。
- `layer`：判断 Service / UseCase / Repo 分层职责与依赖方向。
- `repo`：改查询、filter、relation、回填与数据装配收口。
- `usecase`：改编排、事务边界、领域动作组织。
- `testing`：补业务测试与失败路径覆盖（不要求在此步骤跑测试）。

若同时命中多个主题：以“业务行为变化的主改动点”为主类型，其他作为联动点处理。

### 先复用检查再读 references

- 先检索项目是否已有同类聚合/用例/Repo 的实现模式、目录落点与命名方式。
- 能复用就复用，避免新造平行体系（新聚合概念、新 usecase 编排入口、新 repo 组装方式）。

### 按类型按需加载 references

- `aggregate` -> 读 `references/aggregate-spec.md`
- `layer` -> 读 `references/layer-spec.md`
- `repo` -> 读 `references/repo-spec.md`；若新增/调整 Repo 接口或改动影响聚合边界，同步读 `references/aggregate-spec.md`
- `usecase` -> 读 `references/usecase-spec.md`；若新增/调整 UseCase 职责或编排影响聚合边界，同步读 `references/aggregate-spec.md`
- `testing` -> 读 `references/testing-spec.md`

### 实施领域改动（只做领域职责）

- 先收敛聚合与层级边界，再实施 UseCase、Repo 与测试改动。
- 保持 relation 组装收口在 Repo，业务动作收口在 UseCase，协议适配留在 Service。
- 优先复用已有聚合术语与跨层命名，避免概念漂移。

### 边界自检（不做测试验收）

- 是否以项目既有实现为主（给出证据：找到的目录落点/相似实现/命名模式）。
- 是否对照了本次主题对应的 references（至少确认 MUST/MUST NOT 没踩）。
- 是否保持分层一致：Service 不编排业务，UseCase 不补查 relation，Repo 收口 relation 组装。
- 若业务行为变化：是否补齐对应测试与失败路径（不要求在此处执行测试）。

## 约束

### MUST

- MUST 先定义聚合边界，再决定 UseCase、Repo 和测试结构。
- MUST 让 Service 只做协议适配，UseCase 只做编排与事务，Repo 统一负责查询与 relation 组装。
- MUST 让同一领域概念在 aggregate、usecase、repo、service、proto 间保持同一术语。
- MUST 在业务行为变化后补齐对应测试，而不是只覆盖 happy path。
- MUST 只加载当前主题直接相关的 references。

### MUST NOT

- MUST NOT 在 Service 层编排业务流程或补领域状态机。
- MUST NOT 在 UseCase 或 Service 层补查 relation，relation 组装必须收口在 Repo。
- MUST NOT 把业务问题转移成 proto、wire、codegen 或注册层问题处理。
- MUST NOT 只补表面路径，不补失败路径和边界测试。

### SHOULD

- SHOULD 优先用聚合和用例名称表达职责，而不是用页面名或接口名反推领域。
- SHOULD 在修改事务边界时同步检查 Repo 接口粒度和测试覆盖是否仍一致。
- SHOULD 在单一主题任务中只加载必要 reference，避免把技能变成通用手册。

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `references/aggregate-spec.md` | 改聚合根、实体归属、关系边界 | 涉及聚合建模时加载 |
| `references/layer-spec.md` | 判断 Service / UseCase / Repo 职责边界 | 涉及分层治理时加载 |
| `references/repo-spec.md` | 改查询、filter、relation、回填 | 涉及 Repo/Data 设计时加载 |
| `references/usecase-spec.md` | 改编排、事务、领域动作组织 | 涉及 UseCase 设计时加载 |
| `references/testing-spec.md` | 补业务测试和失败路径 | 涉及测试设计时加载 |

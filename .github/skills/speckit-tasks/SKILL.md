---
name: speckit-tasks
description: 基于 plan.md 生成按 Kratos 工作流排序的可执行任务清单。
---

# Spec Kit Tasks Skill（中文）

## 何时使用

- `plan.md` 完成后，需要将技术设计拆解为可执行任务。

## 输入

- 必需：`specs/<feature>/plan.md`、`specs/<feature>/spec.md`

## 依赖 Skill

- **kratos-patterns**：tasks agent 同时加载本 skill 与 kratos-patterns skill，后者提供各项目类型的「进行工作」步骤顺序，决定任务的 Phase 划分和执行顺序。

## 工作流

### §前置检查

1. 读取 `specs/<feature>/plan.md`，确认存在且 Status 为 Ready。
2. 读取 `specs/<feature>/spec.md`，提取 RQ 编号列表。
3. 若 plan 不存在或 Status 不为 Ready，终止并提示先运行 `/plan`。
4. 读取 `.specify/templates/tasks-template.md` 获取任务模板。

### §上下文提取

从 plan.md 提取：

1. **项目类型**（BaseService / 业务 / 网关）
2. **Phase 0 领域调研结果**：实体列表、状态流转、依赖服务
3. **Phase 1 技术设计表**：Schema / Proto / 枚举 / 异常 / 国际化 / 路由
4. **框架适配扩充**：补充的框架考量
5. **风险项**

从 spec.md 提取：

1. **RQ 列表**及其描述
2. **关键实体**
3. **约束与边界**

### §Phase 划分

根据项目类型，对照 **kratos-patterns** skill 的「进行工作」步骤，将任务划分为 Phase：

#### BaseService 项目

| Phase | 工作步骤 | 说明 |
|-------|---------|------|
| Phase 1 | Proto 接口设计与定义 | 设计/修改 proto 文件 |
| Phase 2 | 代码生成 | Proto → Go 代码 |
| Phase 3 | 枚举与异常定义 | 定义枚举和异常类型 |
| Phase 4 | InnerRPC 依赖包装 | 包装所需的依赖服务 |
| Phase 5 | 国际化 | 编写 i18n 内容 |
| Phase 6 | 收尾 | 文档、代码审查 |

#### 业务项目

| Phase | 工作步骤 | 说明 |
|-------|---------|------|
| Phase 1 | Ent Schema | 编写数据模型定义 |
| Phase 2 | Ent 代码生成 | 生成 Ent Go 代码 |
| Phase 3 | Biz 层实现 | 业务逻辑 |
| Phase 4 | Data 层实现 | 数据访问 |
| Phase 5 | Service 层实现 | API 实现 |
| Phase 6 | 服务注册 | 路由注册 |
| Phase 7 | Wire 代码生成 | 依赖注入 |
| Phase 8 | 测试 | 单元测试 + 集成测试 |

> **Phase 8 测试任务生成规则**：
> 从 plan.md 的「测试策略」章节读取测试目标表，为每个测试目标生成具体任务：
> - 每个需要单测的 biz/data 方法 → 一个任务（标注 `_test.go` 文件路径）
> - Mock 生成 → 一个前置任务（若需要 mockgen）
> - 集成测试 → 按用例场景拆分任务
> - 同层不同文件的测试任务可标记 `[P]`
>
> 测试文件路径约定：
> | 层 | 路径模式 |
> |----|--------|
> | Biz 单测 | `internal/biz/<domain>/<domain>_test.go` |
> | Data 单测 | `internal/data/<domain>/<domain>_test.go` |
> | Mock | `internal/biz/<domain>/mock_*_test.go` |
> | 集成测试 | `internal/biz/<domain>/integration_test.go` |
| Phase 9 | 收尾 | 文档、代码审查 |

#### 网关类项目

| Phase | 工作步骤 | 说明 |
|-------|---------|------|
| Phase 1 | Proxy 层实现 | 代理逻辑 |
| Phase 2 | 服务注册 | 路由注册 |
| Phase 3 | Wire 代码生成 | 依赖注入 |
| Phase 4 | 收尾 | 文档、代码审查 |

### §任务生成

在每个 Phase 内，根据 plan.md 的设计表生成具体任务：

1. **拆分粒度**：每个任务对应一个可独立完成的代码变更（一个文件或一组强关联文件）
2. **RQ 追踪**：每个任务标注其实现的 RQ 编号（来自 spec）
3. **并行标记**：同一 Phase 内操作不同文件且无依赖的任务标记 `[P]`
4. **文件路径**：每个任务描述中包含精确的产出文件路径（使用 Kratos 项目真实路径）

路径约定：

| 层 | 路径模式 |
|----|---------|
| Proto | `api/<domain>/<scope>/v1/<service>.proto` |
| Ent Schema | `internal/data/ent/schema/<entity>.go` |
| Biz | `internal/biz/<domain>/` |
| Data | `internal/data/<domain>/` |
| Service | `internal/service/<domain>/` |
| Enum | `internal/enum/<domain>/` |
| Error | `internal/error/<domain>/` |
| Server | `internal/server/` |
| i18n | `assets/i18n/active.*.toml` |
| Depend | `internal/biz/depend/`、`internal/data/depend/` |

### §任务格式

每条任务必须严格遵循：

```
- [ ] [TaskID] [P?] [RQ-xxx] 描述 + 文件路径
```

规则：
- **TaskID**：`T001` 起顺序编号，全文件唯一
- **`[P]`**：仅在可并行时添加（同 Phase 内、不同文件、无依赖冲突）
- **`[RQ-xxx]`**：标注该任务实现的 spec 需求编号；基建型任务（如代码生成、wire）可省略
- 描述必须包含明确的文件路径

示例：
```
- [ ] T001 [RQ-001] 定义 Portfolio Proto 服务接口 `api/portfolio/inner/v1/portfolio.proto`
- [ ] T002 执行 Proto 代码生成
- [ ] T005 [P] [RQ-001] 编写 Portfolio Ent Schema `internal/data/ent/schema/portfolio.go`
- [ ] T006 [P] [RQ-002] 编写 PortfolioItem Ent Schema `internal/data/ent/schema/portfolio_item.go`
```

### §写入规则

- 写入路径：`specs/<feature>/tasks.md`
- 使用 tasks-template.md 的章节结构
- Phase 顺序严格按照 kratos-patterns「进行工作」顺序

### §质量校验

| # | 校验项 | 标准 |
|---|--------|------|
| 1 | RQ 全覆盖 | 每条 RQ 至少被一个任务引用 |
| 2 | Phase 顺序正确 | 严格匹配 kratos-patterns 工作步骤顺序 |
| 3 | 依赖无环 | 后置 Phase 的任务不被前置 Phase 依赖 |
| 4 | 路径准确 | 文件路径符合项目真实结构 |
| 5 | 粒度合理 | 每个任务可在一次编辑中完成 |
| 6 | 并行标记正确 | `[P]` 任务之间确实无文件冲突和依赖 |
| 7 | 测试覆盖 | plan 测试策略中每个测试目标至少有一个对应任务 |

## 输出

- `specs/<feature>/tasks.md`（可执行任务清单）
- 统计摘要：总任务数、各 Phase 任务数、可并行任务数

---
name: speckit-implement
description: |
  用于 Speckit 研发流程的核心代码落地实施，包括读取 tasks.md 并逐 Phase 进行业务代码编写、脚手架工具链生成调用以及补充单元测试代码，完成闭环。适用于所有的设计、拆解与排位均已就绪，当前操作只需要纯粹进入“敲码落地与执行系统构建验证”阶段的场景。触发关键词包括 implement、实现任务、落地、写代码、编码。
---

# Spec Kit Implement Skill（中文）

## 何时使用

- 已有 tasks.md，准备按任务清单正式开发。
- 没有 tasks.md，但用户希望直接实现（以用户输入为准）；此时 implement 必须先生成一个“可执行的最小任务清单（草案）”用于推进与验收。

## 输入

- 必需：用户输入（要实现什么、范围/约束、可跳过项、验收口径）
- 可选：`specs/<feature>/tasks.md`
- 可选：`specs/<feature>/plan.md`
- 可选：`specs/<feature>/spec.md`
- 用户约束（范围、跳过的 Phase 等）

## 依赖 Skill

- **kratos-patterns**：implement agent 同时加载本 skill 与 kratos-patterns skill，后者提供各层代码编写规范（Ent Schema 怎么写、Proto 怎么定义、biz/data/service 分层结构等）。

## 工作流

### §前置检查

1. 若存在 `specs/<feature>/tasks.md`：进入按 tasks 执行模式。
2. 若 tasks 不存在：基于用户输入生成“最小任务清单（草案）”，创建`specs/<feature>/tasks.md` 再进行实施执行
3. 若存在 `specs/<feature>/plan.md`：读取项目类型与设计要点作为执行约束；若不存在则由或仓库事实判定项目类型。

### §上下文加载

1. 若存在 tasks.md：解析 Phase 列表及顺序、任务 ID/描述/文件路径/标记、已完成任务（`[x]`）用于断点续做。
2. 若 plan/spec 存在：读取设计/边界/验收口径作为执行约束；若不存在：从用户输入提炼“需求点清单（临时编号）+ 验收口径”。
3. 执行并遵循 `kratos-patterns` 路由：点名加载本次变更涉及的 `kratos-*` 子技能（以及需要的 `golang-patterns`），用它们的 MUST 规则约束实现。

### §逐 Phase 执行

**按 Phase 逐步推进，每个 Phase 完成后暂停汇报、等用户确认再继续。**

每个 Phase 内：

1. **加载对应规范**：执行 `kratos-patterns` 路由并点名加载本 Phase 涉及的 `kratos-*` 子技能（必要时 `golang-patterns`），以这些子技能的 MUST 规则作为执行依据
2. **逐任务执行**：
   - 根据 plan.md 设计表的具体设计 + kratos-patterns 规范，编写代码
   - 同 Phase 内 `[P]` 标记的任务顺序无关，按序执行即可
   - 每完成一个任务，在 tasks.md 标记 `[x]`
3. **Phase 完成汇报**：
   - 列出本 Phase 完成的任务
   - 标注产出文件
   - 提示下一个 Phase 的内容
   - 等待用户确认

#### 代码生成 Phase 特殊处理

以下 Phase 涉及代码生成（Ent codegen / Proto codegen / Wire codegen）：
- 执行完代码生成命令后，必须验证生成结果无报错
- 若生成失败，停止并报告错误，等用户决策
- 生成成功后再继续后续 Phase

### §断点续做

若 tasks.md 中已有 `[x]` 标记的任务：
1. 跳过已完成任务
2. 从第一个未完成任务所在的 Phase 开始
3. 汇报续做起点

### §错误处理

- 任务执行失败 → 立即停止，报告错误详情，等用户决策
- 不跳过失败任务继续后续（Kratos 分层有严格依赖链）
- 用户可选：修复后重试 / 跳过该任务 / 终止实施

### §完成校验

所有 Phase 执行完毕后：

| # | 校验项 | 标准 |
|---|--------|------|
| 1 | 任务全完成 | tasks.md 中所有任务均已 `[x]` |
| 2 | 需求覆盖 | spec 存在：每条 RQ 都有对应代码实现；spec 不存在：每条需求点清单都有对应实现或明确不做 |
| 3 | 设计一致 | plan 存在：代码与 plan.md 设计一致；plan 不存在：代码与“最小任务清单（草案）/用户确认的约束”一致 |
| 4 | 编译通过 | 代码可编译（`go build`） |

## 输出

- 代码实现改动
- 更新后的 `specs/<feature>/tasks.md`（所有任务打勾）

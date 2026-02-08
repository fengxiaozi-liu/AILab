---
name: speckit-implement
description: 按 tasks.md 逐 Phase 执行实现，完成任务闭环。
---

# Spec Kit Implement Skill（中文）

## 何时使用

- tasks.md 已就绪，analyze 质量门通过（或用户决定跳过），准备正式开发。

## 输入

- 必需：`specs/<feature>/tasks.md`、`plan.md`、`spec.md`
- 用户约束（范围、跳过的 Phase 等）

## 依赖 Skill

- **kratos-patterns**：implement agent 同时加载本 skill 与 kratos-patterns skill，后者提供各层代码编写规范（Ent Schema 怎么写、Proto 怎么定义、biz/data/service 分层结构等）。

## 工作流

### §前置检查

1. 确认 `specs/<feature>/tasks.md`、`plan.md`、`spec.md` 均存在。
2. 若任一缺失，终止并提示先运行对应命令。
3. 读取 plan.md 中的**项目类型**。

### §上下文加载

1. 解析 tasks.md：
   - Phase 列表及顺序
   - 每个任务的 ID、描述、文件路径、并行标记 `[P]`、RQ 引用
   - 已完成的任务（`[x]`），用于断点续做
2. 读取 plan.md：
   - Phase 1 设计表（Schema/Proto/枚举/异常/国际化等）
   - 框架适配扩充内容
3. 加载 kratos-patterns skill 对应项目类型的「准备工作」—— 读取所有 reference 规范文档

### §逐 Phase 执行

**按 Phase 逐步推进，每个 Phase 完成后暂停汇报、等用户确认再继续。**

每个 Phase 内：

1. **加载对应规范**：根据当前 Phase 的工作步骤，读取 kratos-patterns 中对应的 reference 文档
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
| 2 | RQ 全覆盖 | spec 中每条 RQ 都有对应的代码实现 |
| 3 | 设计一致 | 代码与 plan.md 设计表一致 |
| 4 | 编译通过 | 代码可编译（`go build`） |

## 输出

- 代码实现改动
- 更新后的 `specs/<feature>/tasks.md`（所有任务打勾）

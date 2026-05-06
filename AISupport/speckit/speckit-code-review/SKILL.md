---
name: speckit-code-review
description: 按 spec、plan、tasks 审查实现质量并生成 review.md。用于实现后复审或提交流程前。
---

# Spec Kit Code Review Skill

## 何时使用

- `implement` 执行完成后做正式审查
- 修复一轮问题后做复审
- 提交前做风险检查、回归检查、完成度检查

## 职责边界

- 在代码审查的时候以不遗漏为准则，避免二次审查
- `speckit-code-review` 负责通用代码审查流程
- 它永远先做通用审查，再按需叠加项目专属 skill
- 进入 `Feature review` 或 `Project-aware review` 时，若仓库已能识别项目类型，必须加载对应项目 skill 并依据该 skill 的约束补充审查
- 它是只读技能，不自动修改代码
- 它输出的 `review.md` 除了审查结论外，还必须维护可被 `/fix` 消费的全局问题列表

## 输入

- 优先：`specs/<feature>/tasks.md`
- 优先：`specs/<feature>/plan.md`
- 可选：`specs/<feature>/spec.md`
- 实际代码、测试、配置、迁移、脚本等相关文件
- 输出模板：本技能目录下 `./assets/review-template.md`

## 工作流

### 1. 判定审查模式

按上下文进入以下三种模式之一：

1. **Feature review**
   - 存在 `tasks.md` 与 `plan.md`
   - 目标：检查实现是否覆盖 spec/plan/tasks 约束
2. **Project-aware review**
   - 缺少完整 Speckit 上下文，但能识别项目类型
   - 目标：在通用审查维度上叠加项目适配规则
3. **General review**
   - 既无完整 Speckit 上下文，也无法可靠识别项目类型
   - 目标：退化为通用代码审查，并显式记录假设

### 2. 加载上下文

- Feature review：优先读取已完成任务关联的文件、测试和配置
- Project-aware review：先识别项目类型，再读取核心入口、业务逻辑、对外接口和测试
- General review：聚焦本次变更涉及的入口逻辑、数据访问、接口层和测试
- depend SKILL：加载项目相关的 skill 以及规范内容
- 如果仓库已识别出项目类型，输出中必须显式说明：
  - 已加载了哪个项目 skill
  - 为什么需要加载该 skill
  - 本次审查额外采用了哪些项目约束

### 3. 应用审查维度

始终以下列通用审查维度作为主干：

- 完成度
- 正确性
- 架构与规范
- 性能
- 安全
- 测试

同时在输出阶段优先沿用本技能目录下 `./assets/review-template.md` 的结构：

- 全局问题列表
- 审查模式与上下文
- 项目知识与规范依据
- 完成度
- 正确性
- 架构与规范
- 性能
- 安全
- 假设与未知项
- 结论与下一步

### 4. 输出前合并审查问题

在生成最终 `review.md` 前，必须先收集并合并本次审查的全部标准化 findings。
这些 findings 可以来自：

- 通用代码审查
- 项目审查 skill
- feature 上下文偏差检查

所有输入 findings 必须先标准化为统一结构后再参与合并。
禁止直接拼接不同 skill 的自然语言结论。

合并时必须满足：

- 对重复 findings 去重
- 若项目审查 finding 比通用 finding 更具体，则以更具体版本覆盖
- 最终所有保留 findings 统一编号、统一分类、统一落入全局问题列表
- 最终只输出一份合并后的 `review.md`

### 5. 输出报告

- 生成结构化审查报告`review.md`
- 报告结构优先遵循本技能目录下 `./assets/review-template.md`
- 审查发现的问题需要同时归并到 `review.md` 的全局问题列表中，并按已解决/未解决分组维护
- 若存在 feature 上下文，写入 `specs/<feature>/review.md`
- 若不存在 feature 上下文，则在当前任务范围内输出审查结果

## 输出

- 审查报告：`specs/<feature>/review.md`
- 提示: 
  1. 提示用户审查报告
  2. 提示用户报告无误后进入`/fix`

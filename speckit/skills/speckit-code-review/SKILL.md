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

- `speckit-code-review` 负责**通用代码审查流程**
- 它永远先做通用审查，再按需叠加项目或语言专属 skill
- 它是只读技能，不自动修改代码

## 输入

- 优先：`specs/<feature>/tasks.md`
- 优先：`specs/<feature>/plan.md`
- 可选：`specs/<feature>/spec.md`
- 实际代码、测试、配置、迁移、脚本等相关文件
- 参考：`./references/review-rubric.md`

## 工作流

### 1. 判定审查模式

按上下文进入以下三种模式之一：

1. **Feature review**
   - 存在 `tasks.md` 与 `plan.md`
   - 目标：检查实现是否覆盖 spec/plan/tasks 约束
2. **Project-aware review**
   - 缺少完整 Speckit 上下文，但能识别项目类型
   - 目标：在通用 rubric 上叠加项目适配规则
3. **General review**
   - 既无完整 Speckit 上下文，也无法可靠识别项目类型
   - 目标：退化为通用代码审查，并显式记录假设

### 2. 加载上下文

- Feature review：优先读取已完成任务关联的文件、测试和配置
- Project-aware review：先识别项目类型，再读取核心入口、业务逻辑、对外接口和测试
- General review：聚焦本次变更涉及的入口逻辑、数据访问、接口层和测试

### 3. 应用审查 rubric

始终以 `./references/review-rubric.md` 中的通用维度作为主干：

- 完成度
- 正确性
- 架构与可维护性
- 性能
- 安全
- 测试

如果识别到额外项目 skill 或语言 skill：

- 将其 `MUST` 规则作为补充输入
- 不允许用专属规则替代通用审查主干

### 4. 输出结论

- 生成结构化审查报告
- 若存在 feature 上下文，写入 `specs/<feature>/review.md`
- 若不存在 feature 上下文，则在当前任务范围内输出审查结果

## 输出

- 审查报告：`specs/<feature>/review.md` 或当前任务对应的审查结果
- 审查完成后：
  - 若存在高优先级问题，提示修复后重新 `/code-review`
  - 若审查通过，提示进入 `/summary`

---
name: speckit-implement
description: 按 tasks.md 分 Phase 实施并回写任务状态。用于正式编码阶段。
---

# Spec Kit Implement Skill

## 何时使用

- 已存在 `specs/<feature>/tasks.md`，准备进入正式实施阶段
- 用户明确要求按任务清单逐项落地实现

## 职责边界

- `speckit-implement` 只负责基于既有 `tasks.md` 实施并闭环
- 不负责在缺少 `tasks.md` 时直接开工
- 不擅自改写任务边界；若任务本身失真，应先向用户报告

## 输入

- 必需：`specs/<feature>/tasks.md`
- 可选：`specs/<feature>/plan.md`
- 可选：`specs/<feature>/spec.md`
- 可选：当前任务直接相关的代码、测试、脚本、配置
- 参考：`./references/phase-execution.md`
- 模板：`./templates/issue-template.md`

## 核心原则

- `tasks.md` 是唯一执行清单
- 一次只推进一个 Phase
- 当前 Phase 未闭环前，不进入下一 Phase 的具体实现
- 允许在同一 Phase 内按模块连续实现多个强关联任务
- 任务状态必须真实回写，不允许事后集中“补打钩”

## 工作流

### 1. 前置检查

1. 确认 `tasks.md` 存在
2. 识别当前第一个未闭环的 Phase
3. 按需加载 `plan.md` / `spec.md` 中与当前 Phase 直接相关的约束
4. 按需读取 `./references/phase-execution.md`

若缺少 `tasks.md`，立即停止并提示补齐前置产物。

### 2. Phase 执行

围绕当前 Phase 执行以下循环：

1. 扫描任务描述、目标文件、验收动作
2. 判断是否命中其他相关 skill，并在实施前加载
3. 实现当前批次最相关的一组任务
4. 做必要的本地验证
5. 将已真实落地的任务更新为 `[✅️]`

### 3. 阻断处理

若出现以下情况，停止推进并向用户报告：

- 缺少必要上下文，无法安全实施
- 存在严重设计冲突，无法自行裁决
- 当前任务明显需要补充任务清单或重排边界
- 关键依赖、脚本或必需 skill 无法正常工作

## 输出

- 代码实现改动
- 更新后的 `specs/<feature>/tasks.md`
- 当所有 Phase 闭环后，创建 `specs/<feature>/issue.md`
  - 生成前明确说明本次 `issue.md` 将基于 `./templates/issue-template.md` 生成
  - `issue.md` 只记录本次需求实施过程中暴露的问题，不引入无关需求
  - 告知人工审查者只需填写：问题标题、现象、期望行为、复现或证据
  - 其他字段由后续 AI 修复流程补全
- 如全部完成，提示进入 `/code-review`

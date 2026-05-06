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
- `tasks.md` 是唯一执行清单
- 禁止当前 Phase 未闭环前进入下一 Phase 的具体实现
- 标记任务项完成状态仅使用 `[✅️]`，不允许使用 `[x]`
- 仅当当前 Phase 下的所有任务都标记为 `[✅️]` 时，视为阶段完成

## 输入

- 必需：`specs/<feature>/tasks.md`
- 可选：`specs/<feature>/plan.md`
- 可选：`specs/<feature>/spec.md`
- 可选：当前任务直接相关的代码、测试、脚本、配置
- 参考：`./references/phase-execution.md`

## 工作流

### 1. 前置检查

1. 确认 `tasks.md` 存在
2. 识别当前第一个未闭环的 Phase
3. 按需加载 `plan.md` / `spec.md` 中与当前 Phase 直接相关的约束
4. 按需读取 `./references/phase-execution.md`

若缺少 `tasks.md`，立即停止并提示补齐前置产物。

### 2. Phase 执行

围绕当前 Phase 执行以下循环：

1. 锁定“第一个仍包含未完成任务的 Phase”作为当前唯一执行上下文；在该 Phase 闭环前，不得实现后续 Phase 的具体任务。
2. 扫描当前 Phase 内的任务描述、目标文件、验收动作，并加载与当前 Phase 直接相关的执行依据。
3. `speckit-implement` 作为外层流程，只负责锁定当前 Phase、组织实施顺序、约束确认与回写时机，以及控制 `tasks.md` 的状态推进。
4. 当前 Phase 的实现必须遵循已加载的执行依据；`speckit-implement` 不扩展、不改写这些依据，仅负责在 Phase 内推动其被执行。
5. 在当前 Phase 开始执行时，先输出当前阶段已加载的执行依据列表，使用以下格式：

   ```md
   ### 🧩 Loaded SKILLs
   - `speckit-implement`
   - `<matched-skill-1>`
   - `<matched-skill-2>`
   ```

   若当前 Phase 没有额外执行依据，也必须输出，仅保留 `speckit-implement`。
6. 完成当前 Phase 的代码实现后，必须依据当前 Phase 已加载的内容完成确认：
   - 若存在明确确认要求，则按要求逐项确认
   - 若仅存在原则性或约束性内容，则完成最小必要确认
   - 若当前 Phase 缺少足够确认依据，则必须显式说明当前确认边界与剩余风险
7. 完成当前 Phase 中的任务。
8. 仅当以下条件同时满足时，才允许将当前批次中已真实落地的任务更新为 `[✅️]`（`tasks.md` 对应任务项）：
   - 代码已实际落地
   - 与改动范围匹配的验证已完成
   - 当前 Phase 所需确认已完成
   - 当前不存在阻断该 Phase 关闭的未解决问题
9. 若确认结果表明当前实现仍不满足当前 Phase 的执行依据，则不得回写完成状态，必须继续修复或向用户报告阻塞原因。
10. 仅当当前 Phase 已完成后，才允许进入下一个 Phase，并重复步骤 2-9，直到所有任务完成。

### 3. 阻断处理

若出现以下情况，停止推进并向用户报告：

- 缺少必要上下文，无法安全实施
- 存在严重设计冲突，无法自行裁决
- 当前任务明显需要补充任务清单或重排边界
- 关键依赖、脚本或必需 skill 无法正常工作

## 输出

- 更新后的 `specs/<feature>/tasks.md`
- 每个 Phase 开始执行时输出当前阶段加载的执行依据列表，标题固定为 `### 🧩 Loaded SKILLs`
- 若 `tasks.md` 没有 Phase 级状态字段，则默认以该 Phase 下任务是否全部回写为阶段完成依据，不额外修改 Phase 标题样式
- 若当前 Phase 的确认依据不足，需显式说明当前确认边界与剩余风险，不得将其伪装为完整确认
- 如全部完成且无额外问题待记录，提示进入 `/code-review` 进行代码审查，或通过 `/issue` 提出本次执行的相关问题

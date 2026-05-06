---
name: speckit-checklist
description: 审计 spec.md 的完整性、边界和可测性，并生成需求检查单。用于进入 /plan 前的质量门检查。
---

# Spec Kit Checklist Skill

## 何时使用

- `spec.md` 已存在，并且希望在进入 `/plan` 前做一次需求质量门检查
- spec 发生了较大变动，需要重新审视完整性、边界和可测性
- 检查的是“需求文档质量”，不是代码实现正确性

## 职责边界

- `speckit-checklist` 负责产出**持久化需求检查单**
- 它可以把发现的严重缺口追加为新的 CQ，但**不负责**替用户完成澄清
- 发现新 CQ 后，后续闭环应回到 `/clarify`

## 输入

- 必需：`specs/<feature>/spec.md`
- 可选：用户指定的检查焦点（如安全、边界、非功能、接口设计）
- 可选：`plan.md` / `tasks.md`，仅用于做一致性参考，不替代 spec 本体

## 工作流

### 1. 前置检查

1. 确认 `specs/<feature>/spec.md` 存在
2. 若 spec 不存在，终止并提示先运行 `/specify`
3. 若 spec 仍有大量未解决 CQ，可继续做审计，但必须在结果中标注“当前文档仍处于 Clarifying”

### 2. 主题与范围确定

根据用户输入确定本次 checklist 的重点：

- 功能完整性
- 数据模型
- 接口设计
- 边界与失败处理
- 安全与合规
- 非功能属性
- 术语一致性与可测性

若用户未指定，则做通用质量审计。

### 3. 文档读取

按需提取：

- `spec.md`：RQ、关键实体、边界、约束、已有 CQ
- `plan.md`：仅用于检查设计假设是否和 spec 漂移
- `tasks.md`：仅用于检查任务覆盖是否提前暴露需求空洞

### 4. 检查单生成

输出持久化文件 `specs/<feature>/checklists/requirements.md`，至少包含：

- 检查范围与焦点
- 检查维度
- 通过项
- 风险项 / 缺口项
- 建议新增 CQ
- 下一步建议

检查项必须聚焦需求质量，不写实现验证语句。

### 5. CQ 反推规则

如果发现会影响后续设计、任务拆解或验收定义的重大缺口：

1. 将缺口转化为新的 CQ
2. 追加到 `spec.md` 的 `## 待澄清问题`
3. 避免生成与现有 CQ 重复的问题
4. 若追加了新的未解 CQ，则将 `spec.md` 状态置为 `Clarifying`

如果没有新增 CQ：

- 保持原有状态
- 若所有 CQ 都已解决，则可视为 `Ready`

## 输出

- `specs/<feature>/checklists/requirements.md`
- 如发现新缺口：更新后的 `specs/<feature>/spec.md`
- 若新增了 CQ，提示用户先运行 `/clarify`
- 若未新增 CQ 且 spec 为 `Ready`，提示用户进入 `/plan`

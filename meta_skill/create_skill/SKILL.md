---
name: create-skill
description: 用于新建或重构SKILL技能（SKILL.md +reference/ 文件体系）。适用于从零创建技能、改写现有技能结构、诊断技能为何未被正确触发的场景。
 触发关键词：创建技能、新建 skill、写 SKILL.md、技能不触发、重构技能、create skill。
 DO NOT USE FOR：通用编码问题（用默认 agent）、MCP 服务器配置、VS Code 插件开发。
---

# Create Skill

## 输入

- 必需：技能的领域 / 主题
- 必需：适用场景描述（用户什么时候会用到它）
- 必需：核心约束（MUST / MUST NOT 级别的规则）
- 可选：已有的 reference 文件或示例代码
- 可选：不适用场景（DO NOT USE FOR）

缺少必需输入时，MUST 先向用户提问收集，不得猜测继续。

## 工作流

1. 收集输入（领域、场景、约束、示例），缺少时提问补全
2. 撰写 `description`（路由层，见约束）
3. 撰写 `SKILL.md` 正文（按骨架顺序，见参考文件）
4. 按需创建 `reference/` 文件（决策表 + 正反例，见参考文件）
5. 输出完成后结构化状态

## 约束

### MUST
- `description` MUST 内联 USE FOR、触发关键词、DO NOT USE FOR，不得只写功能描述
- `SKILL.md` MUST 包含：输入、工作流、约束、强制输出、参考文件 五个章节
- 约束 MUST 按 MUST / MUST NOT / SHOULD 三级分类，内联在 `SKILL.md` 中
- `reference/` 文件 MUST 包含决策表或正反例，不得只写 prose 说明
- 参考文件清单 MUST 以三列表格呈现（文件 / 适用场景 / 加载时机）

### MUST NOT
- MUST NOT 在 `SKILL.md` 正文中写 `## 何时使用` 章节（该内容属于 description）
- MUST NOT 创建 `rules/` 文件夹（规则内联到 `## 约束`）
- MUST NOT 在 `reference/` 中写纯 prose 规范文字，无示例代码支撑
- MUST NOT 一个 skill 承担多个不相关职责

### SHOULD
- `reference/` 文件 SHOULD 包含边界案例（⚠️ 灰色地带判定）
- `reference/` 文件 SHOULD 包含组合场景（多规则同时生效的完整示例）
- 强制输出 SHOULD 使用 JSON 结构，而不是自然语言列表

## 强制输出

开始前输出：

```json
{
  "skillName": "待创建技能名称",
  "domain": "领域描述",
  "referenceFiles": ["计划创建的 reference 文件列表"]
}
```

完成后输出：

```json
{
  "descriptionHasTriggers": true,
  "skillMdHasAllSections": true,
  "noRulesFolder": true,
  "referenceHasExamples": true
}
```

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `reference/skill-spec.md` | 撰写 SKILL.md 骨架、description、约束分级、强制输出 | 创建任何技能前必读 |
| `reference/reference-spec.md` | 撰写 reference/ 文件，包含决策表、正反例、边界案例、组合场景 | 需要创建 reference 文件时 |

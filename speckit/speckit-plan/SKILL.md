---
name: speckit-plan
description: 基于用户目标和已有 spec 生成可执行的 plan.md，包括阶段划分、技术方案、设计边界、风险识别和交付规划。
  触发关键词：plan、计划、规划、roadmap、技术方案、阶段划分、出计划。
  DO NOT USE FOR：需求分析与澄清（用 speckit-clarify / speckit-specify）、任务拆解（用 speckit-tasks）、代码实施（用 speckit-implement）。
---

# Spec Kit Plan Skill

## 输入

- 必需：用户输入，包括目标、范围、约束、预期交付物
- 可选：`specs/<feature>/spec.md`
- 模板：`./templates/plan-template.md`

缺少必需输入时，MUST 先向用户提问收集目标与交付物，不得跳过直接生成计划。

## 工作流

### 1. 前置检查

1. 读取 `./templates/plan-template.md`
2. 若存在 `specs/<feature>/spec.md`：提取需求、关键实体、约束和边界
3. 若 spec 不存在或信息不足：继续生成草案，但显式列出待澄清项

### 2. 项目类型识别

结合仓库事实识别当前项目类型，作为后续设计视角的选择依据：

| 项目类型 | 核心设计视角 |
|---------|------------|
| 服务端项目 | 核心实体、接口边界、数据访问、测试与验证 |
| 库 / SDK 项目 | 公共接口、兼容性、依赖边界、示例与测试 |
| CLI 工具 | 命令结构、参数输入、输出格式、错误处理 |
| 前端应用 | 页面结构、状态管理、接口交互、组件边界 |
| 网关 / 接入层 | 路由、代理、协议适配、依赖边界 |
| 数据处理项目 | 输入输出格式、批处理流程、容错策略、监控 |

若已识别到项目适配 skill，将其 MUST 约束作为规划输入；若无法识别，退化为通用规划流程并显式记录假设。

### 3. Phase 0 领域调研

1. 分析关键实体及关系
2. 分析状态流转与异常路径
3. 识别外部依赖与上下游服务
4. 每个调研项输出：决策、理由、排除方案

### 4. Phase 1 技术设计

按识别到的项目类型选择对应设计视角，逐项输出设计要点。

若已加载项目适配 skill，在设计要点中补充：
- 框架或平台特有约束
- 关键设计决策的规范来源
- 与项目现有规范可能冲突的风险点

### 5. 写入 plan.md

- 输出路径：`specs/<feature>/plan.md`
- 保持模板章节结构
- 若 spec 存在，Phase 0 / Phase 1 内容应可追溯到需求点
- 若 spec 不存在，显式标注临时需求点和待澄清项

### 6. 质量校验

| # | 校验项 | 标准 |
|---|--------|------|
| 1 | 需求覆盖 | 每条需求或需求点至少映射到一个设计项 |
| 2 | 实体完整 | 关键实体与关系已分析 |
| 3 | 交付明确 | 每个阶段有明确产出物 |
| 4 | 依赖识别 | 关键依赖与约束已显式说明 |
| 5 | 风险可验证 | 关键风险有验证或缓解策略 |
| 6 | 假设可见 | 若项目类型识别不确定，假设与待澄清项已显式记录 |

## 输出

- `specs/<feature>/plan.md`

## 约束

### MUST
- 生成 plan 前 MUST 读取 `./templates/plan-template.md`，输出结构必须与模板章节对齐
- 若 spec 存在，设计项 MUST 可追溯到 spec 中的需求点
- 若项目类型识别不确定，MUST 显式记录假设与待澄清项，不得假装已确认

### MUST NOT
- MUST NOT 在 spec 和用户输入均不足时直接生成计划，必须先列出待澄清项
- MUST NOT 在 plan.md 中直接写实现代码（plan 只写设计与边界，不写代码）
- MUST NOT 跳过质量校验直接输出 plan.md

### SHOULD
- SHOULD 在调研阶段对每个决策输出排除方案，使设计选择可追溯
- SHOULD 在 plan 中标注哪些设计项依赖项目适配 skill 的约束

## 强制输出

生成 plan 前输出：

```json
{
  "specExists": true,
  "projectType": "识别到的项目类型或 unknown",
  "adaptSkillLoaded": "skill 名称或 none",
  "clarifyItems": ["待澄清项列表，若无则为空数组"]
}
```

完成后输出：

```json
{
  "outputFile": "specs/<feature>/plan.md",
  "requirementsCovered": true,
  "entitiesAnalyzed": true,
  "risksIdentified": true
}
```

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `./templates/plan-template.md` | plan.md 章节结构与格式 | 生成任何 plan 前必读 |
| `specs/<feature>/spec.md` | 需求、关键实体、约束与边界 | 前置检查时读取（若存在） |

---
name: speckit-doc
description: 基于仓库事实初始化或更新 `docs/architecture.md`。用于沉淀标准化架构文档，支持 `init` 基于工程事实建档，支持 `update` 基于 `spec.md`、`plan.md`、`tasks.md` 增量更新架构说明。
---

# Speckit Doc

## 何时使用

- 需要为当前仓库首次建立标准化架构文档时，使用 `init`
- 已存在 `docs/architecture.md`，需要结合当前 feature 的 `spec.md`、`plan.md`、`tasks.md` 更新架构文档时，使用 `update`
- 用户明确要求生成、维护、补全系统架构文档、核心流程图、实体说明时

## 职责边界

- `speckit-doc` 只负责生成和维护 `docs/architecture.md`
- `speckit-doc` 只基于仓库事实、现有文档、以及 speckit 产物进行归纳，不捏造不存在的模块、流程、实体或依赖
- 若信息来自代码与配置，按“已实现事实”表述
- 若信息来自 `spec.md`、`plan.md`、`tasks.md` 但尚未完全落地，按“目标设计”或“规划中”表述
- 文档结构必须遵循标准模板，不允许随意增删核心章节
- 若仓库事实不足以支撑某节内容，明确标注“待补充”或“不适用”，不得虚构

## 输入

### `init`

- 仓库根目录下的代码、配置、README、现有系统文档
- 可按需读取项目适配技能或相关规范，抽取工程事实

### `update`

- 已存在的 `docs/architecture.md`
- 当前 feature 的 `specs/<feature>/spec.md`
- 当前 feature 的 `specs/<feature>/plan.md`
- 当前 feature 的 `specs/<feature>/tasks.md`
- 与该 feature 相关的代码、配置和文档变更

## 输出

- 标准模板：本技能目录下 `./templates/architecture-template.md`
- 标准规范：`./references/doc-standard.md`
- 工作流说明：`./references/update-workflow.md`
- 目标文档：`docs/architecture.md`

## 固定文档结构

`docs/architecture.md` 必须包含以下部分：

1. 文档概览
2. 总体架构说明
3. 总架构流程图
4. 核心功能与核心流程图
5. 实体设计
6. 关键约束与边界

其中：

- “总架构流程图”必须使用 Mermaid，描述系统级模块关系、数据流或调用关系
- “核心功能与核心流程图”至少覆盖当前系统最核心的 1 个主流程；若存在多个关键能力，可拆为多个子节
- “实体设计”必须抽象核心业务实体，描述职责、关键属性和关系

## 工作流

### 1. 模式判定

- 若用户要求首次建档、初始化架构文档、补齐仓库架构全貌，执行 `init`
- 若用户要求根据某个 feature 的 speckit 产物更新架构文档，执行 `update`

### 2. 事实抽取

优先从以下信息源抽取事实：

1. 代码目录结构
2. README、已有架构文档、接口说明
3. 关键配置文件
4. `spec.md`、`plan.md`、`tasks.md`
5. 与 feature 直接相关的实现文件

抽取时至少识别：

- 系统边界与主要模块
- 模块间调用关系或数据流
- 核心业务流程
- 核心业务实体
- 外部依赖或基础设施

### 3. 标准化归档

1. 读取本技能目录下 `./references/doc-standard.md`
2. 读取本技能目录下 `./templates/architecture-template.md`
3. 按模板填充 `docs/architecture.md`
4. 所有章节使用统一命名、统一术语、统一图示风格
5. 最终文档只保留项目架构结果，不保留 speckit 过程痕迹

### 4. `init` 规则

- 若 `docs/architecture.md` 不存在，则创建
- 若已存在，但用户明确要求重新初始化，则基于当前仓库事实整体重建
- 必须区分“仓库已实现事实”与“推断/待确认内容”

### 5. `update` 规则

- 先读取已有 `docs/architecture.md`
- 再读取 `spec.md`、`plan.md`、`tasks.md`，识别新增或变更的模块、流程、实体、依赖、边界
- 仅更新受影响章节，保持未受影响内容稳定
- 若 `tasks.md` 显示功能尚未完成，则在文档中以“规划中”或“待实现”标注，不写成既成事实
- 更新后确保术语、图示、章节编号和实体命名保持一致
- 更新结果直接并入对应章节，不追加变更流水或过程记录

### 6. 质量要求

- 架构图和流程图必须可读、可维护，不堆砌过多节点
- 实体定义必须服务于架构理解，而不是照搬数据库字段
- 文档必须面向团队协作，突出模块职责、交互边界、关键流程、实体关系
- 任何新增内容都必须能追溯到仓库事实或 speckit 产物

## 执行提示

- 如果仓库内已存在其他架构文档，可将其作为事实来源之一，但不得把某个具体项目的章节设计或业务术语固化进通用技能
- 如果项目存在明显的工程适配技能，应结合适配技能理解项目分层、模块职责与交付边界
- 若实体较多，仅保留对系统理解最关键的实体，不做 ER 全量穷举

---
name: workflow-skill
description: 用于 OpenClaw 中枢式多 agent 工作流编排。适用于一个主 agent 接收任务后，按任务类型把工作流转给产品经理、UI 设计、架构、开发、代码审查、项目管理等子 agent，并在主 agent 汇总结果的场景。触发关键词包括多 agent 流转、工作流编排、sessions_spawn、跨角色协作、需求分析到实现、项目跟踪。
---

# Workflow Skill

用于 OpenClaw 的“主 agent 做中枢，子 agent 做叶子任务”模式。

## 何时使用

当任务需要以下任一能力时使用本技能：

- 一个主 agent 根据任务类型选择不同角色 agent
- 一个需求需要按阶段流转，例如 `pm -> architect -> developer -> reviewer`
- 一个复杂任务需要拆成多个互相独立的子任务并发处理
- 一个项目经理 agent 需要接收阶段性进展，但不直接做实现

如果只是单个 agent 自己完成任务，不要使用本技能。

## 核心原则

1. 所有外部入口先进入主 agent。
2. 主 agent 负责分类、拆解、转交、汇总。
3. 子 agent 只处理单一明确子任务。
4. 子 agent 的结果先回主 agent，再由主 agent 决定是否继续流转。
5. 除非有明确需要，否则不要让子 agent 彼此直接通信。

## 角色模板

本技能内置 6 个推荐角色模板：

- `references/pm.md`
- `references/ui.md`
- `references/architect.md`
- `references/developer.md`
- `references/reviewer.md`
- `references/project-manager.md`

这些模板是工作流侧的“推荐角色定义”，用于：

- 帮主 agent选择是否需要某个角色
- 帮主 agent构造传给 `create-sub-agent` 的模板参数
- 帮主 agent生成更清晰的 handoff task

每个模板都包含：

- 建议的 agent id
- 职责边界
- 不负责事项
- 推荐输入
- 推荐输出
- `IDENTITY.md` 模板
- `openclaw.json` 配置片段

## 缺角色时怎么做

如果发现所需角色不存在：

1. 优先调用 `create-sub-agent`
2. 把当前选中的角色模板作为输入参数传给 `create-sub-agent`
3. 只创建当前缺失角色，不要无关扩装
4. 如果当前环境没有 `create-sub-agent`，明确建议用户先安装该技能

如果没有现成模板：

1. 调用 `create-sub-agent`
2. 明确告诉它当前没有模板
3. 让它向用户补问最小必要信息并生成自定义模板

不要在本技能内直接承担“创建 agent 配置和 workspace”的职责。

## 调度规则

主 agent 收到任务后按下列规则决策：

- 纯需求定义问题，转给 `pm`
- 纯界面或交互问题，转给 `ui`
- 纯架构设计问题，转给 `architect`
- 纯实现问题，转给 `developer`
- 纯代码审查问题，转给 `reviewer`
- 纯计划、排期、推进问题，转给 `project-manager`

如果任务跨多个阶段，使用串行流转：

1. 先把需求交给 `pm`
2. 再把 `pm` 的结论交给 `architect` 或 `developer`
3. 再把阶段结果交给 `reviewer`
4. 再把里程碑和阻塞项同步给 `project-manager`

如果任务可拆成多个独立子问题，使用并行流转：

- 每个子任务只派给一个最匹配的 agent
- 每个子任务都写清楚输入、输出、约束
- 等全部完成后，由主 agent 汇总

## 工具使用

使用 `sessions_spawn` 创建子 agent 任务。

每次转交都要：

- 显式指定 `agentId`
- 把上一步结论写进新的 `task`
- 指定清晰产出

任务描述建议包含：

- 本阶段目标
- 上游输入摘要
- 需要产出的内容
- 禁止事项

常见串行流转提示参考 `references/handoff-template.md`。

## 与 create-sub-agent 的协作方式

当缺少角色时，主 agent 应把模板摘要作为参数传给 `create-sub-agent`，例如：

```text
调用 create-sub-agent。
模板参数：
- id: architect
- name: Architect
- responsibilities: 技术方案、模块边界、数据流、风险与取舍
- not_responsible_for: 不负责需求定义、不负责视觉设计、不直接负责排期
- input: 需求摘要、验收标准、仓库约束
- output: 技术方案、模块划分、风险、最小实现路径
```

如果没有模板，则传递“无模板，需要补问用户”的明确说明。

## 禁止事项

- 不要先读一堆无关 skill 再决定怎么调度
- 不要把一个模糊的大任务直接丢给子 agent
- 不要让多个子 agent 同时修改同一块代码
- 不要轮询 `sessions_list`、`sessions_history` 或用 `sleep` 等待
- 不要让 workflow skill 自己直接实现 create-sub-agent 的安装逻辑

OpenClaw 的 sub-agent 完成通知是 push-based，等待完成事件即可。

## 输出要求

主 agent 的最终输出应包含：

- 当前任务经过了哪些 agent
- 每个阶段的关键结论
- 当前状态是完成、待实现、待审查还是阻塞
- 如果仍有后续动作，给出下一步建议

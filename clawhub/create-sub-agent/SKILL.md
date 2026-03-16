---
name: create-sub-agent
description: 用于在 OpenClaw 中新增一个子 agent，或按传入模板批量安装一组 sub-agents，并指导如何描述其职责、如何配置主 agent 调度、以及需要修改哪些文件。适用于新增产品、设计、架构、开发、测试、审查等角色 agent 的场景。触发关键词包括新增子 agent、创建 agent、角色 agent、agent 配置、allowAgents、workspace、IDENTITY.md、模板创建 agent。
---

# Create Sub Agent

用于在 OpenClaw 中新增一个可被主 agent 调度的角色 agent。

## 何时使用

当你需要新增一个具有明确边界的角色 agent 时使用，例如：

- 新增 `developer`
- 新增 `qa`
- 新增 `data-analyst`
- 新增 `ops`

如果只是给现有 agent 增加一条行为规则，不要创建新 agent。

当其他 skill 已经准备好了角色模板，并希望你根据模板完成创建时，也使用本技能。

## 单一职责

本技能只负责两件事：

1. 根据传入模板创建 sub-agent
2. 如果没有模板，则补问用户并生成一个自定义模板后再创建

本技能不负责判断工作流应该使用哪个角色。  
角色选择和工作流编排不属于本技能的职责范围。

## 输入约定

优先接受“调用方显式传入的模板参数”。

推荐模板字段：

- `id`
- `name`
- `responsibilities`
- `not_responsible_for`
- `input`
- `output`
- `workspace`（可选）

如果这些字段已提供，直接基于这些字段创建 agent。  
不要再隐式猜测模板来源，不要假设必须从别的 skill 目录读取模板。

## 没有模板时怎么办

如果调用方没有提供模板：

1. 明确说明当前没有模板
2. 向用户补问最小必要信息
3. 生成一个自定义模板
4. 再基于该模板创建 agent

最小必要信息至少包括：

- agent id
- agent name
- 职责
- 不负责什么
- 输入
- 输出

## 创建标准

只有在以下条件同时满足时才新增子 agent：

1. 该角色有稳定、可复用的职责边界
2. 主 agent 能明确判断何时把任务交给它
3. 该角色的输出和其他角色不同
4. 该角色单独存在能减少主 agent 的上下文负担

## 新 agent 的描述方式

定义一个子 agent 时，必须描述清楚：

- `id`: 机器使用的唯一标识，使用小写字母、数字、连字符
- `name`: 人类可读名称
- `职责`: 这个 agent 负责什么
- `不负责`: 这个 agent 不应该做什么
- `输入`: 主 agent 传给它什么
- `输出`: 它完成后必须返回什么
- `workspace`: 它在哪个工作目录运行

推荐描述模板：

```text
agent id: developer
职责: 根据既有需求和技术方案完成代码实现、测试补齐和本地验证
不负责: 不负责重新定义需求，不负责做最终代码审查，不负责排期
输入: 来自主 agent 的需求摘要、技术方案、目标文件或模块
输出: 修改摘要、关键文件、验证结果、剩余风险
```

## 必改文件

新增一个子 agent 时，通常要改这些地方：

1. `openclaw.json`
2. 新 agent 对应的 workspace 目录
3. 可选的 `IDENTITY.md`

详细检查单见 `references/config-checklist.md`。

## 批量安装一组角色

如果调用方已经决定要创建哪些角色，可以通过脚本批量创建，但角色必须显式传入：

```bash
node {baseDir}/scripts/install-subagents.mjs --main-agent main --workspace-root . --roles pm,architect,reviewer
```

如果希望同时传角色 id 和显示名，使用 `id:name`：

```bash
node {baseDir}/scripts/install-subagents.mjs --main-agent main --workspace-root . --roles pm:Product Manager,architect:Architect,reviewer:Code Reviewer
```

脚本行为：

- 先对每个目标角色执行 `openclaw agents add <id> --workspace <dir> --non-interactive`
- 确保主 agent 存在并允许调度这些 agent
- 为每个目标 workspace 补基础 `IDENTITY.md`
- 如果主 agent workspace 下存在 `skills/`，则复制一份到每个 sub-agent workspace，已有同名技能默认保留

脚本直接调用 `openclaw` 命令。运行前确保当前环境可以直接执行 `openclaw`。

角色必须由调用方显式传入。  
不允许依赖脚本内置默认角色集合。

如果未显式传 `--config`，默认读取 `~/.openclaw/openclaw.json`。

## openclaw.json 需要改什么

至少修改：

- `agents.list` 中新增一个 entry
- 给它配置 `workspace`
- 如果主 agent 要能调度它，把它加入 `main.subagents.allowAgents`

如果新 agent 也要作为外部入口目标，再额外考虑：

- `bindings`
- 是否设为默认 agent

最小示例：

```json
{
  "agents": {
    "list": [
      {
        "id": "main",
        "default": true,
        "workspace": "./workspace-main",
        "subagents": {
          "allowAgents": ["developer"]
        }
      },
      {
        "id": "developer",
        "workspace": "./workspace-developer"
      }
    ]
  }
}
```

## workspace 需要改什么

为新 agent 创建独立 workspace，例如：

- `workspace-developer/`
- `workspace-qa/`

如果该角色需要专门身份提示，可在 workspace 下新增：

- `IDENTITY.md`

`IDENTITY.md` 里只写这个角色的职责和边界，不要写全局编排规则。

## 主 agent 侧需要同步什么

新增子 agent 后，主 agent 的调度规则也要同步更新：

- 何时路由到这个新 agent
- 需要给它什么输入
- 它返回后是否还要继续流转到别的 agent

如果主 agent 规则没更新，新 agent 即使存在，也不会被正确使用。

## 创建步骤

1. 获取调用方传入的模板，或向用户补问生成模板
2. 先使用 `openclaw agents add` 创建 agent
3. 给主 agent 的 `subagents.allowAgents` 加上新 agent id
4. 补写对应 workspace 的 `IDENTITY.md`
5. 如果主 agent workspace 有 `skills/`，复制到子 agent workspace
6. 验证主 agent 能通过 `sessions_spawn(agentId=...)` 转交给它

## 不要这样做

- 不要只创建一个文件夹就认为 agent 已经存在
- 不要新增一个职责高度重叠的 agent
- 不要忘记把新 agent 加入主 agent 的 `allowAgents`
- 不要把“如何创建 agent”的说明写进 `MEMORY.md`
- 不要自己决定应该创建哪个角色，除非调用方明确授权你这样做

## 输出要求

当用户要求“新增一个子 agent”时，输出至少应包含：

- 使用的模板内容，或补问后生成的模板摘要
- 建议的 `agent id`
- 建议职责描述
- 需要修改的 `openclaw.json` 片段
- 需要创建的 workspace 路径
- 是否建议新增 `IDENTITY.md`
- 是否需要把主 agent 的调度规则同步更新

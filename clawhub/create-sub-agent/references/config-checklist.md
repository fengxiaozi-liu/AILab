# Config Checklist

新增子 agent 时逐项检查：

## openclaw.json

- `agents.list[]` 新增角色 entry
- `id` 唯一且规范
- `workspace` 已设置
- 如需专属模型，设置 `model`
- 如需专属身份提示，设置 `identity` 或使用 workspace 下的 `IDENTITY.md`
- 在主 agent 的 `subagents.allowAgents` 中加入该 id

## Workspace

- 创建对应目录，例如 `workspace-developer/`
- 如需单独角色提示，新增 `workspace-developer/IDENTITY.md`

## Routing

- 如果该 agent 要接外部入口，更新 `bindings`
- 如果只作为 worker，被主 agent 调度即可，不一定需要 `bindings`

## Validation

- `openclaw agents list` 能看到新 agent
- 主 agent 的配置中已允许调度该 agent
- 主 agent 的提示词或身份规则已包含该角色的路由条件

## Bulk Install

- 如果要批量安装一整套角色，使用 `scripts/install-subagents.ts`
- 批量安装后确认每个 `workspace-*` 下都已有 `IDENTITY.md`
- 如果主 agent workspace 下存在 `skills/`，确认子 agent workspace 下也已复制对应技能

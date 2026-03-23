---
name: kratos-patterns
description: 用于 Kratos 项目路由判定与子技能选择，是所有 Kratos 代码修改的入口技能。
  适用于任何 Kratos 代码修改、重构、代码评审、生成物更新，或不确定应加载哪个 skill 的场景。
  触发关键词：SERVER_NAME、项目类型判定、BaseService、业务服务、网关、路由、技能选择、实现前审计、kratos。
  DO NOT USE FOR：非 Kratos 项目（→ golang-patterns）、已明确子技能时可直接加载对应 skill。
---

# Kratos Patterns

## 输入

- 必需：本次变更目标或问题描述
- 可选：`.env.*` 中的 `SERVER_NAME`、涉及文件路径或目录
- 可选：`specs/<feature>/tasks.md`

缺少必需输入时，MUST 先向用户提问，不得猜测继续。

## 工作流

1. 判定工作域（`WorkDomain`）：读取 `.env.*` 中的 `SERVER_NAME`
   - `SERVER_NAME` 以 `BaseService` 结尾 → `BaseService`
   - `SERVER_NAME` 包含 `GatewayService` 或 `OpenapiService` → 网关
   - 其他 → 业务项目
   - IF 无法读取 `SERVER_NAME` → 依次按目录特征（gateway/proxy 目录）、仓库内容（proto/枚举/错误码）兜底判定
2. 输出开始前结构化状态（见强制输出）
3. 按本次改动职责变化选择最小聚合 skill 集合：
   - 聚合边界 / layer / UseCase / Repo / 测试 → `kratos-domain`
   - Ent / EventBus / Crontab / Depend / Config → `kratos-components`
   - context / metadata / middleware / proto helper / schema / seata / util → `kratos-pkg`
   - Error / Enum / i18n / Logging / Comment → `kratos-conventions`
   - Server / Gateway / Codegen / Proto / Wire / Build → `kratos-entry`
4. 加载选定的聚合 skill，按各 skill 工作流执行变更
5. 执行所有 codegen / build / test 验证
6. 输出完成后结构化状态（见强制输出）

## 约束

### MUST
- 任何 Kratos 代码修改前 MUST 先完成路由判定，再选择聚合 skill
- 项目类型判定 MUST 优先看职责与语义，路径只作证据提示
- 修改前 MUST 先验证仓库事实，不凭猜测实施
- 变更 proto / wire / ent schema 时 MUST 进入生成与校验链路
- 仅实现与需求直接相关的变更，Contract First：跨边界能力先定协议再落实现

### MUST NOT
- MUST NOT 跳过 `kratos-patterns` 直接选择子 skill（已明确职责范围时除外）
- MUST NOT 仅凭文件路径决定业务边界
- MUST NOT 未确认需求边界时扩大改动范围
- MUST NOT 修改生成物代替修改源文件
- MUST NOT 只给结论不附带路由证据

### SHOULD
- 一个任务可同时加载多个聚合 skill，但 SHOULD 保持最小集合

## 强制输出

开始前输出：

```json
{
  "workDomain": "BaseService | 业务 | 网关",
  "evidence": "用于判定的证据（SERVER_NAME / 目录特征 / 仓库内容）",
  "subSkills": ["kratos-domain", "kratos-entry"],
  "codegenPlan": ["make api", "wire gen", "ent generate"]
}
```

完成后输出：

```json
{
  "subSkillChecksAllPassed": true,
  "codegenDone": true,
  "buildPassed": true
}
```

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `kratos-domain/reference/*` | 聚合、分层、命名、Repo、UseCase、测试 | 加载 kratos-domain 时进入对应目录 |
| `kratos-components/reference/*` | Ent、EventBus、Crontab、Depend、Config | 加载 kratos-components 时进入对应目录 |
| `kratos-pkg/reference/*` | context、metadata、middleware、proto helper、schema、seata、util | 加载 kratos-pkg 时进入对应目录 |
| `kratos-conventions/reference/*` | error、enum、i18n、logging、comment | 加载 kratos-conventions 时进入对应目录 |
| `kratos-entry/reference/*` | proto、server、gateway、wire、codegen | 加载 kratos-entry 时进入对应目录 |

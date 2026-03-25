---
name: kratos-routing
description: Kratos 项目类型与代码落层判定。仅用于判断仓库是 BaseService、GatewayService 还是业务服务，或判断代码应落在 `server`、`service`、`biz/usecase`、`data/repo`、`listener/consumer/cron` 哪一层。不要把它当作总入口、主技能选择器或实现技能。
---

# Kratos Routing

## 何时使用

- 需要判断当前仓库属于 `BaseService`、`GatewayService` 还是业务服务。
- 需要判断某段改动应落在 `server`、`service`、`biz/usecase`、`data/repo`、`listener/consumer/cron` 哪一层。
- 用户给出的需求较模糊，先要确定代码坐标、层级归属和目录落点。

## 职责边界

- 本技能只做项目类型判定与层级判定。
- 本技能不负责主技能选择，不负责多个技能编排，也不负责具体代码实现。

## 输入

- 必需：当前任务目标、待修改文件范围，或待判断落层的代码职责描述。
- 可选：`SERVER_NAME`、`.env.*`、配置文件、目录结构、现有文件路径等仓库证据。
- 可选：候选落点，例如 `internal/service`、`internal/biz`、`internal/data`、`internal/server`。

缺少必需输入时，MUST 先从工作区和任务上下文补齐；仍无法判断项目类型或层级时，再向用户提问，不得直接猜测。

## 工作流

### 收集证据与补齐输入

- 优先使用：`SERVER_NAME`、`.env.*`、配置文件、目录结构、待修改文件路径与职责描述。
- 输入不足时：先在仓库中检索项目类型证据与目录落点证据；仍不足再向用户追问（不要猜）。

### 判定项目类型（只选一个主类型）

- `BaseService`：`SERVER_NAME` 以 `BaseService` 结尾或存在明确约定证据。
- `GatewayService`：`SERVER_NAME` 包含 `GatewayService` 或存在明确约定证据。
- 业务服务：`SERVER_NAME` 包含具体业务名或存在明确约定证据。=

### 按需加载 references

- 读 `references/layer-spec.md`，用目录结构与职责边界校验代码落层判断。

### 输出判定与建议（只做判定职责）

- 给出项目类型、证据与建议落层（`server` / `service` / `biz/usecase` / `data/repo` / `listener/consumer/cron`）。
- 若证据不足，明确说明假设来源与缺失信息点。

### 边界自检（不做测试验收）

- 是否给出了可核验的仓库证据（而不是抽象结论）。
- 是否避免承担主技能选择与具体实现责任。

## 约束

### MUST

- MUST 只承载项目类型判定和代码落层判定。
- MUST 优先使用仓库事实作为证据，例如 `SERVER_NAME`、目录结构、模块命名、实际文件路径。
- MUST 在输出结果时给出明确证据，而不是只给抽象结论。
- MUST 在证据不足时说明不确定性和假设来源。

### MUST NOT

- MUST NOT 把自己表述成任何 Kratos 任务的总入口或前置总控。
- MUST NOT 负责主技能选择、辅助技能组合或后续实现兜底。
- MUST NOT 在证据不足时伪造项目类型或强行解释层级归属。

### SHOULD

- SHOULD 先判定项目类型，再判定层级归属。
- SHOULD 优先使用实际文件路径和改动面判断层级，而不是只依赖自然语言描述。
- SHOULD 在判定完成后给出最小、明确的目录落点建议。

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `references/layer-spec.md` | 判断代码应落在 `server/service/biz/usecase/data/repo` 哪一层，以及 `{Business}Service` 的常见层级结构 | 需要做层级判定时加载 |

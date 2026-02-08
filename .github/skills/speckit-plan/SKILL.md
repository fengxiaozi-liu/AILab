---
name: speckit-plan
description: 基于 spec 生成 Kratos 微服务实施计划。
---

# Spec Kit Plan Skill（中文）

## 何时使用

- spec 已就绪（Status 为 Ready），需要生成技术实施计划。

## 输入

- `specs/<feature>/spec.md`（由 specify + clarify 产出）

## 依赖 Skill

- **kratos-patterns**：plan agent 同时加载本 skill 与 kratos-patterns skill，后者提供框架适配知识（Ent Schema 规范、Proto 规范、分层规范等），用于将领域设计扩充为 Kratos 项目可执行计划。

## 工作流

### §前置检查

1. 读取 `specs/<feature>/spec.md`，确认 Status 为 Ready。
2. 若 Status 不为 Ready 或 spec 不存在，终止并提示先运行 `/specify` 或 `/clarify`。
3. 读取 `.specify/templates/plan-template.md` 获取计划模板。

### §项目类型判定

通过 `.env.*` 中的 `SERVER_NAME` 判定当前项目类型：

| 标识 | 项目类型 |
|------|---------|
| 以 `BaseService` 结尾 | BaseService 项目（抽象定义） |
| 含 `GatewayService` 或 `OpenapiService` | 网关类项目（接口实现） |
| 其他 | 业务项目（业务实现） |

判定结果写入 plan.md，后续流程根据项目类型选择对应设计内容。

### §Phase 0 — 领域调研

从 spec 的**核心需求（RQ）**和**关键实体**推导调研内容：

1. **实体与关系分析**：
   - 从 RQ 和关键实体中提取所有业务实体
   - 分析实体间关系（1:N、M:N、1:1）
   - 识别关键属性和约束

2. **状态流转**（若 spec 涉及状态变化）：
   - 定义状态枚举
   - 明确流转触发条件和规则

3. **依赖服务调研**：
   - 从 RQ 的约束与边界中识别需要调用的上下游服务
   - 检查 `internal/biz/depend/` 和 `internal/data/depend/` 确认 InnerRPC 是否已有
   - 标记需要新增或修改的依赖

4. **每个调研项输出**：
   - 决策（选择了什么）
   - 理由（为什么选择）
   - 排除方案（考虑过但排除的选项）

### §Phase 1 — 技术设计

根据项目类型，从 Phase 0 结果推导技术设计。

#### BaseService 项目

1. **Proto 接口设计**：从 RQ 推导 gRPC 服务和方法，区分接口类型（Admin / Inner / Open），定义请求/响应关键字段
2. **枚举定义**：列出需要的枚举名、值、用途
3. **异常定义**：列出错误码、错误名、触发场景
4. **InnerRPC 依赖包装**：列出需包装的依赖服务及方法
5. **国际化 Key**：列出需要的 i18n key 及中英文

#### 业务项目

1. **Ent Schema 设计**：从实体分析推导字段、类型、约束、索引、关系映射
2. **Proto 接口设计**：同 BaseService
3. **枚举定义**
4. **异常定义**

#### 网关类项目

1. **Proto 接口设计**：同 BaseService
2. **Proxy 路由设计**：列出需代理的接口及对应的上游服务

### §框架适配

加载 **kratos-patterns** skill，用框架知识**扩充 Phase 1 的技术设计**：

1. 对照 kratos-patterns 中对应项目类型的规范，补充 Phase 1 设计表中未覆盖的框架层面考量（如 Ent Mixin 选型、Proto 包路径规划、接口归属判定等）
2. 标注每条设计决策对应的框架约束来源（哪条 kratos-patterns 规范）
3. 识别设计中可能与框架规范冲突的点，标记为风险

### §写入规则

- 写入路径：`specs/<feature>/plan.md`
- 保持模板章节结构
- Phase 0/1 的调研和设计内容必须可追溯到 spec 中的 RQ 编号

### §质量校验

| # | 校验项 | 标准 |
|---|--------|------|
| 1 | RQ 全覆盖 | 每条 RQ 至少映射到一个设计项或实施步骤 |
| 2 | 实体完整 | spec 中的关键实体均已在 Phase 0 分析 |
| 3 | 接口完整 | 每个用户行为都有对应的 Proto 方法 |
| 4 | 步骤可执行 | 每步有明确产出物 |
| 5 | 依赖已识别 | 所有 InnerRPC 调用已调研并标记状态 |
| 6 | 项目类型正确 | 判定结果与 SERVER_NAME 一致 |
| 7 | 测试策略完整 | 每个含业务逻辑的 RQ 至少有一个测试目标 |

## 输出

- `specs/<feature>/plan.md`（填充后的实施计划）

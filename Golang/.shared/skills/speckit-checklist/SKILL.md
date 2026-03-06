---
name: speckit-checklist
description: |
  用于 Speckit 研发流程的交付物质量保障，包括基于当前特性和文档自动生成多维度的交叉验收检查清单 (Checklist)。适用于在编写测试用例前，或进入交付阶段需明确列出需求覆盖路径、边界碰撞测试、异常阻断用例等可供人工排查或自动触发验证项的场景。触发关键词包括 checklist、检查单、验收清单、覆盖边界、可测试性。
---

# Spec Kit Checklist Skill（中文）

## 何时使用

- **推荐时机**：clarify 完成、spec 变为 Ready 后，plan 之前——确认需求本身质量过关再进入技术设计
- 也可在流水线任意阶段使用（spec 存在即可）
- 检查的是“需求写得是否好”，而非验证代码实现是否正确

## 核心理念

- 清单是"需求的单元测试"
- 检查的是需求质量：完整性、清晰性、一致性、可测性、覆盖度
- **不是**验证实现行为（不能写"调用接口是否返回成功"这类条目）

## 输入

- 用户说明的清单焦点与范围
- 可用的 spec/plan/tasks（按需读取）

## 工作流

### §前置检查

1. 确认 `specs/<feature>/` 目录存在且至少有 `spec.md`。
2. 若 spec 不存在，终止并提示先运行 `/specify`。
3. 读取 `./templates/checklist-template.md` 获取清单模板。

### §澄清问题

动态提问（最多 3 个，必要时再加 2 个）：

- 只问会实质改变清单内容的问题
- 可问的维度：范围、风险优先级、深度、受众、排除边界
- 已明确的信息不重复追问
- 若无法互动，默认：深度=Standard，受众=Author，焦点=相关性最高两类

### §主题确定

合并用户输入 + 澄清答案，确定清单主题。常见主题：

| 主题 | 适用场景 |
|------|---------|
| 功能完整性 | 检查 RQ 是否覆盖所有业务场景 |
| 接口设计 | 检查 Proto 接口定义的需求是否清晰 |
| 数据模型 | 检查实体/属性/关系的需求是否明确 |
| 安全 | 检查权限、认证、数据保护的需求是否完整 |
| 边界场景 | 检查异常流、限制条件是否有需求覆盖 |
| 非功能属性 | 检查性能、可用性等非功能需求是否量化 |

### §文档读取

按需读取 `specs/<feature>/` 下的文档：

| 文档 | 提取内容 |
|------|---------|
| spec.md | RQ 列表、关键实体、约束与边界 |
| plan.md（若存在） | 设计表、框架适配、风险项 |
| tasks.md（若存在） | 任务列表、Phase 覆盖 |

### §清单生成

1. 按主题分类生成条目
2. 条目编号从 `CHK001` 顺序递增
3. 每条必须是**检查需求质量**的问句，不是验证实现行为

#### 条目写作规则

**推荐维度**：
- 需求完整性（Completeness）
- 需求清晰性（Clarity）
- 需求一致性（Consistency）
- 可测性（Measurability）
- 场景覆盖（Scenario Coverage）
- 边界覆盖（Edge Case Coverage）
- 非功能需求（Non-Functional）
- 依赖与假设（Dependencies）
- 歧义与冲突（Ambiguities）

**条目模式示例**：
- "是否为所有失败模式定义了错误处理需求？[Gap]"
- "'高效'是否被量化为明确阈值？[Clarity, Spec §RQ-003]"
- "实体 A 和实体 B 的所属关系在 spec 和 plan 中是否一致？[Consistency]"
- "该需求能否客观验证？[Measurability]"

**追踪性要求**：
- 至少 50% 条目必须含追踪引用：
  - `[Spec §RQ-xxx]`
  - `[Gap]` / `[Ambiguity]` / `[Conflict]` / `[Assumption]`

**严禁内容**：
- "Verify/Test/Confirm/Check + 实现行为"
- 调用、执行、返回等实现测试语句
- 测试用例/测试流程/QA 执行步骤
- 框架、算法、接口实现细节

### §写入规则

- 写入路径：`specs/<feature>/checklists/<domain>.md`
- 使用 checklist-template.md 的章节结构
- 每次运行创建新文件；若同名存在则追加
- domain 即主题名（如 `completeness.md`、`api-design.md`、`security.md`）

### §输出报告

完成后报告：
- 清单文件路径
- 条目总数
- 焦点主题
- 深度级别
- 纳入的用户必选项

## 输出

- `specs/<feature>/checklists/<domain>.md`

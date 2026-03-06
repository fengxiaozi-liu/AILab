---
description: 对已有 spec 进行深度澄清，逐题与用户对话并将答案回写到 spec。
agent: clarify
---

# /clarify 命令

调用 **clarify** agent 对已有 spec 进行深度澄清。

## 何时使用

- `/specify` 生成的 spec 中存在待澄清项（CQ）
- 需要对 spec 做更深入的补漏和消歧
- specify 完成后建议的下一步

## 前置条件

已使用 `/specify` 生成：
- `specs/<feature>/spec.md` — 包含需求描述和待澄清问题

## 产出物

- 更新后的 `specs/<feature>/spec.md`（CQ 已填入澄清结果，需求描述已同步更新）

## 工作流程

```
/clarify [澄清方向或约束（可选）]
    ↓
加载 spec → 按 9 大维度评估覆盖状态
    ↓
合并已有 CQ + 新发现问题 → 生成问题队列
    ↓
逐题提问（一次一题，含推荐答案）
    ↓ 用户回答
增量回写 → 同步更新需求描述 → 写后校验
    ↓ 全部澄清或用户结束
Status → Ready，汇报结果
```

## 示例

### 深度澄清

```
User: /clarify

Agent:
已加载 spec：`specs/manage-portfolio/spec.md`

覆盖评估：
- 功能范围与行为：Clear
- 领域与数据模型：Partial
- 边界与失败处理：Missing

CQ-003（领域与数据模型）：股单中的标的是否支持多种资产类型（股票、ETF、期权）？
  推荐：仅支持股票（降低初版复杂度）
  - A. 仅股票
  - B. 股票 + ETF
  - C. 全部资产类型

请选择或回答。

WAIT_CONFIRM
```

```
User: A

Agent:
✅ CQ-003 已回写 → 仅支持股票
需求描述已同步更新。

CQ-004（边界与失败处理）：删除股单时，若股单内有标的，是否允许强制删除？
  推荐：禁止删除（需先清空标的）
  - A. 禁止删除，提示先清空
  - B. 级联删除标的
  - C. 软删除（标记归档）

WAIT_CONFIRM
```

```
User: A

Agent:
DONE

已更新：`specs/manage-portfolio/spec.md`
- 已澄清 2 题（CQ-003, CQ-004）
- 变更章节：关键实体、约束与边界
- 覆盖状态：全部 Clear
- Status → Ready

下一步建议由主 Agent 决定。
```

## 后续命令

- 由主 Agent 根据项目流程决定下一步

---

## 调用模板（主 Agent 使用）

调用 clarify 时，除通用字段外，需附带以下专用字段：

```markdown
## Spec 路径
- [specs/xxx/spec.md]

## 澄清方向（可选）
- [用户指定的澄清重点或约束，如"聚焦权限设计"或"关注并发场景"]
```

**说明**：clarify 根据 spec 中已有 CQ 和覆盖评估结果自行决定提问顺序

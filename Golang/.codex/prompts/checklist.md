适用场景：
- 需要为当前特性生成需求质量或验收检查清单

执行：
- 读取必要上下文
- 调用 `speckit-checklist`
- 生成或更新 `specs/<feature>/checklists/<domain>.md`

关键约束：
- 检查需求质量，不验证实现行为
- 条目应可检查、可复审

用户输入：`$ARGUMENTS`

适用场景：
- `spec.md`、`plan.md`、`tasks.md` 已生成，需要做只读一致性检查

执行：
- 读取 `specs/<feature>/spec.md`、`plan.md`、`tasks.md`
- 调用 `speckit-analyze`
- 输出分析报告

关键约束：
- 严格只读
- 优先识别高严重度问题

用户输入：`$ARGUMENTS`

适用场景：
- 需求、计划、任务、审查结果都已具备，需要生成交付总结

执行：
- 读取 `specs/<feature>/spec.md`、`plan.md`、`tasks.md`、`review.md`
- 调用 `speckit-summary`
- 生成或更新 `specs/<feature>/summary.md`

关键约束：
- 只做压缩汇总，不重复长篇原文
- 数字与结论以源文档为准

用户输入：`$ARGUMENTS`

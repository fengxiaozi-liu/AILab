适用场景：
- `plan.md` 已确定，需要拆成可执行任务

执行：
- 读取 `specs/<feature>/plan.md`、`specs/<feature>/spec.md` 与 `skills/speckit-tasks/templates/tasks-template.md`
- 调用 `speckit-tasks`
- 必要时加载当前项目适配器 skill（如 `kratos-patterns`）
- 生成或更新 `specs/<feature>/tasks.md`

关键约束：
- 任务应可执行、可验证、可追踪
- Phase 顺序由项目适配器或当前项目事实决定

完成后：
- 建议进入 `/analyze` 或 `/implement`

用户输入：`$ARGUMENTS`

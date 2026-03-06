适用场景：
- `spec.md` 已基本收敛，需要生成实施计划

执行：
- 读取 `specs/<feature>/spec.md` 与 `skills/speckit-plan/templates/plan-template.md`
- 调用 `speckit-plan`
- 必要时加载当前项目适配器 skill（如 `kratos-patterns`）
- 生成或更新 `specs/<feature>/plan.md`

关键约束：
- 先确认 spec 达到可规划状态
- speckit 负责标准化规划流程，项目判定与框架约束由项目适配器提供

完成后：
- 建议进入 `/tasks`

用户输入：`$ARGUMENTS`

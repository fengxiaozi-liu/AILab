适用场景：
- 实现已完成，需要做代码审查

执行：
- 读取 `specs/<feature>/spec.md`、`plan.md`、`tasks.md` 与实现代码
- 调用 `speckit-code-review`
- 必要时加载当前项目适配器 skill
- 输出 review 结论

关键约束：
- 发现项优先
- 优先识别 bug、回归、边界遗漏、缺少测试和架构违规

用户输入：`$ARGUMENTS`

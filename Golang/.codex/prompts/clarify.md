适用场景：
- `spec.md` 中仍有待澄清项，或需求边界尚未收敛

执行：
- 读取 `specs/<feature>/spec.md`
- 调用 `speckit-clarify`
- 按需回写 spec，并保留未解决项

关键约束：
- 只澄清会影响设计、实现、验收的关键信息
- 不做高风险默认假设

完成后：
- 若 spec 已就绪，建议进入 `/plan`

用户输入：`$ARGUMENTS`

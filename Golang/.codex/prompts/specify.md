适用场景：
- 用户给出自然语言需求，希望整理成结构化 spec

执行：
- 读取必要仓库事实与 `skills/speckit-specify/templates/spec-template.md`
- 调用 `speckit-specify`
- 生成或更新 `specs/<feature>/spec.md`

关键约束：
- 只写需求，不提前展开实现方案
- 缺失信息整理为 CQ

完成后：
- 有 CQ 则建议进入 `/clarify`
- 无 CQ 则建议进入 `/plan`

用户输入：`$ARGUMENTS`

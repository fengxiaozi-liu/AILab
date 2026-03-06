---
description: 基于已就绪的 spec 生成 Kratos 微服务实施计划。
---

# /plan 工作流

在 `Status` 为 `Ready` 的 `spec.md` 基础上生成架构实施计划，包含技术调研和框架设计适配。

## 适用场景
在完成需求说明 `/specify`（和可能的澄清 `/clarify`）之后，你需要启动本动作来规划系统技术架构、框架与测试方案（Phase 0/1）。 

## 执行步骤

1. **执行基础检查**
   - 读取 `.gemini/skills/speckit-plan/SKILL.md` 的计划方法论。
   - 打开包含功能的 `specs/<feature>/spec.md`。如果 `Status` 值不是 Ready 或该文件不存在，向用户说明必须先进行 `/specify` 或 `/clarify`，并中断计划流程。
   - 读取 `.gemini/templates/plan-template.md` 模板内容获取架构大纲要求。
   - 若发现虽然标了 Ready 但业务逻辑存在严重跳跃性遗漏，主动提示这可能影响计划、但不再强行要求返回 clarifying。

2. **Phase 0 调研与微服务架构适配 (Kratos Patterns)**
   - 根据用户的 $ARGUMENTS，调用代码阅读及文件列出工具去梳理本工程 `README.md` 与可能存在的 `docs`。
   - 使用 `.gemini/skills/kratos-patterns` 查阅 Kratos 微服务的架构指引和原则设计。
   - 判断当前项目和新需求的**项目类型**。
   - 梳理需要做改动的服务或库。

3. **Phase 1 设计展开**
   - 设计阶段包含模型结构定义（Entity/Proto 等）、层级定义（biz/data/service）。
   - 将这部分的细节依照 `plan-template.md` 骨架及上一步获得的 `kratos-patterns` 标准要求进行深度转化规划。

4. **架构审查输出方案及提问**
   - 检查最终的方案是否满足 `.gemini/skills/speckit-plan/SKILL.md` 中说明的质量校验标准。
   - 所有的设计都应当是在本地项目实际基础上的结果（即如果文件路径是啥、模块名是啥你要从代码里抽）。
   - 最后写入输出至 `specs/<feature>/plan.md`。

5. **下一步编排建议**
   - 生成后，总结汇报**设计决策摘要**给用户评审。
   - 在本回合向用户发问："架构方案已出炉，待您确认后，我们可以触发 **/tasks** 为这段代码拆解可执行任务并准备编码，需要继续么？" 

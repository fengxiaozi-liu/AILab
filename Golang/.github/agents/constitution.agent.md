---
name: constitution
description: 创建或更新项目宪法，并同步校验关键规则与运行时指引。
---

# 宪法治理专员

## 职责

根据用户输入与项目上下文抽象治理原则，维护 `.github/.specify/memory/constitution.md`，
并同步更新规则细则与运行资产，保持“宪法总纲 + rules 落地”的一致性。

## 执行步骤

1. 路径判定
- 优先使用 `.github/.specify/memory/constitution.md`。
- 若不存在则创建该文件，不回退到其他路径。

2. 原则提炼
- 用户输入优先；缺失部分从仓库上下文补全。
- 从 `.github/rules/project/` 抽取并归并为宪法原则，不在宪法重复细则。
- 若已有宪法，先继承再修订，不全量重写。
- 每条原则保持可引用锚点 `Principle-ID`（如 `[KRATOS_PRINCIPLE_I]`）。
- 每条原则仅绑定一个 `Rule Source` 文件，禁止多文件并列引用。

3. 治理完善
- 补充修订流程、版本策略、合规检查要求。
- 明确 Rule Application：细则以 `.github/rules/project/*.md` 与其 demo 为落地依据。
- 顶部写入 Sync Impact Report（HTML 注释）。

4. 联动更新
- 同步更新 `rules/project/*`（新增/修订/归档规则与 demo）。
- 同步更新相关运行资产：`.github/agents/*`、`.github/prompts/*`。

5. 一致性检查
- 核对 `.github/copilot-instructions.md` 是否声明宪法地位与调用入口。
- 核对 `rules/project/*` 与宪法章节是否一一可映射。
- 发现冲突时列出冲突文件与修订建议。

6. 输出结果
- 输出新版本号、升级原因、受影响文件、待办项。

## 关键约束

- 宪法内容优先复用现有规则，不平行造新体系。
- 宪法仅定义项目治理原则，不吸收 AI 编排规则为项目原则。
- 宪法不重复 rules 细则；细则与 demo 仅保留在 `rules/project/*`。

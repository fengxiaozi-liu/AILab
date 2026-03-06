# Codegen Rule

## Principles

- 生成链路必须完整，生成物不可手改。

## Specification

- 修改 proto、wire、ent schema 后运行对应生成命令。
- 提交前至少完成最小 build 验证。

## Prohibit

- 禁止只改源文件不更新生成物。

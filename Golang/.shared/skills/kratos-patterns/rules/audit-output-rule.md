# Audit Output Rule

## Specification

- 开始前输出：
  - `WorkDomain`
  - `Evidence`
  - `SubSkills`
  - `Rules`
- 提交前输出：
  - 各 skill 规则检查结果
  - codegen / build / test 验证结果

## Prohibit

- 禁止只给结论不附带证据。
- 禁止遗漏本次实际加载的规则文件。

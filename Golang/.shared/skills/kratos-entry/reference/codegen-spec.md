# Codegen Reference

## 这个主题解决什么问题

说明 Kratos 项目里 proto、wire、ent 等生成链路通常如何执行，以及常见的本地开发流程。

## 适用场景

- 修改 proto、wire provider、ent schema
- 排查生成物与源码不一致
- 补齐本地生成与构建验证步骤

## 设计意图

Codegen 参考的作用是帮助把“改源文件”和“同步生成链路”视为同一个实现动作，而不是两个割裂步骤。

- proto、wire、ent 都属于源定义驱动的生成体系，理解这点后更容易主动补生成和校验。
- 生成链路清楚时，可以更快判断问题出在源文件、工具版本还是产物漂移。
- 本地流程稳定后，评审和 CI 看到的差异也会更可预测。

## 实施提示

- 先确认改动触发的是哪条生成链路。
- 再执行对应命令，并检查差异是否与预期一致。
- 如果生成结果超出预期范围，优先回看源定义而不是直接修补产物。

## 典型生成对象

- Proto：API 契约、DTO、validate 代码
- Wire：依赖注入装配代码
- Ent：ORM schema 生成代码

## 推荐执行方式

### Wire

```bash
cd cmd/server
wire
```

### Proto

```bash
make generate
```

### 常见本地流程

```text
修改源文件
-> 运行对应生成命令
-> 检查差异
-> 执行最小 build 验证
```

## 代码示例参考

```bash
# proto 变更后
make generate

# wire 变更后
cd cmd/server && wire

# 最小验证
go build ./...
```

## 常见坑

- 本地工具版本与 CI 不一致
- 修改了 proto 或 schema，但漏跑生成命令
- 生成结果漂移后难以判断问题是工具还是源文件

## 相关 Rule

- `../rules/codegen-rule.md`
- `../rules/wire-rule.md`

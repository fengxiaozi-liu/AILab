---
name: golang-patterns
description: |
  Go 语言模式库：
  功能: Go 语言层面的惯用法与可复用代码模式（errors/context/concurrency/defer/generics/performance/design/style）
---

# Go Language Patterns

## 必读规则

- `./rules/errors-rule.md`
- `./rules/context-rule.md`
- `./rules/concurrency-rule.md`
- `./rules/defer-resource-rule.md`
- `./rules/generics-rule.md`
- `./rules/performance-rule.md`
- `./rules/style-organization-rule.md`

## 按需参考

- 错误处理：`patterns/errors.md`
- Context：`patterns/context.md`
- 并发：`patterns/concurrency.md`
- 资源释放：`patterns/defer-resource.md`
- 泛型：`patterns/generics.md`
- 性能：`patterns/performance.md`
- 设计模式：`patterns/design-patterns.md`
- 风格组织：`patterns/style-organization.md`

## 读取顺序

先按当前任务读取对应 `./rules/*.md`，再读取相关 `patterns/*.md` 获取推荐写法、模板和示例，不要一次性全部加载。

## 触发时机

涉及 Go 语言层面的实现与重构时加载本技能。

## 目录结构

```
golang-patterns/
└── patterns/
    ├── errors.md           错误处理
    ├── context.md          Context 使用
    ├── concurrency.md      并发与同步
    ├── defer-resource.md   Defer 与资源释放
    ├── generics.md         泛型
    ├── performance.md      性能意识
    ├── design-patterns.md  设计模式
    └── style-organization.md 风格与可读性
```

## 加载策略

| 场景 | 加载文件 |
|------|----------|
| 错误处理 | `patterns/errors.md` |
| Context 使用 | `patterns/context.md` |
| 并发编程 | `patterns/concurrency.md` |
| 资源释放 | `patterns/defer-resource.md` |
| 泛型 | `patterns/generics.md` |
| 性能优化 | `patterns/performance.md` |
| 设计模式 | `patterns/design-patterns.md` |
| 代码风格 | `patterns/style-organization.md` |
| Go 语言规则 | `rules/*.md` |

## 使用范围

- ✅ Go 语言层面的写法/惯用法
- ✅ 可复用的代码模式
- ❌ 不负责框架/项目级规范

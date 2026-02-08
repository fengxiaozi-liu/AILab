# Go Language Patterns

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

## 使用范围

- ✅ Go 语言层面的写法/惯用法
- ✅ 可复用的代码模式
- ❌ 不负责框架/项目级规范（交给 kratos-patterns）

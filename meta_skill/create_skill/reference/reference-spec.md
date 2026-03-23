# Reference Spec

## reference/ 文件的职责

SKILL.md 里的约束是声明性的（做什么 / 不做什么）。
reference/ 解决的是：**LLM 知道规则，但不知道规则在代码上长什么样。**

每个 reference 文件 = 一组知识点的精确示例集，不是规范文档。

---

## 四种内容类型

### 1. 最小正反例

去掉所有无关字段，只保留能说明边界的最少代码量。

```proto
// ✅ 局部结构嵌套
message ReviewAccountRequest {
  message RejectPageItem { string page_code = 1; }
  repeated RejectPageItem reject_page_list = 3;
}

// ❌ 局部结构平铺顶层
message RejectPageItem { string page_code = 1; }
```

### 2. 边界案例（⚠️）

正反例只覆盖黑白两端，边界案例覆盖灰色地带。

```proto
// ⚠️ Filter 只被 1 个 RPC 用 → 嵌套
message ListAccountRequest {
  message Filter { uint32 status = 1; }
  Filter filter = 1;
}

// ⚠️ Filter 被 2 个以上 RPC 用 → 提升为顶层
message AccountFilter { uint32 status = 1; }
message ListAccountRequest   { AccountFilter filter = 1; }
message SearchAccountRequest { AccountFilter filter = 1; }
```

### 3. 决策表

多个合法选项时，正反例无法指导选哪个，用决策表。

```markdown
| 条件 | 做法 |
|------|------|
| 只服务于 1 个 RPC | 嵌套 message |
| 被 2+ 个 RPC 共用 | 顶层 message |
| 跨 proto 文件共用 | 提取到 base/* |
```

### 4. 组合场景

多条规则同时生效时，代码整体长什么样。
用一个完整真实示例覆盖本文件所有核心约束。

---

## 文件结构模板

```markdown
# {主题} Spec

## {知识点 A}

决策表（可选，多选项时必须有）

| 条件 | 做法 |
|...   |...   |

正反例：
// ✅ 正例
// ❌ 反例
// ⚠️ 边界案例

---

## {知识点 B}
...

---

## 组合场景

涵盖本文件所有约束的完整示例

---

## 常见错误模式

// ❌ 错误模式 1
// ❌ 错误模式 2
```

---

## 质量检查

| 检查项 | 判断标准 |
|--------|----------|
| 能否去掉文字只留代码，仍能传达同等信息 | 如果能 → 写得足够精简 |
| 正例和反例的差异是否最小化 | 差异越小，边界越清晰 |
| 每个知识点是否独立成块 | 不混杂，LLM 能精准定位 |
| 是否重复了 SKILL.md 的约束文字 | MUST NOT 重复，只做实例化 |
| 是否有组合场景 | 有 → 多规则叠加时有参照 |

---

## 正反例：reference 文件整体写法

```markdown
# ✅ 正确写法 —— 以示例为主，prose 极少

## ID 字段命名

| 条件 | 做法 |
|------|------|
| 处于聚合根语境内 | 直接用 id |
| 跨聚合引用 | 用 {aggregate}_id |

// ✅
message GetAccountRequest { uint32 id = 1; }

// ❌ 当前语境已是 Account，仍追加前缀
message GetAccountRequest { uint32 account_id = 1; }
```

```markdown
# ❌ 错误写法 —— 纯 prose，无示例

## ID 字段命名

当 request 处在聚合根语境时，主对象 ID 直接使用 id。
只有跨聚合根引用时才使用 account_id 这类命名。
优先使用简洁命名，避免冗余的聚合根前缀。
```

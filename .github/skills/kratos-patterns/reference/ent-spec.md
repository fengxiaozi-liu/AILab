# Ent Schema 规范

## 概述

本项目使用 [Ent](https://entgo.io/) 作为 ORM 框架，Schema 定义位于 `internal/data/ent/schema/` 目录。

## Schema 模板

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
	"linksoft.cn/trader/internal/data/ent/mixin"
	baseenum "linksoft.cn/trader/internal/enum/base"
	orderenum "linksoft.cn/trader/internal/enum/order"
)

// OrderTradeCharge holds the schema definition for the OrderTradeCharge entity.
type OrderTradeCharge struct {
	ent.Schema
}

// Fields of the OrderTradeCharge.
func (OrderTradeCharge) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("id").Comment("主键ID"),
		field.Uint32("order_id").Default(0).Comment("订单ID"),
		field.String("order_no").MaxLen(32).Default("").Immutable().Comment("订单号"),
		field.Uint32("charge_id").Default(0).Comment("收费项ID"),
		field.String("currency").MaxLen(10).Default("").GoType(baseenum.Currency("")).Immutable().Comment("收费币种"),
		field.Float("charge_amount").DefaultFunc(func() decimal.Decimal { return decimal.Zero }).Comment("收费金额").GoType(decimal.Decimal{}).SchemaType(map[string]string{dialect.MySQL: "decimal(20,4)"}),
		field.Float("charge_measure").DefaultFunc(func() decimal.Decimal { return decimal.Zero }).Comment("计费度量").GoType(decimal.Decimal{}).SchemaType(map[string]string{dialect.MySQL: "decimal(20,4)"}),
		field.Uint8("type").Default(1).GoType(orderenum.TradeChargeType(0)).Immutable().Comment("类型 1 券商收费项 2 代收费项"),
		field.Uint32("charge_package_id").Default(0).Comment("收费套餐ID"),
		field.String("security_type").MaxLen(512).Default("").Comment("收费证券类型 逗号拼接多个"),
		field.String("side").MaxLen(512).Default("").Comment("收费方向 逗号拼接多个"),
		field.Uint8("calculate_type").Default(0).GoType(orderenum.TradeChargeCalculateType(0)).Comment("计算方式 1 以金额计算 2 以数量计算 3 每笔固定"),
		field.Float("calculate_value").DefaultFunc(func() decimal.Decimal { return decimal.Zero }).Comment("计算数值").GoType(decimal.Decimal{}).SchemaType(map[string]string{dialect.MySQL: "decimal(20,6)"}),
		field.Float("max").DefaultFunc(func() decimal.Decimal { return decimal.Zero }).Comment("最高收费").GoType(decimal.Decimal{}).SchemaType(map[string]string{dialect.MySQL: "decimal(20,4)"}),
		field.Float("min").DefaultFunc(func() decimal.Decimal { return decimal.Zero }).Comment("最低收费").GoType(decimal.Decimal{}).SchemaType(map[string]string{dialect.MySQL: "decimal(20,4)"}),
		field.Uint8("round_type").Default(0).GoType(baseenum.RoundType(0)).Comment("舍入方式 1 向上舍入 2 四舍五入 3 向下舍入"),
		field.Uint32("round_precision").Default(0).Comment("保留小数位数"),
		field.Text("extra").Comment("拓展参数").GoType(map[string]interface{}{}).ValueScanner(field.ValueScannerFunc[map[string]interface{}, *sql.NullString]{
			V: func(m map[string]interface{}) (driver.Value, error) {
				b, err := json.Marshal(m)
				if err != nil {
					return nil, err
				}
				return string(b), nil
			},
			S: func(ns *sql.NullString) (map[string]interface{}, error) {
				if !ns.Valid {
					return map[string]interface{}{}, nil
				}
				var m map[string]interface{}
				err := json.Unmarshal([]byte(ns.String), &m)
				return m, err
			},
		}),
	}
}

// Edges of the OrderTradeCharge.
func (OrderTradeCharge) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("order_trade_charge").Unique().Required().Field("order_id"),
	}
}

func (OrderTradeCharge) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("order_id").StorageKey("idx_orderid"),
		index.Fields("charge_id").StorageKey("idx_chargeid"),
		index.Fields("charge_package_id").StorageKey("idx_chargepackageid"),
	}
}

func (OrderTradeCharge) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.BaseMixin{},
		mixin.TimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}

func (OrderTradeCharge) Annotations() []schema.Annotation {
	return []schema.Annotation{
		schema.Comment("订单交易收费项表"),
		entsql.Annotation{Table: "link_order_trade_charge"},
	}
}

```

## 字段定义规范

### 基础字段类型

```go
func (Example) Fields() []ent.Field {
    return []ent.Field{
        // 主键 - 使用 uint32
        field.Uint32("id").Comment("主键ID"),
        
        // 字符串 - 指定最大长度
        field.String("name").MaxLen(50).Default("").Comment("名称"),
        
        // 文本 - 不限长度
        field.Text("description").Optional().Comment("描述"),
        
        // 整数
        field.Uint32("status").Default(0).Comment("状态"),
        field.Int64("amount").Default(0).Comment("金额"),
        
        // 浮点数 - 推荐使用 decimal
        field.Float("rate").DefaultFunc(func() decimal.Decimal { return decimal.Zero }).Comment("比率").GoType(decimal.Decimal{}).SchemaType(map[string]string{dialect.MySQL: "decimal(20,4)"}),
        
        // 布尔值
        field.Bool("is_active").Default(true).Comment("是否激活"),
    }
}
```

### JSON 字段

```go
// 使用 GoType + ValueScanner 实现 JSON 字段
field.Text("extra").Comment("拓展参数").
    GoType(map[string]interface{}{}).
    ValueScanner(field.ValueScannerFunc[map[string]interface{}, *sql.NullString]{
        V: func(m map[string]interface{}) (driver.Value, error) {
            b, err := json.Marshal(m)
            if err != nil {
                return nil, err
            }
            return string(b), nil
        },
        S: func(ns *sql.NullString) (map[string]interface{}, error) {
            if !ns.Valid {
                return map[string]interface{}{}, nil
            }
            var m map[string]interface{}
            err := json.Unmarshal([]byte(ns.String), &m)
            return m, err
        },
    }),
```

### 可选字段

```go
// 可选字段使用 Optional() 和 Nillable()
field.String("remark").Optional().Nillable().Comment("备注"),
field.Uint32("parent_id").Optional().Comment("父级ID"),
```

## Mixin 使用

### BaseMixin

提供基础数据库配置：

```go
// internal/data/ent/mixin/base.go
type BaseMixin struct {
    mixin.Schema
}

func (b BaseMixin) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.WithComments(true),
        entsql.Annotation{Charset: "utf8mb4", Collation: "utf8mb4_general_ci"},
    }
}
```

### TimeMixin

自动管理创建/更新时间：

```go
// internal/data/ent/mixin/time.go
type TimeMixin struct {
    mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Uint32("create_time").
            DefaultFunc(func() uint32 { return uint32(time.Now().Unix()) }).
            Immutable().
            Comment("创建时间"),
        field.Uint32("update_time").
            DefaultFunc(func() uint32 { return uint32(time.Now().Unix()) }).
            UpdateDefault(func() uint32 { return uint32(time.Now().Unix()) }).
            Comment("更新时间"),
    }
}

func (TimeMixin) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("create_time").StorageKey("idx_createtime"),
    }
}
```

### SoftDeleteMixin

软删除支持：

```go
// internal/data/ent/mixin/soft_delete.go
type SoftDeleteMixin struct {
    mixin.Schema
}

func (SoftDeleteMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Uint32("delete_time").Nillable().Optional().Comment("删除时间"),
    }
}

// 使用 SkipSoftDelete 跳过软删除过滤
ctx = mixin.SkipSoftDelete(ctx)
```

## 索引定义

```go
func (Example) Indexes() []ent.Index {
    return []ent.Index{
        // 单字段索引
        index.Fields("status").StorageKey("idx_status"),
        
        // 复合索引
        index.Fields("user_id", "create_time").StorageKey("idx_userid_createtime"),
        
        // 唯一索引
        index.Fields("code").Unique().StorageKey("uniq_code"),
    }
}
```

## 关联定义

```go
func (Order) Edges() []ent.Edge {
    return []ent.Edge{
        // 一对多：订单包含多个订单项
        edge.To("items", OrderItem.Type),
        
        // 多对一：订单属于用户
        edge.From("user", User.Type).
            Ref("orders").
            Unique(),
    }
}
```

## 表命名规范

```go
func (Example) Annotations() []schema.Annotation {
    return []schema.Annotation{
        schema.Comment("示例表"),           // 表注释
        entsql.Annotation{
            Table: "link_example",         // link_表名
        },
    }
}
```

## 代码生成

```shell
cd internal/data/ent
ent generate ./schema --feature privacy,entql,sql/lock,sql/modifier,intercept,schema/snapshot,sql/upsert --template ./template
```

生成的功能：
- `privacy` - 隐私保护
- `entql` - EntQL 支持
- `sql/lock` - 行级锁支持
- `sql/modifier` - SQL 修改器
- `intercept` - 拦截器支持
- `schema/snapshot` - Schema 快照
- `sql/upsert` - Upsert 支持

## 数据库迁移

项目启动时自动执行迁移（可配置关闭）：

```yaml
# configs/config.yaml
data:
  database:
    disable_migrate: "${DATABASE_DISABLE_MIGRATE:false}"
```

## 最佳实践

1. **字段注释**：每个字段都要添加 `Comment()`
2. **默认值**：合理使用 `Default()` 避免 NULL
3. **长度限制**：字符串字段使用 `MaxLen()` 限制
4. **索引命名**：使用 `StorageKey()` 显式命名索引
5. **软删除**：业务表统一使用 `SoftDeleteMixin`
6. **时间戳**：使用 Unix 时间戳 (uint32)，不使用 time.Time

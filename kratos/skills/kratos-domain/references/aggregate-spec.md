# Aggregate Reference

## 约束先看
必须遵守：

- 先识别稳定业务对象，再确定聚合根
- 聚合边界必须先于接口动作、页面结构和临时返回体
- `repo`、`usecase`、`proto`、`service` 的命名和职责都必须回收到同一聚合根
- 如果对象名称无法稳定映射回聚合根，说明建模还没收敛

## 使用说明

说明如何先识别实体，再确定聚合根，并由聚合根平级派生出 `biz` 领域对象、应用层输入输出对象、`repo`、`proto` 等结构。

## 常见场景

- 新增业务能力
- 重构现有聚合边界
- 判断一个对象、文件、接口应归属哪个聚合根
- 规划新功能的最小文件集

## 使用边界

聚合根是领域中心，不只是一个结构体名字。

- 实体先被识别，聚合根再负责收口边界
- 应用层输入输出对象和 `proto message` 都是聚合根的投影
- 两者是平级产物，不是上下游派生链
- `repo`、`usecase`、`proto`、`service` 的命名和职责都应围绕同一聚合根稳定下来

## 实施提示

- 先识别稳定实体，再决定谁是聚合根、谁是从属实体、谁只是关系视图
- 不要先按页面、接口动作或临时返回结构切模型
- 如果 `repo`、应用层输入输出对象、`proto` 的名字无法回到同一聚合根，通常说明建模还没收敛
- 对象在跨层传递、上下文传递、事件传递、依赖调用传递时，优先复用语义稳定的现有对象

## 推荐结构

- 一个聚合根围绕一个稳定业务对象组织
- 聚合根下可以有实体、值对象、状态枚举和关系视图
- 应用层输入输出对象、`repo`、`proto` 命名优先与聚合根保持一致

## 创建顺序

1. 先识别实体与稳定业务对象
2. 再确定聚合根名称与边界
3. 由聚合根平级派生 `biz` 领域对象、应用层输入输出对象、`repo`、`proto`
4. 再创建 `service`、`server`
5. 最后补 `wire`、`codegen` 与接入收口文件

## 最小文件集

通常包括：

- 聚合根或实体文件
- `usecase` 相关入参出参对象
- `repo` 接口与实现文件
- `usecase` 文件
- 对应 `proto` 文件

按需补充：

- `service` 文件
- `server` 注册文件
- `wire` / provider 文件
- gateway 代理文件

## 示例

### 标准模板

```text
Account
|- 核心字段：id、status、type
|- 从属实体：account_collect、account_flow_page
|- 关系视图：first_check_user、second_check_user
|- 平级派生：Account / AccountRepo / Account proto
```

### 典型实现

```go
type Account struct {
    ID                  uint32                 `json:"id"`
    Status              openenum.AccountStatus `json:"status"`
    Type                openenum.AccountType   `json:"type"`
    CreateTime          uint32                 `json:"create_time"`
    UpdateTime          uint32                 `json:"update_time"`
    FirstCheckUserInfo  *adminbiz.AdminUser    `json:"first_check_user_info"`
    SecondCheckUserInfo *adminbiz.AdminUser    `json:"second_check_user_info"`
    AccountCollectInfo  *AccountCollect        `json:"account_collect_info"`
    AccountFlowPageList []*AccountFlowPage     `json:"account_flow_page_list"`
}
```

### 反例

```go
type OpenAccountPageData struct{}
type SubmitAccountReplyStore struct{}
```

## 常见坑

- 先按接口动作拆对象，再反推领域模型
- 把临时展示结构误建成长期领域模型
- 应用层对象机械复制 `proto message`
- 新文件名无法稳定映射回聚合根

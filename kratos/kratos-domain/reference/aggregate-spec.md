# Aggregate Reference

## 这个主题解决什么问题
说明如何先识别实体，再确定聚合根，并由聚合根平级派生出 `biz` 领域对象、应用层输入输出对象、`repo`、`proto` 等文件与结构。

## 适用场景

- 新增业务能力
- 重构现有聚合边界
- 判断一个对象、文件、接口应归属哪个聚合根
- 规划新功能最小文件集

## 设计意图

聚合根是领域中心，不只是一个结构体名称。

- 实体先被识别，聚合根再负责收口边界。
- 应用层输入输出对象和 `proto message` 都是聚合根的投影。
- 两者是平级产物，不是上下游关系。
- `repo`、`usecase`、`proto`、`service` 的命名和职责都应围绕同一个聚合根稳定下来。

## 实施提示

- 先识别稳定实体，再决定谁是聚合根、谁是从属实体、谁只是关系视图。
- 不要先按页面、接口动作或临时返回结构来切模型。
- 如果 `repo`、应用层输入输出对象、`proto` 的名字无法回到同一个聚合根，通常说明建模还没收敛。
- 对象在跨层传递、上下文传递、事件传递、依赖调用传递时，优先复用已有且语义稳定的对象，不要仅为传递链路再包一层近义结构。

## 推荐结构

- 一个聚合根围绕一个稳定业务对象组织
- 聚合根下可以有实体、值对象、状态枚举和关系视图
- 应用层输入输出对象、`repo`、`proto` 命名优先与聚合根保持一致

## 聚合根驱动的平级文件创建

推荐顺序：

1. 先识别实体与稳定业务对象
2. 再确定聚合根名称与边界
3. 由聚合根平级派生 `biz` 领域对象、应用层输入输出对象、`repo`、`proto`
4. 再创建 `service`、`server`
5. 最后补 `wire`、codegen 与接入收口文件

最小文件集通常包括：

- 聚合根或实体文件
- `usecase` 相关入参出参对象
- `repo` 接口与实现文件
- `usecase` 文件
- 对应 `proto` 文件

按需出现的文件通常包括：

- `service` 文件
- `server` 注册文件
- `wire` / provider 文件
- gateway 代理文件

## 结构体字段排序

聚合根和领域结构体的字段顺序更适合稳定表达为：

1. 普通字段
2. 时间类型字段
3. `Info` 字段
4. `List` 字段

参与协议对齐、JSON 编解码或事件投递的聚合根对象、实体对象，建议统一补全 `json` tag，并使用 `snake_case`。

示例：

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

## 典型实现方式

1. 先确定业务主对象和边界
2. 再识别哪些字段属于聚合根自身，哪些属于从属实体或关系视图
3. 决定聚合根需要暴露哪些稳定能力给应用层输入输出对象、`repo` 和 `proto`
4. 由聚合根平级展开应用层对象、协议层对象和数据访问结构
5. 按聚合根扩展对应文件，而不是按单个接口动作临时创建镜像文件

## 对象复用说明

领域对象和应用层对象应优先复用，而不是随着传递场景不断派生近义结构。

- 复用已有对象的典型场景：
  - UseCase 把聚合对象继续传给下游依赖
  - EventBus 直接传递已有聚合对象或 store
  - 上下文相关处理直接复用已有稳定输入对象
- 需要新建对象的典型场景：
  - 对外协议投影，需要和 proto/message 对齐
  - 脱敏、裁剪、隔离边界
  - 原对象语义过大，无法安全表达当前边界

正例：

```go
return u.eventBus.Publish(ctx, &eventbus.Event{
    Topic:   topic,
    Payload: account,
})
```

```go
err := u.adminUserDepend.FillUsers(ctx, accountFlowPageList)
```

反例：

```go
type AccountAfterOpenEventPayload struct {
    Account *Account
}
```

```go
type AccountFlowPageContextDTO struct {
    Store *AccountFlowPageStore
}
```

## 标准模板

```text
Account
|- 核心字段：id、status、type
|- 从属实体：account_collect、account_flow_page
|- 关系视图：first_check_user、second_check_user
|- 平级派生：Account / AccountRepo / Account proto
```

## Good Example

```text
Account 聚合根
- AccountUseCase 负责业务编排
- Account 负责应用层输入输出
- AccountRepo 负责数据访问与装配
- Account proto message 负责协议表达
```

## 代码示例参考

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

type AccountCollect struct {
    ID        uint32 `json:"id"`
    AccountID uint32 `json:"account_id"`
    PageNo    uint32 `json:"page_no"`
}

type AccountFlowPage struct {
    ID        uint32 `json:"id"`
    AccountID uint32 `json:"account_id"`
    Title     string `json:"title"`
}
```

## 常见坑

- 先按接口动作拆对象，后续再倒推领域模型
- 把临时展示结构误建成长期领域模型
- 应用层对象机械复制 `proto message`
- 新文件名无法稳定映射回聚合根

## 相关 Rule

- `../rules/aggregate-rule.md`
- `../rules/layer-rule.md`

## 相关 Reference

- `./naming-spec.md`
- `./repo-spec.md`
- `./usecase-spec.md`

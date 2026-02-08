---
name: kratos-patterns
description: Kratos 微服务新增功能/修Bug/重构时的工作流程规范
---

# Kratos Patterns

## 触发时机

涉及**代码更新**（新增功能/修 Bug/重构）时加载本技能。

## 工作流程

### 获取当前工作域

通过.env.*中的SERVER_NAME来确认当前所在项目类型：

- BaseService项目：SERVER_NAME通常以`BaseService`结尾。
- 网关类项目：SERVER_NAME通常包含`GatewayService`或`OpenapiService`。
- 业务项目：不属于上述两类的项目。

### 根据项目类型执行以下工作流之一

- 以下三种工作流互斥，你只能执行其中一种，请根据当前工作域选择合适的工作流。
- 选择工作流后，你必须输出所选择的工作流名称。
- 你必须加载这个工作流中列出的所有参考文件，然后再进行工作。

#### BaseService项目（抽象定义）

##### 准备工作

1. 阅读[命名规范](./reference/naming-spec.md)，[项目规范](./reference/project-spec.md)，遵循其中的项目约定。
2. 阅读[Proto 规范](./reference/proto-spec.md)，了解如何定义proto文件（新增或者修改）。
3. 阅读[代码生成规范](./reference/codegen-spec.md)，了解如何生成代码。
4. 阅读[枚举规范](./reference/enum-spec.md)，了解如何编写枚举。
5. 阅读[异常规范](./reference/error-spec.md)，了解如何编写异常。
6. 阅读[InnerRPC 依赖包装规范](./reference/depend-spec.md)，了解如何包装innerRpc的depend。
7. 阅读[国际化规范](./reference/i18n-spec.md)，了解如何编写国际化。

##### 进行工作

1. 根据需求，设计或修改proto文件，定义服务接口和消息结构。
2. 进行代码生成，生成proto对应的go代码。
3. 定义或修改枚举类型，确保符合业务需求。
4. 定义或修改异常类型，确保覆盖所有可能的错误场景。
5. 包装所需的InnerRPC依赖，确保服务间通信的正确性。
6. 编写国际化内容，确保支持多语言环境。

#### 业务项目（业务实现）

##### 准备工作

1. 阅读[命名规范](./reference/naming-spec.md)，[项目规范](./reference/project-spec.md)，遵循其中的项目约定。
2. 向用户确认需求细节，需要操作哪些领域，明确业务逻辑与边界条件。
3. 阅读[Proto 规范](./reference/proto-spec.md)，了解如何使用proto文件（新增或者修改）。
4. 阅读[枚举规范](./reference/enum-spec.md)，了解如何使用枚举。
5. 阅读[异常规范](./reference/error-spec.md)，了解如何使用异常。
6. 阅读[InnerRPC 依赖包装规范](./reference/depend-spec.md)，了解如何使用innerRpc。
7. 阅读相关可能用到的innerRpc，参考[InnerRPC 依赖包装规范](./reference/depend-spec.md)。
8. 阅读[Ent 代码生成规范](./reference/ent-spec.md)，了解如何编写ent schema，如何进行代码生成。
9. 阅读[Layer 层规范](./reference/layer-spec.md)，了解如何实现biz，data，service层，如何编写具体的业务代码。
10. 阅读[Server 注册规范](./reference/server-spec.md)，了解如何进行服务注册。
11. 阅读[代码生成规范](./reference/codegen-spec.md)，了解如何生成wire代码。
12. 阅读[测试规范](./reference/testing-spec.md)，了解如何编写单元测试和集成测试。

##### 进行工作

1. 编写ent schema，定义数据模型。
2. 进行ent代码生成，生成对应的go代码。
3. 实现biz层，编写具体的业务逻辑代码。
4. 实现data层，编写具体的数据访问代码。
5. 实现service层，编写具体的API实现代码。
6. 进行服务注册，确保服务能够被正确发现和调用。
7. 生成wire代码，确保依赖注入的正确性。
8. 编写单元测试，确保各层代码的正确性和稳定性。
9. 进行集成测试，确保各层之间的协作正确性。
10. 编写文档，记录新增或修改的功能和使用方法。
11. 进行代码审查，确保代码质量符合项目标准。

#### 网关类项目（接口实现）

##### 准备工作

1. 阅读[命名规范](./reference/naming-spec.md)，[项目规范](./reference/project-spec.md)，遵循其中的项目约定。
2. 阅读[Proto 规范](./reference/proto-spec.md)，了解如何使用proto文件（新增或者修改）。
3. 阅读[枚举规范](./reference/enum-spec.md)，了解如何使用枚举。
4. 阅读[异常规范](./reference/error-spec.md)，了解如何使用异常。
5. 阅读[InnerRPC 依赖包装规范](./reference/depend-spec.md)，了解如何使用innerRpc。
6. 阅读[Gateway 层规范](./reference/gateway-spec.md)，了解如何实现gateway层，如何编写具体的代码。
7. 阅读[Server 注册规范](./reference/server-spec.md)，了解如何进行服务注册。
8. 阅读[代码生成规范](./reference/codegen-spec.md)，了解如何生成wire代码。

##### 进行工作

1. 实现proxy层，编写具体的代理逻辑代码。
2. 使用proxy进行网关的服务注册，确保服务能够被正确发现和调用。
3. 生成wire代码，确保依赖注入的正确性。
4. 进行代码审查，确保代码质量符合项目标准。

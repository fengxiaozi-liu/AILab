# Common Violations

## 1. 层级越界

- 在 `repo` 中直接推进业务状态
- 在 `service / listener / consumer` 中编排完整业务流程

## 2. 聚合根漂移

- 同一能力拆散到多个无关业务域
- 命名使用其他聚合根的术语

## 3. 稳定协议散落

- header / path / query key / content-type 留在实现文件中
- 第三方固定状态值直接散落在 `biz / data / service`

## 4. 注入前提被运行期兜底

- 对构造注入依赖在方法体内做 `nil` 判断
- 用普通运行期错误掩盖装配期失败

## 5. 公共能力误下沉

- 仍带业务语义的逻辑进入 `internal/pkg`
- 仅因“当前只用一次”就把稳定公共语义留在实现文件里

## 6. data-access 泄漏

- `biz` 直接感知 ent 结构、predicate、mutation 细节
- DTO、持久化模型、领域模型边界混写

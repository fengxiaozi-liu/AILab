# 防御式编程规则

## 原则

- 校验和错误语义要有单一可信来源。
- 不要只为了“防御”就在上层增加参数校验。
- 返回值语义要稳定，让调用方优先只处理 `err`。
- 出现错误就直接返回；只有需要业务错误转换时才配合 `error-rule.md`。

## 规范

- 普通查询链路允许 `Service` 或 `UseCase` 直接把参数透传到 `Repo`，默认不做前置校验。
- 由真正拥有该语义的层判断参数是否合法，以及应该返回什么错误。
- 只有当特殊调用方确实需要额外约束时，才在使用处就地校验；不要要求所有上层重复做同一类校验。
- 对 `(*T, error)`、`([]T, error)` 等二元返回，成功与失败状态必须清晰无歧义。
- 对“必须存在”的单对象查询，查不到时返回明确错误，不返回 `nil, nil`。
- 下层返回错误时，上层通常直接返回该错误；只有现有错误规范要求做语义转换时才转换。

## 禁止项

- 不要在 `UseCase` 或 `Service` 重复校验已经由 `Repo` 或其他下层负责的普通参数。
- 不要对必须存在的结果返回 `nil, nil`。
- 不要通过 `if obj != nil` 之类的补偿逻辑，把状态猜测转嫁给调用方。
- 不要用 `_, _ = call()`、忽略返回值、只记日志后伪装成功等方式吞掉错误。
- 不要把“是否存在 / 是否合法 / 是否成功”的歧义扩散给整条调用链。

## 正例

```go
func (u *AccountUseCase) GetAccountByID(ctx context.Context, id uint32) (*Account, error) {
	return u.repo.GetAccount(ctx, id)
}
```

```go
func (r *accountRepo) GetAccount(ctx context.Context, id uint32) (*Account, error) {
	info, err := query.First(ctx)
	if err != nil {
		return nil, err
	}
	return convert(info), nil
}
```

## 反例

```go
func (u *AccountUseCase) GetAccountByID(ctx context.Context, id uint32) (*Account, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	return u.repo.GetAccount(ctx, id)
}
```

```go
func (r *accountRepo) GetAccount(ctx context.Context, id uint32) (*Account, error) {
	if ent.IsNotFound(err) {
		return nil, nil
	}
	return convert(info), err
}
```

```go
_, _ = repo.GetAccount(ctx, id)
```

```go
if info != nil {
	return handle(info)
}
```

## Review Checklist

- 是否存在仅出于防御目的的上层重复参数校验？
- 是否有必须存在的查询返回了 `nil, nil`？
- 是否存在吞错、忽略错误或把错误伪装成成功？
- 是否因为下层语义不清，导致调用方不得不写 `if obj != nil` 一类补判？
- 业务错误转换是否与默认的“直接返回 err”路径保持分离？

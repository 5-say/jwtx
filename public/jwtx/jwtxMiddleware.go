package jwtx

import "context"

type ctxKey int

const MIDDLEWARE_RESULT ctxKey = iota

type MiddlewareResult struct {
	TokenID             uint64 // token ID
	AccountID           uint64 // 原始的账户 ID
	Terminal            string // 登录的终端名称
	MakeTokenIP         string // 首次请求生成 token 的 IP 地址
	ExpirationTimestamp int64  // token 过期时间戳
	NewToken            string // 刷新后的 token
}

func WithValue(ctx context.Context, response *CheckToken_Response) context.Context {
	return context.WithValue(ctx, MIDDLEWARE_RESULT, MiddlewareResult{
		TokenID:             response.TokenID,
		AccountID:           response.AccountID,
		Terminal:            response.Terminal,
		MakeTokenIP:         response.MakeTokenIP,
		ExpirationTimestamp: response.ExpirationTimestamp,
		NewToken:            response.NewToken,
	})
}

func GetValue(ctx context.Context) MiddlewareResult {
	return ctx.Value(MIDDLEWARE_RESULT).(MiddlewareResult)
}

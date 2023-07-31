package logic

import (
	"context"

	"github.com/5-say/zero-auth/private/jwtx/rpc/internal/svc"
	"github.com/5-say/zero-auth/public/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type MakeTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMakeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MakeTokenLogic {
	return &MakeTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 生成 token（登录）
func (l *MakeTokenLogic) MakeToken(in *jwtx.MakeToken_Request) (*jwtx.MakeToken_Response, error) {
	// todo: add your logic here and delete this line

	return &jwtx.MakeToken_Response{}, nil
}

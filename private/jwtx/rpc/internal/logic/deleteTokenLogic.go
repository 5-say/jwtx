package logic

import (
	"context"

	"github.com/5-say/zero-auth/private/jwtx/rpc/internal/svc"
	"github.com/5-say/zero-auth/public/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTokenLogic {
	return &DeleteTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 移除 token（安全退出）
func (l *DeleteTokenLogic) DeleteToken(in *jwtx.DeleteToken_Request) (*jwtx.DeleteToken_Response, error) {
	// todo: add your logic here and delete this line

	return &jwtx.DeleteToken_Response{}, nil
}

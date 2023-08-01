package logic

import (
	"context"

	"github.com/5-say/go-tools/tools/t"
	"github.com/5-say/zero-auth/private/jwtx/db/dao"
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
	// 初始化数据库
	q := dao.Common()

	// 查找 token
	token, err := q.Token.Where(
		q.Token.ID.Eq(in.TokenID),
		q.Token.AccountID.Eq(in.AccountID),
	).First()
	if err != nil {
		return nil, t.RPCError(err.Error(), "token does not exist")
	}

	// 分组校验
	if token.LoginGroup != in.Group {
		return nil, t.RPCError("group ["+in.Group+"] is not login group ["+token.LoginGroup+"]", "group fail")
	}

	// 获取配置
	c, ok := l.svcCtx.Config.JWTX[in.Group]
	if !ok {
		return nil, t.RPCError("group ["+in.Group+"] config does not exist", "group fail")
	}

	// 移除 token
	{
		if c.SingleEnd { // 强制单端登录

			_, err := q.Token.Where(
				q.Token.AccountID.Eq(token.AccountID),
				q.Token.LoginGroup.Eq(token.LoginGroup),
			).Delete()
			if err != nil {
				return nil, t.RPCError(err.Error(), "delete fail")
			}
		} else { // 多端登录

			_, err := q.Token.Where(
				q.Token.AccountID.Eq(token.AccountID),
				q.Token.LoginGroup.Eq(token.LoginGroup),
				q.Token.LoginTerminal.Eq(in.Terminal), // 支持跨端注销
			).Delete()
			if err != nil {
				return nil, t.RPCError(err.Error(), "delete fail")
			}
		}
	}

	return &jwtx.DeleteToken_Response{}, nil
}

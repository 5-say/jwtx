package logic

import (
	"context"
	"time"

	"github.com/5-say/go-tools/tools/random"
	"github.com/5-say/go-tools/tools/t"
	"github.com/5-say/zero-auth/private/jwtx/common"
	"github.com/5-say/zero-auth/private/jwtx/db/dao"
	"github.com/5-say/zero-auth/private/jwtx/db/dao/model"
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
	// 初始化数据库
	q := dao.Common()

	// 获取配置
	c, ok := l.svcCtx.Config.JWTX[in.Group]
	if !ok {
		return nil, t.RPCError("group ["+in.Group+"] config does not exist", "group fail")
	}

	// 清除旧 token
	{
		if c.SingleEnd { // 强制单端登录

			if _, err := q.Token.Where(
				q.Token.AccountID.Eq(in.AccountID),
				q.Token.LoginGroup.Eq(in.Group),
			).Delete(); err != nil {
				return nil, t.RPCError(err.Error(), "clear fail")
			}

		} else { // 多端登录

			if _, err := q.Token.Where(
				q.Token.AccountID.Eq(in.AccountID),
				q.Token.LoginGroup.Eq(in.Group),
				q.Token.LoginTerminal.Eq(in.Terminal),
			).Delete(); err != nil {
				return nil, t.RPCError(err.Error(), "clear fail")
			}
		}
	}

	var (
		now     = time.Now()
		expTime = now.Add(time.Duration(c.AccessExpireByHour) * time.Hour)
	)

	// 创建新 token
	token := model.Token{
		AccountID:      in.AccountID,
		LoginGroup:     in.Group,
		LoginTerminal:  in.Terminal,
		MakeTokenIP:    in.RequestIP,
		CreatedAt:      now,
		LastRefreshAt:  now,
		FinalRefreshAt: now,
		ExpirationAt:   expTime,
	}
	if err := q.Token.Create(&token); err != nil {
		return nil, t.RPCError(err.Error(), "create fail")
	}

	// 构造 token 字符串
	randomAccountID := random.Simple(l.svcCtx.Config.SimpleRandom).Encode(int64(in.AccountID), 8)
	tokenStr, err := common.MakeTokenStr(l.svcCtx.Config.TokenSecret, now, expTime, token.ID, randomAccountID)
	if err != nil {
		return nil, t.RPCError(err.Error(), "make fail")
	}

	return &jwtx.MakeToken_Response{
		Token: tokenStr,
	}, nil
}

package logic

import (
	"context"
	"time"

	"github.com/5-say/go-tools/tools/random"
	"github.com/5-say/go-tools/tools/t"
	"github.com/5-say/zero-auth/private/jwtx/common"
	"github.com/5-say/zero-auth/private/jwtx/db/dao"
	"github.com/5-say/zero-auth/private/jwtx/rpc/internal/svc"
	"github.com/5-say/zero-auth/public/jwtx"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTokenLogic {
	return &CheckTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 校验 token（拓展校验、刷新 token）
func (l *CheckTokenLogic) CheckToken(in *jwtx.CheckToken_Request) (*jwtx.CheckToken_Response, error) {
	// 解析 token 字符串
	claims, err := common.ParseToken(in.RequestToken, []byte(l.svcCtx.Config.TokenSecret))
	if err != nil {
		return nil, t.RPCError(err.Error(), "parse fail")
	}
	var (
		iat             = int64(claims["iat"].(float64))  // 签发时间
		exp             = int64(claims["exp"].(float64))  // 过期时间
		tid             = uint64(claims["tid"].(float64)) // token ID
		randomAccountID = claims["rai"].(string)          // 加密的账户 ID

		now = time.Now()
	)

	// token 过期时间校验
	if exp < now.Unix() {
		return nil, t.RPCErrorCode(err.Error(), "token has expired", codes.DeadlineExceeded)
	}

	// 初始化数据库
	q := dao.Common()

	// 查找 token
	m := q.Token
	token, err := m.Where(m.ID.Eq(tid)).First()
	if err != nil {
		return nil, t.RPCError(err.Error(), "not found")
	}

	// db 过期时间校验
	if token.ExpirationAt.Unix() < now.Unix() {
		return nil, t.RPCErrorCode(err.Error(), "account has expired", codes.DeadlineExceeded)
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

	// IP 一致性校验
	if c.CheckIP {
		if in.RequestIP != token.MakeTokenIP {
			return nil, t.RPCError("", "ip fail")
		}
	}

	// accountID 防篡改
	if random.Simple(l.svcCtx.Config.SimpleRandom).Decode(randomAccountID) != int64(token.AccountID) {
		return nil, t.RPCError("", "account fail")
	}

	// token 刷新校验
	var (
		newToken = ""
	)
	// l.Debug("iat", iat)
	// l.Debug("最后刷新时间", token.FinalRefreshAt.Unix())
	// l.Debug("刷新时间间隔", in.RefreshInterval)
	// l.Debug("上次刷新时间", token.LastRefreshAt.Unix())
	// l.Debug("并发容错时间", in.FaultTolerance)
	if iat == token.FinalRefreshAt.Unix() { // token 未刷新

		// 原始的 token 过期时间
		expTime := token.ExpirationAt

		// 需要刷新 token
		if iat+c.RefreshInterval < now.Unix() {

			// 自动续期
			if c.AutomaticRenewal {
				expTime = now.Add(time.Duration(c.AccessExpireByHour) * time.Hour)
			}

			// 构造 token 字符串（过期时间不变，签发时间顺延）
			newToken, err = common.MakeTokenStr(l.svcCtx.Config.TokenSecret, now, expTime, token.ID, randomAccountID)
			if err != nil {
				return nil, t.RPCError(err.Error(), "make fail")
			}

			// 更新数据库
			_, err = m.Where(m.ID.Eq(token.ID)).UpdateSimple(
				m.LastRefreshAt.Value(token.FinalRefreshAt),
				m.FinalRefreshAt.Value(now),
			)
			if err != nil {
				return nil, t.RPCError(err.Error(), "refresh fail")
			}
		}

		// 验证通过
		return &jwtx.CheckToken_Response{
			TokenID:             token.ID,
			AccountID:           token.AccountID,
			Terminal:            token.LoginTerminal,
			MakeTokenIP:         token.MakeTokenIP,
			ExpirationTimestamp: expTime.Unix(),
			NewToken:            newToken,
		}, nil

	} else if iat == token.LastRefreshAt.Unix() { // token 已刷新

		// 当前时间 未超出 并发容错时间（允许继续使用）
		if now.Unix() < token.FinalRefreshAt.Unix()+c.FaultTolerance {
			// 验证通过
			return &jwtx.CheckToken_Response{
				TokenID:             token.ID,
				AccountID:           token.AccountID,
				Terminal:            token.LoginTerminal,
				MakeTokenIP:         token.MakeTokenIP,
				ExpirationTimestamp: token.ExpirationAt.Unix(),
				NewToken:            "",
			}, nil
		}
	}

	return nil, t.RPCError("", "token refreshed")
}

package common

import (
	"time"

	"github.com/5-say/go-tools/tools/t"
	"github.com/golang-jwt/jwt"
)

// MakeTokenStr .. 构造 token 字符串
func MakeTokenStr(secretKey string, iatTime, expTime time.Time, tokenID uint64, randomAccountID string) (tokenStr string, err error) {
	// 规避 MySQL 数据库 DATETIME 字段类型存储时，毫秒位四舍五入的问题
	iat := t.Fix_MySQL_DATETIME(iatTime)
	exp := t.Fix_MySQL_DATETIME(expTime)

	claims := make(jwt.MapClaims)
	claims["iat"] = iat             // 签发时间
	claims["exp"] = exp             // 过期时间
	claims["tid"] = tokenID         // token ID
	claims["rai"] = randomAccountID // 加密的账户 ID
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

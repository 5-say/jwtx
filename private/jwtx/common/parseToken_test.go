package common_test

import (
	"testing"
	"time"

	"github.com/5-say/zero-auth/private/jwtx/common"
	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	var (
		secretKey       = "56789"
		iat             = time.Now()
		exp             = iat.Add(time.Second * 3600)
		tokenID         = uint64(1)
		randomAccountID = "KDUE7FL"
	)
	tokenStr, err := common.MakeTokenStr(secretKey, iat, exp, tokenID, randomAccountID)
	if assert.NoError(t, err) {
		claims, err := common.ParseToken(tokenStr, []byte(secretKey))
		if assert.NoError(t, err) {
			assert.Equal(t, uint64(claims["tid"].(float64)), tokenID)
		}
	}
}

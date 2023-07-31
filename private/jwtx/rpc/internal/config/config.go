package config

import (
	"github.com/5-say/go-tools/tools/db"
	"github.com/5-say/go-tools/tools/random"
	"github.com/5-say/zero-auth/public/jwtx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DB           db.Config
	SimpleRandom random.SimpleRandomConfig
	TokenSecret  string
	JWTX         map[string]jwtx.Config
}

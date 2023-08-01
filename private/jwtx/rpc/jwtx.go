package main

import (
	"flag"
	"fmt"

	"github.com/5-say/go-tools/tools/db"
	"github.com/5-say/zero-auth/private/jwtx/db/dao"
	"github.com/5-say/zero-auth/private/jwtx/rpc/internal/config"
	"github.com/5-say/zero-auth/private/jwtx/rpc/internal/server"
	"github.com/5-say/zero-auth/private/jwtx/rpc/internal/svc"
	"github.com/5-say/zero-auth/public/jwtx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/jwtx.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// 初始化公共数据库连接
	dao.SetCommon(c.DB, db.LogWriter())

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		jwtx.RegisterJwtxServer(grpcServer, server.NewJwtxServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

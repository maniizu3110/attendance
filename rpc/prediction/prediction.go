package main

import (
	"flag"
	"fmt"

	"github.com/maniizu3110/attendance/rpc/prediction/internal/config"
	"github.com/maniizu3110/attendance/rpc/prediction/internal/server"
	"github.com/maniizu3110/attendance/rpc/prediction/internal/svc"
	"github.com/maniizu3110/attendance/rpc/prediction/proto/add"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/prediction.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		add.RegisterPredictionServer(grpcServer, server.NewPredictionServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

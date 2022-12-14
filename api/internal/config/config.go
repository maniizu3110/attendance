package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Project    zrpc.RpcClientConf
	Attendance zrpc.RpcClientConf
	Prediction zrpc.RpcClientConf
}

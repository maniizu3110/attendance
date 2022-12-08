package svc

import (
	"github.com/maniizu3110/attendance/api/internal/config"
	"github.com/maniizu3110/attendancse/rpc/project/project"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	Project project.Project
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Project: project.NewProject(zrpc.MustNewClient(c.Project)),
	}
}

package svc

import (
	"github.com/maniizu3110/attendance/api/internal/config"
	"github.com/maniizu3110/attendance/rpc/project/adder"
)

type ServiceContext struct {
	Config config.Config
	Adder  adder.Adder // manual code
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

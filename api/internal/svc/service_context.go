package svc

import (
	"github.com/maniizu3110/attendance/api/internal/config"
	"github.com/maniizu3110/attendance/rpc/attendance/attendanceclient"
	"github.com/maniizu3110/attendance/rpc/prediction/prediction"
	"github.com/maniizu3110/attendance/rpc/project/project"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	Project    project.Project
	Attendance attendanceclient.Attendance
	Prediction prediction.Prediction
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		Project:    project.NewProject(zrpc.MustNewClient(c.Project)),
		Attendance: attendanceclient.NewAttendance(zrpc.MustNewClient(c.Attendance)),
		Prediction: prediction.NewPrediction(zrpc.MustNewClient(c.Prediction)),
	}
}

// Code generated by goctl. DO NOT EDIT!
// Source: attendance.proto

package server

import (
	"context"

	"github.com/maniizu3110/attendance/rpc/attendance/internal/logic"
	"github.com/maniizu3110/attendance/rpc/attendance/internal/svc"
	"github.com/maniizu3110/attendance/rpc/attendance/proto/attendance"
)

type AttendanceServer struct {
	svcCtx *svc.ServiceContext
	attendance.UnimplementedAttendanceServer
}

func NewAttendanceServer(svcCtx *svc.ServiceContext) *AttendanceServer {
	return &AttendanceServer{
		svcCtx: svcCtx,
	}
}

func (s *AttendanceServer) Add(ctx context.Context, in *attendance.AddReq) (*attendance.AddResp, error) {
	l := logic.NewAddLogic(ctx, s.svcCtx)
	return l.Add(in)
}

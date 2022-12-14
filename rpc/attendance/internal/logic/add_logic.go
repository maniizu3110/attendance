package logic

import (
	"context"

	"github.com/maniizu3110/attendance/rpc/attendance/internal/svc"
	"github.com/maniizu3110/attendance/rpc/attendance/proto/attendance"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddLogic) Add(in *attendance.AddReq) (*attendance.AddResp, error) {
	// todo: add your logic here and delete this line

	return &attendance.AddResp{}, nil
}

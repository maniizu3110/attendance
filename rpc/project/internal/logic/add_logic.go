package logic

import (
	"context"
	"fmt"

	"github.com/maniizu3110/attendance/rpc/project/internal/svc"
	"github.com/maniizu3110/attendance/rpc/project/proto/add"

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

func (l *AddLogic) Add(in *add.AddReq) (*add.AddResp, error) {
	fmt.Println("Add called")
	fmt.Println(in)
	return &add.AddResp{}, nil
}

package logic

import (
	"context"

	"github.com/maniizu3110/attendance/api/internal/svc"
	"github.com/maniizu3110/attendance/api/internal/types"
	"github.com/maniizu3110/attendance/rpc/project/project"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.AddReq) (resp *types.AddResp, err error) {
	l.svcCtx.Project.Add(l.ctx, &project.AddReq{Price: req.Price, Book: req.Book})
	return nil, err
}

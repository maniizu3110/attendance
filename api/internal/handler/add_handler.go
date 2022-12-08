package handler

import (
	"fmt"
	"net/http"

	"github.com/maniizu3110/attendance/api/internal/logic"
	"github.com/maniizu3110/attendance/api/internal/svc"
	"github.com/maniizu3110/attendance/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddReq
		fmt.Println("add_handler")
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewAddLogic(r.Context(), svcCtx)
		resp, err := l.Add(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

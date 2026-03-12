package user

import (
	"net/http"

	"easy_chat/apps/user/api/internal/logic/user"
	"easy_chat/apps/user/api/internal/svc"
	"easy_chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewFindUserLogic(r.Context(), svcCtx)
		resp, err := l.FindUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

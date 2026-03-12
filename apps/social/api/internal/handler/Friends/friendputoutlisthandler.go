package Friends

import (
	"net/http"

	"easy_chat/apps/social/api/internal/logic/Friends"
	"easy_chat/apps/social/api/internal/svc"
	"easy_chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FriendPutOutListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendPutOutListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := Friends.NewFriendPutOutListLogic(r.Context(), svcCtx)
		resp, err := l.FriendPutOutList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package Friends

import (
	"net/http"

	"easy_chat/apps/social/api/internal/logic/Friends"
	"easy_chat/apps/social/api/internal/svc"
	"easy_chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 好友申请列表
func FriendPutInListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendPutInListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := Friends.NewFriendPutInListLogic(r.Context(), svcCtx)
		resp, err := l.FriendPutInList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)

	}
}

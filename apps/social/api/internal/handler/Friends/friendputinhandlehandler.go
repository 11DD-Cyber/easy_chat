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

// 好友申请处理
func FriendPutInHandleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendPutInHandleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := Friends.NewFriendPutInHandleLogic(r.Context(), svcCtx)
		resp, err := l.FriendPutInHandle(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)

	}
}

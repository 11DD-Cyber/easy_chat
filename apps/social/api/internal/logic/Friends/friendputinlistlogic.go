// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package Friends

import (
	"context"

	"easy_chat/apps/social/api/internal/svc"
	"easy_chat/apps/social/api/internal/types"
	"easy_chat/apps/social/rpc/social"
	"easy_chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请列表
func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (resp *types.FriendPutInListResp, err error) {
	// todo: add your logic here and delete this line
	friendslist, err := l.svcCtx.Social.FriendPutInList(l.ctx, &social.FriendPutInListReq{
		UserId: ctxdata.GetUId(l.ctx),
	})
	if err != nil {
		return nil, err
	}
	var resplist = []*types.FriendRequests{}
	for _, v := range friendslist.List {
		item := &types.FriendRequests{
			Id:           v.Id,
			UserId:       v.UserId,
			ReqUid:       v.ReqUid,
			ReqMsg:       v.ReqMsg,
			HandleResult: int(v.HandleResult),
		}
		// 把新对象append到切片里（切片会自动扩容）
		resplist = append(resplist, item)
	}
	return &types.FriendPutInListResp{
		List: resplist,
	}, nil
}

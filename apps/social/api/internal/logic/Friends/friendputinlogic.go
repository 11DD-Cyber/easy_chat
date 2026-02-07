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

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请
func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {
	// todo: add your logic here and delete this line
	_, err = l.svcCtx.Social.FriendPutIn(l.ctx, &social.FriendPutInReq{
		UserId: ctxdata.GetUId(l.ctx),
		ReqUid: req.ReqId,
		ReqMsg: req.ReqMsg,
	})
	if err != nil {
		return nil, err
	}
	return &types.FriendPutInResp{}, nil
}

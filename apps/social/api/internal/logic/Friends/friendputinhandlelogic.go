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

type FriendPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请处理
func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(req *types.FriendPutInHandleReq) (resp *types.FriendPutInHandleResp, err error) {
	// todo: add your logic here and delete this line
	_, err = l.svcCtx.Social.FriendPutInHandle(l.ctx, &social.FriendPutInHandleReq{
		UserId:       ctxdata.GetUId(l.ctx),
		FriendReqId:  req.FriendReqId,
		HandleResult: req.HandleResult,
	})
	if err != nil {
		return nil, err
	}
	return &types.FriendPutInHandleResp{}, nil
}

package logic

import (
	"context"

	"easy_chat/apps/social/rpc/internal/svc"
	"easy_chat/apps/social/rpc/social"
	"easy_chat/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取好友申请列表
func (l *FriendPutInListLogic) FriendPutInList(in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	// todo: add your logic here and delete this line
	friendrequests, err := l.svcCtx.FriendRequestsModel.FindByTargetId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendrequests err %v req %v", err, in)

	}
	if len(friendrequests) == 0 {
		return nil, errors.WithStack(xerr.NewMsgErr("没有待处理的请求"))

	}
	var friendrequestslist = []*social.FriendRequests{}
	copier.Copy(&friendrequestslist, friendrequests)
	return &social.FriendPutInListResp{
		List: friendrequestslist,
	}, nil
}

package logic

import (
	"context"

	"easy_chat/apps/social/rpc/internal/svc"
	"easy_chat/apps/social/rpc/social"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutOutListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutOutListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutOutListLogic {
	return &FriendPutOutListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutOutListLogic) FriendPutOutList(in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	reqs, err := l.svcCtx.FriendRequestsModel.FindBySenderId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(err, "find friend requests by sender err, req=%v", in)
	}

	list := make([]*social.FriendRequests, 0, len(reqs))
	if len(reqs) > 0 {
		if err := copier.Copy(&list, reqs); err != nil {
			return nil, err
		}
	}

	return &social.FriendPutInListResp{List: list}, nil
}

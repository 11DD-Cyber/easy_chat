package logic

import (
	"context"
	"log"

	"easy_chat/apps/social/rpc/internal/svc"
	"easy_chat/apps/social/rpc/social"
	"easy_chat/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取我的好友列表
func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	// todo: add your logic here and delete this line
	friends, err := l.svcCtx.FriendsModel.FindByUserId(l.ctx, in.UserId)
	log.Printf("得到的结果friends是%v", friends)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends by userid err %v req %v", err, in)
	}
	var friendsresp = []*social.Friends{}
	copier.Copy(&friendsresp, friends)
	return &social.FriendListResp{
		List: friendsresp,
	}, nil
}

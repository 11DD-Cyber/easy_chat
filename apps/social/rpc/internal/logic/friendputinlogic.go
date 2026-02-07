package logic

import (
	"context"

	"easy_chat/apps/social/rpc/internal/svc"
	"easy_chat/apps/social/rpc/social"
	"easy_chat/apps/social/socialmodels"
	"easy_chat/pkg/constants"
	"easy_chat/pkg/wuid"
	"easy_chat/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 好友业务 ：请求好友、通过或拒绝申请、好友列表
func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// todo: add your logic here and delete this line
	//申请人是否与目标是好友关系
	friends, err := l.svcCtx.FriendsModel.FindOneByUserIdFriendUid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends by uid and fid err %v req %v", err, in)
	}
	if friends != nil {
		return &social.FriendPutInResp{}, err
	}
	//是否已经申请
	friendReqs, err := l.svcCtx.FriendRequestsModel.FindOneByUserIdReqUid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by uid and rid err %v req %v", err, in)
	}
	if friendReqs != nil {
		return &social.FriendPutInResp{}, xerr.NewMsgErr("已申请好友，不可重复申请")
	}
	//创建申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		Id:           wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		UserId:       in.UserId,
		ReqUid:       in.ReqUid,
		ReqMsg:       in.ReqMsg,
		HandleResult: constants.NoHandlerResult,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friendRequest err %v req %v", err, in)

	}

	return &social.FriendPutInResp{}, nil
}

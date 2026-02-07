package logic

import (
	"context"
	"log"

	"easy_chat/apps/social/rpc/internal/svc"
	"easy_chat/apps/social/rpc/social"
	"easy_chat/apps/social/socialmodels"
	"easy_chat/pkg/constants"
	"easy_chat/pkg/wuid"
	"easy_chat/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var (
	ErrFriendReqBeforeRefuse = xerr.New(xerr.FRIEND_REQ_ALREADY_REFUSE, xerr.ErrMsg(xerr.FRIEND_REQ_ALREADY_REFUSE))
	ErrFriendReqBeforePass   = xerr.New(xerr.FRIEND_REQ_ALREADY_PASS, xerr.ErrMsg(xerr.FRIEND_REQ_ALREADY_PASS))
)

// 处理好友申请
func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// todo: add your logic here and delete this line
	//获取好友申请记录
	log.Printf("【关键验证】接收到的 friendReqId 原始值：%s（类型：%T）", in.FriendReqId, in.FriendReqId)
	friendReq, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, in.FriendReqId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendRequest by friendReqId err %v req %v", err, in)
	}
	//验证是否有处理
	if friendReq.HandleResult == int64(constants.RefuseHandlerResult) {
		return nil, errors.WithStack(ErrFriendReqBeforeRefuse)
	} else if friendReq.HandleResult == int64(constants.PassHandlerResult) {
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	}

	//修改申请结果

	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.FriendRequestsModel.UpdateHandleResult(ctx, session, in.FriendReqId, int64(in.HandleResult)); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend request err %v req %v", err, in.FriendReqId)
		}
		if in.HandleResult != constants.PassHandlerResult {
			return nil
		}

		friends := []*socialmodels.Friends{
			{Id: wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
				UserId:    friendReq.UserId,
				FriendUid: friendReq.ReqUid,
			},
			{
				Id:        wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
				UserId:    friendReq.ReqUid,
				FriendUid: friendReq.UserId,
			},
		}
		_, err = l.svcCtx.FriendsModel.Inserts(ctx, session, friends...)
		if err != nil {
			log.Printf("【插入好友失败】friends: %+v, err: %v", friends, err)
			return errors.Wrapf(xerr.NewDBErr(), "friends inserts err %v req %v", err, friends)
		}
		return nil
	})
	return &social.FriendPutInHandleResp{}, nil
}

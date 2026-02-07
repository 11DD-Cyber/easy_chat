// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package Friends

import (
	"context"

	"easy_chat/apps/social/api/internal/svc"
	"easy_chat/apps/social/api/internal/types"
	"easy_chat/apps/social/rpc/social"
	"easy_chat/apps/user/rpc/user"
	"easy_chat/apps/user/rpc/userclient"
	"easy_chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友列表
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// todo: add your logic here and delete this line
	userid := ctxdata.GetUId(l.ctx)
	friends, err := l.svcCtx.Social.FriendList(l.ctx, &social.FriendListReq{
		UserId: userid,
	})
	if err != nil {
		return nil, err
	}
	if len(friends.List) == 0 {
		return &types.FriendListResp{}, nil

	}
	//好友列表，根据好友id获取好友信息
	uids := make([]string, 0, len(friends.List))
	for _, i := range friends.List {
		uids = append(uids, i.FriendUid)
	}
	//根据uids查询用户信息
	users, err := l.svcCtx.UserRpc.FindUser(l.ctx, &user.FindUserReq{
		Ids: uids,
	})
	if err != nil {
		return nil, err
	}
	userRecords := make(map[string]*userclient.UserEntity, len(users.User))
	for i, _ := range users.User {
		userRecords[users.User[i].Id] = users.User[i]
	}
	respList := make([]*types.Friends, 0, len(friends.List))
	for _, v := range friends.List {
		friend := &types.Friends{
			Id:        v.UserId,
			FriendUid: v.FriendUid,
			Remark:    v.Remark,
		}
		if u, ok := userRecords[v.FriendUid]; ok {
			friend.Nickname = u.Nickname
			friend.Avatar = u.Avatar
		}
		respList = append(respList, friend)
	}
	return &types.FriendListResp{
		List: respList,
	}, nil

}

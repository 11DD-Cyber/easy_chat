package Friends

import (
	"context"

	"easy_chat/apps/social/api/internal/svc"
	"easy_chat/apps/social/api/internal/types"
	"easy_chat/apps/social/rpc/social"
	"easy_chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutOutListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutOutListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutOutListLogic {
	return &FriendPutOutListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutOutListLogic) FriendPutOutList(req *types.FriendPutOutListReq) (*types.FriendPutOutListResp, error) {
	uid := ctxdata.GetUId(l.ctx)

	rpcResp, err := l.svcCtx.Social.FriendPutOutList(l.ctx, &social.FriendPutInListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	if rpcResp == nil || len(rpcResp.List) == 0 {
		return &types.FriendPutOutListResp{List: []*types.FriendRequests{}}, nil
	}

	list := make([]*types.FriendRequests, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, &types.FriendRequests{
			Id:           item.Id,
			UserId:       item.UserId,
			ReqUid:       item.ReqUid,
			ReqMsg:       item.ReqMsg,
			HandleResult: int(item.HandleResult),
		})
	}

	return &types.FriendPutOutListResp{List: list}, nil
}

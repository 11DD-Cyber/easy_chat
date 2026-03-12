package user

import (
	"context"

	"easy_chat/apps/user/api/internal/svc"
	"easy_chat/apps/user/api/internal/types"
	"easy_chat/apps/user/rpc/user"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindUserLogic) FindUser(req *types.FindUserReq) (*types.FindUserResp, error) {
	rpcReq := &user.FindUserReq{
		Phone: req.Phone,
		Name:  req.Name,
		Ids:   req.Ids,
	}

	rpcResp, err := l.svcCtx.User.FindUser(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	if len(rpcResp.User) == 0 {
		return &types.FindUserResp{List: []types.User{}}, nil
	}

	list := make([]types.User, len(rpcResp.User))
	for i, u := range rpcResp.User {
		item := types.User{}
		_ = copier.Copy(&item, u)
		item.Mobile = u.Phone
		list[i] = item
	}

	return &types.FindUserResp{List: list}, nil
}

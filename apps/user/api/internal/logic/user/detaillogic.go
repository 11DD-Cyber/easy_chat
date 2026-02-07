// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"easy_chat/apps/user/api/internal/svc"
	"easy_chat/apps/user/api/internal/types"
	"easy_chat/apps/user/rpc/user"
	"easy_chat/pkg/ctxdata"
	"easy_chat/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUId(l.ctx)
	if uid == "" {
		return nil, xerr.New(xerr.TOKEN_EXPIRE_ERROR, "token为空或过期")
	}
	userInfoResp, err := l.svcCtx.User.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}
	var res types.User
	copier.Copy(&res, userInfoResp.User)
	return &types.UserInfoResp{Info: res}, nil
}

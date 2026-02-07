// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"easy_chat/apps/social/api/internal/config"
	"easy_chat/apps/social/rpc/socialclient"
	"easy_chat/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userclient.User
	Social  socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social:  socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
}

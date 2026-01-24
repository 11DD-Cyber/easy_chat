// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"easy_chat/apps/user/api/internal/config"
	"easy_chat/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	User   userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}

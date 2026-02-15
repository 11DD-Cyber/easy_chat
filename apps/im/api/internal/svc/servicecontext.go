// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"easy_chat/apps/im/api/internal/config"
	"easy_chat/apps/im/rpc/imclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	ImRpc  imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		ImRpc:  imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}

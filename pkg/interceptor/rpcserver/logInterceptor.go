package rpcserver

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}
	// 2. 日志：记录“哪个RPC方法+什么错误+请求上下文”，方便排查
	logx.WithContext(ctx).Errorf("【RPC SRV ERR】%v", err)
	// 3. 剥开错误包装，拿到最原始的错误（你的xerr创建的错误）
	rootErr := errors.Cause(err)
	// 4. 从原始错误中获取你定义的错误码和信息
	return resp, rootErr
}

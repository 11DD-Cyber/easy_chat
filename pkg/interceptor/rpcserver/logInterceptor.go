package rpcserver

import (
	"context"
	"strconv"

	"easy_chat/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	rootCode := xerr.Code(rootErr)
	rootMsg := xerr.Message(rootErr)
	if rootMsg == "" && rootErr != nil {
		rootMsg = rootErr.Error()
	}
	st := status.New(statusCodeFromRoot(rootCode), rootMsg)
	stWithDetails, detailErr := st.WithDetails(&errdetails.ErrorInfo{Reason: strconv.Itoa(rootCode)})
	if detailErr != nil {
		return resp, st.Err()
	}
	return resp, stWithDetails.Err()
}

func statusCodeFromRoot(code int) codes.Code {
	switch code {
	case xerr.REQUEST_PARAM_ERROR:
		return codes.InvalidArgument
	case xerr.TOKEN_EXPIRE_ERROR:
		return codes.Unauthenticated
	case xerr.DB_ERROR:
		return codes.Internal
	default:
		return codes.Unknown
	}
}

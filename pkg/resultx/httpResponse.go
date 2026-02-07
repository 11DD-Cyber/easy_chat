package resultx

import (
	"context"
	"easy_chat/pkg/xerr"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	grpcstatus "google.golang.org/grpc/status"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}
func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}
func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errcode := xerr.Code(err)
		errmsg := xerr.Message(err)
		if st, ok := grpcstatus.FromError(err); ok {
			for _, detail := range st.Details() {
				if info, ok := detail.(*errdetails.ErrorInfo); ok {
					if code, codeErr := strconv.Atoi(info.Reason); codeErr == nil {
						errcode = code
						break
					}
				}
			}
			if st.Message() != "" {
				errmsg = st.Message()
			}
		}
		if errmsg == "" {
			errmsg = xerr.ErrMsg(errcode)
		}
		status := statusFromErrCode(errcode)
		logx.WithContext(ctx).Errorf("【%s】 code=%d msg=%s err=%v", name, errcode, errmsg, err)
		return status, Fail(errcode, errmsg)
	}
}

func statusFromErrCode(code int) int {
	switch code {
	case xerr.REQUEST_PARAM_ERROR:
		return http.StatusBadRequest
	case xerr.TOKEN_EXPIRE_ERROR:
		return http.StatusUnauthorized
	case xerr.DB_ERROR, xerr.SERVER_COMMON_ERROR:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

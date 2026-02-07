package xerr

import (
	stdErrors "errors"

	xerrors "github.com/zeromicro/x/errors"
)

func New(code int, msg string) error {
	return xerrors.New(code, msg)
}
func NewMsgErr(msg string) error {
	return xerrors.New(SERVER_COMMON_ERROR, msg)

}
func NewCodeErr(code int) error {
	return xerrors.New(code, ErrMsg(code))
}
func NewInternalErr() error {
	return xerrors.New(SERVER_COMMON_ERROR, ErrMsg(SERVER_COMMON_ERROR))
}
func NewDBErr() error {
	return xerrors.New(DB_ERROR, ErrMsg(DB_ERROR))
}
func NewReqParamErr() error {
	return xerrors.New(REQUEST_PARAM_ERROR, ErrMsg(REQUEST_PARAM_ERROR))
}

func getCodeMsg(err error) *xerrors.CodeMsg {
	if err == nil {
		return nil
	}

	var cm *xerrors.CodeMsg
	if stdErrors.As(err, &cm) {
		return cm
	}
	return nil
}

func Code(err error) int {
	if cm := getCodeMsg(err); cm != nil {
		return cm.Code
	}
	return SERVER_COMMON_ERROR
}

func Message(err error) string {
	if cm := getCodeMsg(err); cm != nil {
		return cm.Msg
	}
	if err != nil {
		return err.Error()
	}
	return ErrMsg(SERVER_COMMON_ERROR)
}

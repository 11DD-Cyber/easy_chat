package xerr

var condeText = map[int]string{
	SERVER_COMMON_ERROR:       "服务异常，请稍后处理",
	REQUEST_PARAM_ERROR:       "参数不正确",
	TOKEN_EXPIRE_ERROR:        "token失效，请重新登录",
	DB_ERROR:                  "数据库繁忙，请稍后再试",
	FRIEND_REQ_ALREADY_REFUSE: "该好友申请已被拒绝，无法重复处理",
	FRIEND_REQ_ALREADY_PASS:   "该好友申请已被同意，无法重复处理",
	FRIEND_REQ_NOT_EXIST:      "好友申请记录不存在，请检查申请ID",
}

func ErrMsg(errcode int) string {
	if msg, ok := condeText[errcode]; ok {
		return msg
	}
	return condeText[SERVER_COMMON_ERROR]
}

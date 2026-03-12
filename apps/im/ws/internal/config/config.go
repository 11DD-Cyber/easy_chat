package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	// 宓屽叆ServiceConf锛氬鐢╣o-zero鐨勯€氱敤鏈嶅姟閰嶇疆
	// 鍖呭惈鏈嶅姟鍚嶏紙Name锛夈€佹棩蹇椼€佺洃鎺с€侀摼璺拷韪€佽繍琛屾ā寮忥紙dev/prod锛夌瓑閫氱敤閰嶇疆椤?
	service.ServiceConf
	ListenOn string
	JwtAuth  struct {
		AccessSecret string
	}
	Mongo struct {
		Url string
		Db  string
	}
	MsgChatTransfer struct {
		Addrs []string
		Topic string
	}
	Redisx redis.RedisConf
}

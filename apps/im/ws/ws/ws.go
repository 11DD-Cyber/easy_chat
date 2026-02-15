package ws

type (
	Push struct {
		ConversationId string `mapstructure:"conversationId" json:"conversationId"`
		SendId         string `mapstructure:"sendId" json:"sendId"`
		RecvId         string `mapstructure:"recvId" json:"recvId"`
		Msg            `mapstructure:"msg" json:"msg"`
		SendTime       int64 `mapstructure:"sendTime" json:"sendTime"`
	}
)

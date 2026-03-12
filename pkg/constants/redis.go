package constants

const (
	REDIS_SYSTEM_ROOT_TOKEN = "system:root:token"

	RedisKeyOnlineUserFmt  = "im:online:%s"
	RedisKeyOfflineMsgFmt  = "im:offline:%s"
	RedisKeyUnreadFmt      = "im:unread:%s:%s"
	RedisOnlineTTLSeconds  = 120
	RedisOfflineTTLSeconds = 7 * 24 * 60 * 60
	RedisOfflineSyncLimit  = 200
)

package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"easy_chat/apps/im/ws/ws"
	"easy_chat/pkg/constants"
)

type OfflineMessageEntry struct {
	Member string
	Chat   *ws.Chat
}

func onlineKey(uid string) string {
	return fmt.Sprintf(constants.RedisKeyOnlineUserFmt, uid)
}

func offlineKey(uid string) string {
	return fmt.Sprintf(constants.RedisKeyOfflineMsgFmt, uid)
}

func unreadKey(uid, conversationId string) string {
	return fmt.Sprintf(constants.RedisKeyUnreadFmt, uid, conversationId)
}

func ttlSeconds(seconds int) int {
	if seconds <= 0 {
		return 0
	}
	return seconds
}

func (s *ServiceContext) MarkUserOnline(ctx context.Context, uid string) error {
	if s.Redis == nil {
		return nil
	}
	ttl := ttlSeconds(constants.RedisOnlineTTLSeconds)
	return s.Redis.SetexCtx(ctx, onlineKey(uid), fmt.Sprintf("%d", time.Now().Unix()), ttl)
}

func (s *ServiceContext) MarkUserOffline(ctx context.Context, uid string) error {
	if s.Redis == nil {
		return nil
	}
	_, err := s.Redis.DelCtx(ctx, onlineKey(uid))
	return err
}

func (s *ServiceContext) AppendOfflineMessage(ctx context.Context, uid string, chat *ws.Chat) error {
	if s.Redis == nil || chat == nil {
		return nil
	}
	payload, err := json.Marshal(chat)
	if err != nil {
		return err
	}
	score := float64(chat.SendTime)
	if score == 0 {
		score = float64(time.Now().UnixMilli())
	}
	if _, err = s.Redis.ZaddFloatCtx(ctx, offlineKey(uid), score, string(payload)); err != nil {
		return err
	}
	return s.Redis.ExpireCtx(ctx, offlineKey(uid), constants.RedisOfflineTTLSeconds)
}

func (s *ServiceContext) SaveOfflineMessage(ctx context.Context, uid string, chat *ws.Chat) error {
	if err := s.AppendOfflineMessage(ctx, uid, chat); err != nil {
		return err
	}
	return s.IncreaseUnread(ctx, uid, chat.ConversationId)
}

func (s *ServiceContext) FetchOfflineMessages(ctx context.Context, uid string, limit int64) ([]OfflineMessageEntry, error) {
	if s.Redis == nil {
		return nil, nil
	}
	if limit <= 0 || limit > int64(constants.RedisOfflineSyncLimit) {
		limit = int64(constants.RedisOfflineSyncLimit)
	}
	end := limit - 1
	values, err := s.Redis.ZrangeCtx(ctx, offlineKey(uid), 0, end)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}
	result := make([]OfflineMessageEntry, 0, len(values))
	for _, val := range values {
		var chat ws.Chat
		if err := json.Unmarshal([]byte(val), &chat); err != nil {
			continue
		}
		result = append(result, OfflineMessageEntry{
			Member: val,
			Chat:   &chat,
		})
	}
	return result, nil
}

func (s *ServiceContext) RemoveOfflineMessages(ctx context.Context, uid string, members ...string) error {
	if s.Redis == nil || len(members) == 0 {
		return nil
	}
	args := make([]any, 0, len(members))
	for _, m := range members {
		args = append(args, m)
	}
	_, err := s.Redis.ZremCtx(ctx, offlineKey(uid), args...)
	return err
}

func (s *ServiceContext) IncreaseUnread(ctx context.Context, uid, conversationId string) error {
	if s.Redis == nil || conversationId == "" {
		return nil
	}
	key := unreadKey(uid, conversationId)
	if _, err := s.Redis.IncrCtx(ctx, key); err != nil {
		return err
	}
	return s.Redis.ExpireCtx(ctx, key, constants.RedisOfflineTTLSeconds)
}

func (s *ServiceContext) ClearUnread(ctx context.Context, uid, conversationId string) error {
	if s.Redis == nil || conversationId == "" {
		return nil
	}
	_, err := s.Redis.DelCtx(ctx, unreadKey(uid, conversationId))
	return err
}

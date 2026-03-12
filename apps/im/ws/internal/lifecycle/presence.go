package lifecycle

import (
	"context"

	"easy_chat/apps/im/ws/internal/svc"
	"easy_chat/apps/im/ws/websocket"
)

type Presence struct {
	svc *svc.ServiceContext
}

func NewPresence(svcCtx *svc.ServiceContext) *Presence {
	return &Presence{svc: svcCtx}
}

func (p *Presence) OnConnect(srv *websocket.Server, conn *websocket.Conn, uid string) {
	if p.svc == nil {
		return
	}
	ctx := context.Background()
	if err := p.svc.MarkUserOnline(ctx, uid); err != nil {
		srv.Logger.Errorf("mark user %s online err: %v", uid, err)
	}
	entries, err := p.svc.FetchOfflineMessages(ctx, uid, 0)
	if err != nil {
		srv.Logger.Errorf("fetch offline messages err uid %s: %v", uid, err)
		return
	}
	if len(entries) == 0 {
		return
	}
	members := make([]string, 0, len(entries))
	cleared := map[string]struct{}{}
	for _, entry := range entries {
		if entry.Chat == nil {
			continue
		}
		if err := srv.Send(websocket.NewMessage(entry.Chat.SendId, entry.Chat), conn.Conn); err != nil {
			srv.Logger.Errorf("send offline message err uid %s: %v", uid, err)
			continue
		}
		members = append(members, entry.Member)
		cleared[entry.Chat.ConversationId] = struct{}{}
	}
	if len(members) > 0 {
		if err := p.svc.RemoveOfflineMessages(ctx, uid, members...); err != nil {
			srv.Logger.Errorf("remove offline cache err uid %s: %v", uid, err)
		}
		for conversationId := range cleared {
			if err := p.svc.ClearUnread(ctx, uid, conversationId); err != nil {
				srv.Logger.Errorf("clear unread err uid %s cid %s: %v", uid, conversationId, err)
			}
		}
	}
}

func (p *Presence) OnDisconnect(srv *websocket.Server, conn *websocket.Conn, uid string) {
	if p.svc == nil {
		return
	}
	ctx := context.Background()
	if err := p.svc.MarkUserOffline(ctx, uid); err != nil {
		srv.Logger.Errorf("mark user %s offline err: %v", uid, err)
	}
}

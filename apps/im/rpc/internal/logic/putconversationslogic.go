package logic

import (
	"context"

	"easy_chat/apps/im/models"
	"easy_chat/apps/im/rpc/im"
	"easy_chat/apps/im/rpc/internal/svc"
	"easy_chat/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新会话
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	// todo: add your logic here and delete this line

	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find conversations by userId err %v req %v", err, in.UserId)

	}
	if conversations.ConversationList == nil {
		//没有建立会话
		conversations.ConversationList = make(map[string]*models.Conversation)
	}
	for k, i := range in.ConversationList {
		var oldTotal int32
		if conversations.ConversationList[k] != nil {
			oldTotal = int32(conversations.ConversationList[k].Total)
		}
		conversations.ConversationList[k] = &models.Conversation{
			ConversationId: i.ConversationId,
			ChatType:       0,
			IsShow:         i.IsShow,
			//更新最新的已读总记录
			Total: i.ToRead + oldTotal,
			Seq:   i.Seq,
		}
	}
	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "update conversation err %v req %v", err, conversations)
	}
	return &im.PutConversationsResp{}, err
}

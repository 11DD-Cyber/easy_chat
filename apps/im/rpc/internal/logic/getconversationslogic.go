package logic

import (
	"context"

	"easy_chat/apps/im/models"
	"easy_chat/apps/im/rpc/im"
	"easy_chat/apps/im/rpc/internal/svc"
	"easy_chat/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	// todo: add your logic here and delete this line
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		if err == models.ErrNotFound {
			return &im.GetConversationsResp{}, nil
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find conversation by userId err %v req %v", err, in)

	}
	var res im.GetConversationsResp
	copier.Copy(&res, &data)
	ids := make([]string, 0, len(data.ConversationList))
	//获取所有会话id
	for _, conversation := range data.ConversationList {
		ids = append(ids, conversation.ConversationId)
	}
	//统计会话的消息情况
	list, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find conversation by conversationIds err %v req %v", err, ids)
	}
	for _, conversation := range list {
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			continue
		}
		total := res.ConversationList[conversation.ConversationId].Total
		if total < int32(conversation.Total) {
			//有新的消息
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			//待读消息量
			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.ToRead)
			// 标记有新消息，前端需要显示提示
			res.ConversationList[conversation.ConversationId].IsShow = true
		}
	}
	return &res, nil
}

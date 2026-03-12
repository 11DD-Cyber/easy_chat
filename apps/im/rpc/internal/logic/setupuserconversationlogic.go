package logic

import (
	"context"

	"easy_chat/apps/im/models"
	"easy_chat/apps/im/rpc/im"
	"easy_chat/apps/im/rpc/internal/svc"
	"easy_chat/pkg/wuid"
	"easy_chat/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	var res im.SetUpUserConversationResp
	switch in.ChatType {
	case 0:
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			if err == models.ErrNotFound {
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &models.Conversation{
					ConversationId: conversationId,
					Msg:            &models.ChatLog{},
					ChatType:       0,
				})
				if err != nil {
					return nil, errors.Wrapf(xerr.NewDBErr(), "create conversation err %v ", err)
				}
			} else {
				return nil, errors.Wrapf(xerr.NewDBErr(), "find conversation by conversationId err %v req %v", err, conversationId)
			}
		}
		err = l.setUpUserConversation(conversationId, in.SendId, true)
		if err != nil {
			return &res, errors.Wrapf(xerr.NewDBErr(), "set up user single conversation err %v", err)
		}
	case 1:
		return nil, nil
	}
	return &res, nil
}

func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId string, userId string, isShow bool) error {
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	isNew := false
	if err != nil {
		if err == models.ErrNotFound {
			conversations = &models.Conversations{
				ID:               bson.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*models.Conversation),
			}
			isNew = true
		} else {
			return err
		}
	}
	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}
	conversations.ConversationList[conversationId] = &models.Conversation{
		ConversationId: conversationId,
		ChatType:       0,
		IsShow:         isShow,
	}
	if isNew {
		return l.svcCtx.ConversationsModel.Insert(l.ctx, conversations)
	}
	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	return err
}

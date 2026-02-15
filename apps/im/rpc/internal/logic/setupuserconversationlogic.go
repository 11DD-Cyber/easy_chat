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

// 建立会话
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// todo: add your logic here and delete this line

	var res im.SetUpUserConversationResp
	switch in.ChatType {
	case 0:
		//建立私聊的关系
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		//建立两者的会话
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
	//发送者
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if err == models.ErrNotFound {
			conversations = &models.Conversations{
				ID:               bson.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*models.Conversation),
			}
		} else {
			return err
		}
	}
	//更新会话记录
	if _, ok := conversations.ConversationList[conversationId]; ok {
		//存在
		return nil
	}
	//需要建立
	conversations.ConversationList[conversationId] = &models.Conversation{
		ConversationId: conversationId,
		ChatType:       0,
		IsShow:         isShow,
	}
	//存在即更新，不存在则修改
	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	return err
}

package logic

import (
	"context"
	"database/sql"
	"easy_chat/apps/user/models"
	"easy_chat/apps/user/rpc/internal/svc"
	"easy_chat/apps/user/rpc/user"
	"easy_chat/pkg/ctxdata"
	"easy_chat/pkg/encrypt"
	"easy_chat/pkg/wuid"
	"easy_chat/pkg/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrPhoneIsRegister = xerr.NewMsgErr("\u624b\u673a\u53f7\u5df2\u6ce8\u518c")

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// todo: add your logic here and delete this line
	//1.楠岃瘉鐢ㄦ埛鏄惁娉ㄥ唽锛屾牴鎹墜鏈哄彿楠岃瘉
	userEntity, err := l.svcCtx.UserModels.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}
	if userEntity != nil {
		return nil, ErrPhoneIsRegister
	}
	//瀹氫箟鐢ㄦ埛鏁版嵁
	userEntity = &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}
	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		userEntity.Password = sql.NullString{
			String: string(genPassword),
			Valid:  true,
		}
	}
	if _, err := l.svcCtx.UserModels.Insert(l.ctx, userEntity); err != nil {
		return nil, err
	}
	//鐢熸垚token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}

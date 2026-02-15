package svc

import (
	"easy_chat/apps/user/models"
	"easy_chat/apps/user/rpc/internal/config"
	"easy_chat/pkg/constants"
	"easy_chat/pkg/ctxdata"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	Redis      *redis.Redis
	UserModels models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redisx),
		UserModels: models.NewUsersModel(conn, c.Cache),
	}
}

func (s *ServiceContext) SetRootToken() error {
	systemToken, err := ctxdata.GetJwtToken(s.Config.Jwt.AccessSecret, time.Now().Unix(), 99999999, constants.SYSTEM_ROOT_UID)
	if err != nil {
		return err
	}
	return s.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN, systemToken)
}

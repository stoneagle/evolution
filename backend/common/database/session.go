package database

import (
	"evolution/backend/common/config"

	"github.com/gin-gonic/contrib/sessions"
)

func SessionByRedis(redisConf config.RedisConf) (store sessions.RedisStore, err error) {
	redisServer := redisConf.Host + ":" + redisConf.Port
	store, err = sessions.NewRedisStore(10, "tcp", redisServer, redisConf.Password, []byte("secret"))
	return
}

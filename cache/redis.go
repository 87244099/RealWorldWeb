package cache

import (
	"RealWorldWeb/config"
	"RealWorldWeb/logger"
	"context"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() {
	ctx := context.Background()
	log := logger.New(ctx)
	log.Infof("redisAddr=%s", config.GetRedisAddr())

	rdb = redis.NewClient(&redis.Options{
		Addr: config.GetRedisAddr(),
		//Password: "", // no password set
		//DB:       0,  // use default DB
	})
	ping := rdb.Ping(context.Background())
	err := ping.Err()
	if err != nil {
		log.WithError(err).Infof("redis ping failed")
		panic(err)
	}
}
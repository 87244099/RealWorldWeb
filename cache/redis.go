package cache

import (
	"RealWorldWeb/config"
	"RealWorldWeb/logger"
	"RealWorldWeb/models"
	"context"
	"encoding/json"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"time"
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

	Locker = redislock.New(rdb)
}

const (
	POPULAR_TAGS_KEY = "popular_tags"
	USER_PROFILE_KEY = "user_profile_"
)

func GetPopularTags(ctx context.Context) ([]string, error) {
	cmd := rdb.ZRange(ctx, POPULAR_TAGS_KEY, 0, 10)
	return cmd.Result()
}

func GetUserProfile(ctx context.Context, userName string) (*models.User, error) {
	log := logger.New(ctx)
	log.WithContext(ctx).Infof("cache GetUserProfile called")
	js, err := rdb.Get(ctx, USER_PROFILE_KEY+userName).Result()
	if err != nil {
		log.WithError(err).Infof("get userProfile cache failed; userName=%s", userName)
		return nil, err
	}
	user := &models.User{}
	err = json.Unmarshal([]byte(js), user)
	if err != nil {
		log.WithError(err).Infof("parse failed; userName=%s", userName)
		return nil, err
	}
	return user, nil
}

func SetUserProfile(ctx context.Context, userName string, user *models.User, ttl int) error {
	log := logger.New(ctx)
	log.WithContext(ctx).Infof("cache SetUserProfile called")
	js, err := json.Marshal(user)
	if err != nil {
		log.WithError(err).Infof("stringify failed; userName=%s", userName)
		return err
	}
	err = rdb.Set(ctx, USER_PROFILE_KEY+userName, string(js), time.Duration(ttl)*time.Second).Err()
	if err != nil {
		log.WithError(err).Infof("set userProfile cache failed; userName=%s", userName)
		return err
	}
	return nil
}

func DeleteUserProfile(ctx context.Context, userName string) error {
	log := logger.New(ctx)
	log.WithContext(ctx).Infof("cache DeleteUserProfile called")
	err := rdb.Del(ctx, USER_PROFILE_KEY+userName).Err()
	if err != nil {
		log.WithError(err).Infof("delete userProfile key failed; userName=%s", userName)
	}
	return err
}

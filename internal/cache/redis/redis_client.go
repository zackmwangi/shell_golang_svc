package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type (
	RedisClientConfig struct {
		Address     string        // "localhost:6379"
		Username    string        // ""
		Password    string        // ""
		DbNum       string        // 0
		DialTimeout time.Duration //100 * time.Millisecond
		ReadTimeout time.Duration //100 * time.Millisecond
		Logger      *zap.Logger
	}
)

var Ctx = context.Background()

func NewRedisClient(redisConfig RedisClientConfig) *redis.Client {

	dbNum, err := strconv.Atoi(redisConfig.DbNum)

	if err != nil {
		redisConfig.Logger.Sugar().Fatal("could not connect to redis db num %s, error: %v", dbNum, err)
		//return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:        redisConfig.Address,
		Username:    redisConfig.Username,
		Password:    redisConfig.Password,
		DB:          dbNum,
		DialTimeout: redisConfig.DialTimeout * time.Millisecond,
		ReadTimeout: redisConfig.ReadTimeout * time.Millisecond,
	})

	if err := client.Ping(Ctx).Err(); err != nil {
		redisConfig.Logger.Sugar().Fatalf("error connecting to redis at "+redisConfig.Address+" with message=[%v]", err)
		//return nil, err
	}

	redisConfig.Logger.Sugar().Infof(" successfully connected to redis host %s at DB %s", redisConfig.Address, redisConfig.DbNum)
	return client
}

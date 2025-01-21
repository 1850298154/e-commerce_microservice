package redis

import (
	"context"

	"2501YTC/app/order/conf"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client // RedisClient Redis客户端

// Init 初始化Redis
func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.GetConf().Redis.Address,
		Username: conf.GetConf().Redis.Username,
		Password: conf.GetConf().Redis.Password,
		DB:       conf.GetConf().Redis.DB,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	klog.Infof("Redis initial successfully, addr: %s, username: %s, password: %s, DB: %s",
		conf.GetConf().Redis.Address,
		conf.GetConf().Redis.Username,
		conf.GetConf().Redis.Password,
		conf.GetConf().Redis.DB)
}

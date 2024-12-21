package db

import "github.com/go-redis/redis/v8"

var redisClient *redis.Client

func InitRedis() {
	// 创建 Redis 客户端
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址，默认端口是 6379
		Password: "",               // 密码，默认空
		DB:       0,                // 默认数据库是 0
	})
}

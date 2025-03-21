package redis

import (
	"context"
	"fmt"
	"myApp/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

// 定义常量，用于判断缓存是否存在
var Nil = redis.Nil

// InitRedis 初始化Redis客户端
func InitRedis() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Conf.Redis.Host, config.Conf.Redis.Port),
			Password: "", // 如果有密码，可以在配置中添加
			DB:       0,  // 使用默认DB
		})

		// 测试连接
		_, err := redisClient.Ping(ctx).Result()
		if err != nil {
			panic(fmt.Sprintf("Redis连接失败: %v", err))
		}
		fmt.Println("Redis连接成功")
	}
	return redisClient
}

// GetRedisClient 获取Redis客户端实例
func GetRedisClient() *redis.Client {
	if redisClient == nil {
		redisClient = InitRedis()
	}
	return redisClient
}

// Set 设置缓存
func Set(key string, value interface{}, expiration time.Duration) error {
	return GetRedisClient().Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func Get(key string) (string, error) {
	return GetRedisClient().Get(ctx, key).Result()
}

// Delete 删除缓存
func Delete(key string) error {
	return GetRedisClient().Del(ctx, key).Err()
}

// DeleteByPattern 根据模式删除缓存
func DeleteByPattern(pattern string) error {
	keys, err := GetRedisClient().Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return GetRedisClient().Del(ctx, keys...).Err()
	}

	return nil
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	result, err := GetRedisClient().Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return GetRedisClient().Expire(ctx, key, expiration).Err()
}

// Incr 自增
func Incr(key string) (int64, error) {
	return GetRedisClient().Incr(ctx, key).Result()
}
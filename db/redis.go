package db

import (
	"fmt"
	"game_server/core/base"
	"time"

	"game_server/core/logger"

	"github.com/go-redis/redis" // 实现了redis连接池
)

type RedisClient struct {
	client *redis.Client
}

var RedisMgr = NewRedisClient()

//连接redis，并返回连接实例
func NewRedisClient() *RedisClient {
	rm := new(RedisClient)
	rm.client = redis.NewClient(&redis.Options{
		Addr:        base.Setting.Redis.Host + ":" + base.Setting.Redis.Port, // Redis地址
		Password:    base.Setting.Redis.Password,                             // Redis账号
		DB:          base.Setting.Redis.DbName,                               // Redis库
		PoolSize:    base.Setting.Redis.PoolSize,                             // Redis连接池大小
		MaxRetries:  3,                                                       // 最大重试次数
		IdleTimeout: 10 * time.Second,                                        // 空闲链接超时时间
	})
	pong, err := rm.client.Ping().Result()
	if err == redis.Nil {
		logger.Info("Redis异常")
	} else if err != nil {
		logger.Info("失败:", err)
	} else {
		logger.Info(pong)
	}
	logger.Info("redis connect succc")
	return rm
}

// 向key的hash中添加元素field的值
func (this *RedisClient) HSet(key, field string, data interface{}) {
	err := this.client.HSet(key, field, data)
	if err != nil {
		logger.Errorf("Redis HSet Error:", err)
	}
}

// 批量向key的hash添加对应元素field的值
func (this *RedisClient) BatchHashSet(key string, fields map[string]interface{}) string {
	val, err := this.client.HMSet(key, fields).Result()
	if err != nil {
		logger.Errorf("Redis HMSet Error:", err)
	}
	return val
}

// 通过key获取hash的元素值
func (this *RedisClient) HGet(key, field string) string {
	result := ""
	val, err := this.client.HGet(key, field).Result()
	if err == redis.Nil {
		logger.Infof("Key Doesn't Exists:", field)
		return result
	} else if err != nil {
		logger.Infof("Redis HGet Error:", err)
		return result
	}
	return val
}

// 批量获取key的hash中对应多元素值
func (this *RedisClient) BatchHashGet(key string, fields ...string) map[string]interface{} {
	resMap := make(map[string]interface{})
	for _, field := range fields {
		var result interface{}
		val, err := this.client.HGet(key, fmt.Sprintf("%s", field)).Result()
		if err == redis.Nil {
			logger.Infof("Key Doesn't Exists:", field)
			resMap[field] = result
		} else if err != nil {
			logger.Infof("Redis HMGet Error:", err)
			resMap[field] = result
		}
		if val != "" {
			resMap[field] = val
		} else {
			resMap[field] = result
		}
	}
	return resMap
}

// 获取自增唯一ID
func (this *RedisClient) Incr(key string) int {
	val, err := this.client.Incr(key).Result()
	if err != nil {
		logger.Errorf("Redis Incr Error:", err)
	}
	return int(val)
}

// 添加集合数据
func (this *RedisClient) SetAdd(key, val string) {
	this.client.SAdd(key, val)
}

// 从集合中获取数据
func (this *RedisClient) SetGet(key string) []string {
	val, err := this.client.SMembers(key).Result()
	if err != nil {
		logger.Errorf("Redis SMembers Error:", err)
	}
	return val
}

// 设置数据
func (this *RedisClient) SetCode(key string, value interface{}) error {
	err := this.client.Set(key, value, 0).Err()
	if err != nil {
		logger.Errorf("set key Error:", key)
		return err
	}
	return nil
}

// 获取数据
func (this *RedisClient) GetCode(key string) interface{} {
	val, err := this.client.Get("key").Result()
	if err != nil {
		logger.Errorf("set key Error:", key)
		return ""
	}
	return val
}

//hash key 是否存在
func (this *RedisClient) HashExistsKey(key string, field string) bool {
	ret, err := this.client.HExists(key, field).Result()
	if err != nil {
		logger.Errorf("HashExistsKey Error:", key, field)
		return false
	}
	return ret
}

//key 是否存在
func (this *RedisClient) KeyExist(key string) bool {
	ret, err := this.client.Exists(key).Result()
	if err != nil {
		logger.Errorf("ExistsKey Error:", key)
		return false
	}
	return ret > 0
}

//返回指定key的hash表所有filed和value
func (this *RedisClient) HGetAll(key string) map[string]string {
	ret, err := this.client.HGetAll(key).Result()
	if err != nil {
		logger.Errorf("HGetAll Error:", key)
		return nil
	}
	return ret
}

// GetRedisClient 返回redis连接.
func (r *RedisClient) GetRedisClient() *redis.Client {
	return r.client
}

/**
设置过期时间
*/
func (this *RedisClient) Expire(key string, expiration int) bool {
	number := time.Duration(expiration)
	_, err := this.client.Expire(key, number*time.Second).Result()
	if err != nil {
		logger.Errorf("Expire Error:", key)
		return false
	}
	return true
}

//  获取hel 的数量
func (this *RedisClient) HLen(key string) int {
	val, err := this.client.HLen(key).Result()
	if err != nil {
		logger.Errorf("Redis HLen Error:", err)
	}
	return int(val)
}

// 向key的hash中添加元素field的值
func (this *RedisClient) HIncrBy(key string, field string, incr int64) int64 {
	ret, err := this.client.HIncrBy(key, field, 1).Result()
	if err != nil {
		logger.Errorf("Redis HIncrBy Error:", err)
	}
	return ret
}

// 获取Get数据
func (this *RedisClient) Get(key string) string {
	val, err := this.client.Get(key).Result()
	if err != nil {
		logger.Errorf("Redis Get Error:", err)
	}
	return val
}

package cache

import (
	"fmt"
	util "picture-oss-proxy/pkg/utils"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// 定义 MyRedis 结构体
type MyRedis struct {
	Client *redis.Client
}

type RateLimit struct {
	Key      string
	Expire   time.Duration
	Max      int64
	Decrease bool
}

// 使用单例模式进行封装
var once sync.Once

// RedisClient Redis缓存客户端单例
var RedisClient *MyRedis = new(MyRedis)

// 封装 redis 实例，提供获取
func GetInstance() *MyRedis {
	return RedisClient
}

func NewRedis(RedisAddr string, RedisDbName string, RedisPw string) *redis.Client {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	myRedis := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		//Password: "",
		DB: int(db),
	})
	_, err := myRedis.Ping().Result()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("redis已经连接")

	once.Do(func() {
		RedisClient.Client = myRedis
	})

	return myRedis
}

func (mr *MyRedis) Set(key string, value interface{}, ttl time.Duration) {
	mr.Client.Set(key, value, ttl)
}
func (mr MyRedis) Get(key string) *redis.StringCmd {
	return mr.Client.Get(key)
}

func (mr MyRedis) Incr(key string) interface{} {
	return mr.Client.Incr(key)
}

func (mr MyRedis) ZIncrBy(key string, increment float64, member string) interface{} {
	return mr.Client.ZIncrBy(key, increment, member)
}

func (mr *MyRedis) AllowIp(param *RateLimit) bool {

	key := param.Key
	expire := param.Expire
	max := param.Max
	decrease := param.Decrease

	//使用毫秒时间戳
	now := time.Now().UnixMilli()
	start := now - int64(expire)
	redisKey := "rateLimiter:" + key

	pipe := mr.Client.TxPipeline()
	pipe.ZRemRangeByScore(redisKey, strconv.Itoa(0), strconv.Itoa(int(start)))
	count := pipe.ZCard(redisKey)
	pipe.ZRange(redisKey, 0, 0)
	if decrease {
		pipe.ZAdd(redisKey, redis.Z{Score: util.Unwrap(now, 0), Member: now})
	}
	pipe.ZRange(redisKey, -max, -max)
	pipe.PExpire(redisKey, expire)
	cmd, err := pipe.Exec()
	if err != nil && err != redis.Nil {
		util.LogrusObj.Errorln("[cache]AllowIp err", err)
	}
	fmt.Println("cmd redis pipeline 命令打印", cmd)

	if count.Val() < max {
		return true
	}
	return false

}

// func main() {
// 	// 项目启动时初始化 redis
// 	NewRedis("localhost", "")
// 	fmt.Println("redis 连接成功")
// }

// package user
// func getUser() {
//     result := cache.GetInstance().Exist("user_001")
//     if !result {
//         fmt.Println("不存在该数据")
//     }
// }

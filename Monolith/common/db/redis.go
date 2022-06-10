package db

import (
	"io/ioutil"
	"mini-tiktok/common/xerr"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Redis Configure Struct
type RedisConfig struct {
	Host     string `yaml:"redis_host"`
	Port     string `yaml:"redis_port"`
	DBName   int    `yaml:"redis_dbname"`
	Password string `yaml:"redis_password"`
}

// Redis Data Access Object
type RedisConn struct {
	rdb *redis.Client
}

var redisConn *RedisConn
var redisOnce sync.Once

// Get redis configure from yaml file.
// If read or unmarshal failed, exit instantly.
func readRedisConf() *RedisConfig {
	config := new(RedisConfig)
	file, _ := ioutil.ReadFile("config.yaml")
	yaml.Unmarshal(file, config)
	return config
}

// connect to the redis
func InitRedis() *redis.Client {
	conf := readRedisConf()
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: "",
		DB:       conf.DBName,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		zap.L().Fatal("Redis init failed")
	}
	zap.L().Info("Redis init success")
	return rdb
}

// return redis dao with singleten pattern
func NewRedisDaoInstance() *RedisConn {
	redisOnce.Do(
		func() {
			rdb := InitRedis()
			redisConn = &RedisConn{rdb}
		})
	return redisConn
}

// store (token, userID) in redis
func (r *RedisConn) SetToken(token string, userID int64) error {
	// 0 means this key(token) never expire
	err := r.rdb.Set(token, userID, 0).Err()
	if err != nil {
		zap.L().Error("redis: set token userid failed")
		return xerr.ErrDatabase
	}
	return nil
}

// get token's value from redis
func (r *RedisConn) GetToken(token string) (int64, error) {
	val, err := r.rdb.Get(token).Result()
	if err == redis.Nil {
		return 0, xerr.ErrTokenNotFound
	} else if err != nil {
		zap.L().Error("redis: get token failed")
		return 0, xerr.ErrDatabase
	}
	id, _ := strconv.ParseInt(val, 10, 64)
	return id, nil
}

// once user logout use this method to remove token from redis.
func (r *RedisConn) ClearToken(token string) error {
	err := r.rdb.Del(token).Err()
	if err != nil {
		zap.L().Error("redis: del token failed")
		return xerr.ErrDatabase
	}
	return nil
}

// check token is expired or not.
func (r *RedisConn) IsTokenValid(token string) (bool, error) {
	val, err := r.rdb.TTL(token).Result()
	if err != nil {
		zap.L().Error("redis: ttl token failed")
		return false, xerr.ErrDatabase
	}
	// key 不存在返回 -2
	// key 存在但是没有关联超时时间返回 -1
	if val == -2 {
		return false, nil
	}
	return true, nil
}

func (r *RedisConn) AddToNameList(username string) error {
	err := r.rdb.SAdd("namelist", username).Err()
	if err != nil {
		zap.L().Error("redis: sadd namelist user failed")
		return xerr.ErrDatabase
	}
	return nil
}

func (r *RedisConn) IsUserNameExist(username string) (bool, error) {
	exist, err := r.rdb.SIsMember("namelist", username).Result()
	if err != nil {
		zap.L().Error("redis: sismember namelist username failed")
		return false, xerr.ErrDatabase
	}
	return exist, nil
}

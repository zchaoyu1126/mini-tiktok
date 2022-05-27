package db

import (
	"io/ioutil"
	"log"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
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
type RedisDao struct {
	rdb *redis.Client
}

var redisDao *RedisDao
var once2 sync.Once

// Get redis configure from yaml file.
// If read or unmarshal failed, exit instantly.
func getRedisConf() *RedisConfig {
	var c RedisConfig
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Fatal(err)
	}
	return &c
}

// connect to the redis
func InitRedis() *redis.Client {
	conf := getRedisConf()
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: "",
		DB:       conf.DBName,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("redis init success.")
	// read mysql user exist or not
	return rdb
}

// return redis dao with singleten pattern
func NewRedisDaoInstance() *RedisDao {
	once2.Do(
		func() {
			rdb := InitRedis()
			redisDao = &RedisDao{rdb}
		})
	return redisDao
}

// store (token, userID) in redis
func (r *RedisDao) SetToken(token string, userID int64) error {
	// 0 means this key(token) never expire
	err := r.rdb.Set(token, userID, 0).Err()
	return err
}

// get token's value from redis
func (r *RedisDao) GetToken(token string) (int64, error) {
	val, err := r.rdb.Get(token).Result()
	if err == redis.Nil {
		return -1, errors.New("repository key does not exist")
	} else if err != nil {
		return -1, errors.New("repository get failed")
	}
	id, _ := strconv.ParseInt(val, 10, 64)
	return id, nil
}

// once user logout use this method to remove token from redis.
func (r *RedisDao) ClearToken(token string) error {
	return r.rdb.Del(token).Err()
}

// check token is expired or not.
func (r *RedisDao) IsTokenValid(token string) (bool, error) {
	val, err := r.rdb.TTL(token).Result()
	if err != nil {
		return false, err
	}
	// once the key is expired, ttl key's result is -2,
	// so val == -2 means this token is expired.
	if val == -2 {
		return false, nil
	}
	return true, nil
}

func (r *RedisDao) AddToNameList(username string) error {
	return r.rdb.SAdd("namelist", username).Err()
}

func (r *RedisDao) IsUserNameExist(username string) (bool, error) {
	return r.rdb.SIsMember("namelist", username).Result()
}

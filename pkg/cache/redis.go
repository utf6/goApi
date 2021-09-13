package cache

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/utf6/goApi/pkg/config"
	"time"
)

var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			c,err := redis.Dial("tcp", config.Redis.Host)
			if err != nil {
				return nil, err
			}

			if config.Redis.Password != "" {
				if _, err := c.Do("AUTH", config.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         config.Redis.MaxIdle,
		MaxActive:       config.Redis.MaxActive,
		IdleTimeout:     config.Redis.IdleTimeout,
		Wait:            false,
	}
	return nil
}

//redis 写入
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}
	return  nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

// 获取 key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return  nil,err
	}
	return  reply, nil
}

func Del(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDel(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return  err
	}

	for _, key := range keys{
		_, err = Del(key)
		if err != nil {
			return  err
		}
	}

	return  nil
}

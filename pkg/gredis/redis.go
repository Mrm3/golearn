package gredis

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"immortality-demo/config"
	"time"
)

var RedisConn *redis.Pool

// Setup Initialize the Redis instance
func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     config.Config.RedisMaxIdle,
		MaxActive:   config.Config.RedisMaxActive,
		IdleTimeout: config.Config.RedisIdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Config.RedisHost)
			if err != nil {
				return nil, err
			}
			if config.Config.RedisPassword != "" {
				if _, err := c.Do("AUTH", config.Config.RedisPassword); err != nil {
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
	}

	return nil
}

func Push(key string, data interface{}) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("rpush", key, value)
	if err != nil {
		return err
	}

	return nil
}

func TryPop(key string) (interface{}, error) {
	var reply interface{}
	count, err := Len(key)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		conn := RedisConn.Get()
		defer conn.Close()

		reply, err = conn.Do("LPOP", key)
		if err != nil {
			return nil, err
		}
	}
	return reply, nil
}

func Len(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	lenqueue, err := conn.Do("llen", key)
	if err != nil {
		return 0, err
	}

	count, ok := lenqueue.(int64)
	if !ok {
		return 0, errors.New("类型转换错误!")
	}
	return int(count), nil
}

// Set a key/value
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

	return nil
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

// Get get a key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

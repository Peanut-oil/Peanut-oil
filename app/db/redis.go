package db

import (
	"github.com/gin-gonic/gin/app/def"
	"github.com/gin-gonic/gin/app/pkg"
	"github.com/gomodule/redigo/redis"
	"log"
	"runtime/debug"
	"time"
)

type redisClient struct {
	pool *redis.Pool
}

var MainRedis *redisClient

func ConnectRedis() {
	MainRedis = newRedisClient(def.RedisAddr, def.RedisPassWord)
}

func newRedisClient(addr, password string) *redisClient {
	c := new(redisClient)
	c.pool = &redis.Pool{
		Wait:        true,
		MaxIdle:     50,
		MaxActive:   200,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if len(password) != 0 {
				if _, err := conn.Do("AUTH", password); err != nil {
					return nil, err
				}
			}
			return conn, nil
		},
	}
	//测试连接
	_, err := redis.String(c.Do("PING"))
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func (c *redisClient) Do(commandName string, args ...interface{}) (interface{}, error) {
	conn := c.pool.Get()
	reply, err := conn.Do(commandName, args...)
	if err != nil && err != redis.ErrNil {
		pkg.Errorln("[Do] redis err stack", err, string(debug.Stack()))
	}
	conn.Close()
	return reply, err
}

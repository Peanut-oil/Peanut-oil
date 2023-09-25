package db

import (
	"fmt"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gin-gonic/gin/app/pkg"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"log"
	"runtime/debug"
	"time"
)

type SendCommand struct {
	CommandName string
	Args        []interface{}
}

type redisClient struct {
	pool *redis.Pool
}

var MainRedis *redisClient

const (
	usePipeSendCommand = true
)

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

func (c *redisClient) Send(commands []SendCommand) (interface{}, error) {
	if usePipeSendCommand {
		return c.SendV2(commands)
	}
	conn := c.pool.Get()
	conn.Send("MULTI")
	for _, command := range commands {
		err := conn.Send(command.CommandName, command.Args...)
		if err != nil && err != redis.ErrNil {
			logrus.Errorf("redis send err. %+v", err)
		}
	}
	reply, err := conn.Do("EXEC")
	if err != nil && err != redis.ErrNil {
		logrus.Errorf("redis err. %+v", err)
	}
	conn.Close()
	return reply, err
}

func (c *redisClient) SendV2(commands []SendCommand) (interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	var err error
	for _, command := range commands {
		err = conn.Send(command.CommandName, command.Args...)
		if err != nil && err != redis.ErrNil {
			logrus.WithFields(logrus.Fields{"command": command.CommandName, "args": command.Args}).Errorln(
				"[Send] redis send error:%s", err.Error())
		}
	}
	err = conn.Flush()
	if err != nil {
		logrus.Errorln("[Send] flush command error:%s", err.Error())
		return nil, err
	}
	reply := make([]interface{}, 0, len(commands))
	for i := 0; i < len(commands); i++ {
		singleReply, recErr := redis.ReceiveWithTimeout(conn, time.Millisecond*500)
		if recErr != nil && recErr != redis.ErrNil {
			err = fmt.Errorf("name:%s,args:%v,exec err:%w", commands[i].CommandName, commands[i].Args, recErr)
		}
		reply = append(reply, singleReply)
	}
	if err != nil && err != redis.ErrNil {
		logrus.Errorf("[Send] redis err. %+v", err)
	}

	return reply, err
}

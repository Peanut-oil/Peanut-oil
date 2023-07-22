package store

import (
	"database/sql"
	"github.com/gin-gonic/gin/app/dao"
	"github.com/gin-gonic/gin/app/db"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gin-gonic/gin/app/helper"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

func saveCacheUser(info *dao.UserInfo) error {
	key := def.HSetUserInfo + info.DeviceId
	m, err := helper.Struct2Map(*info, "redis")
	if err != nil {
		logrus.Errorf("deviceId:%s, err:%s", info.DeviceId, err.Error())
		return err
	}
	db.MainRedis.Do("HMSet", redis.Args{}.Add(key).AddFlat(m)...)
	db.MainRedis.Do("expire", key, def.UserInfoExpire)
	return nil
}

func getDBUserInfo(deviceId string) *dao.UserInfo {
	info := new(dao.UserInfo)
	selectSql := "select * from " + def.TableUserInfo + " where device_id = ?"
	err := db.MysqlDB.Unsafe().QueryRowx(selectSql, deviceId).StructScan(info)
	if err != nil {
		if err != sql.ErrNoRows {
			logrus.WithField("deviceId", deviceId).WithField("err", err).Warn("get db user info err")
		}
		return nil
	}
	return info
}

func checkUserInfoTTL(deviceId string) {
	key := def.HSetUserInfo + deviceId
	ttl, err := redis.Int64(db.MainRedis.Do("ttl", key))
	if err != nil {
		logrus.Error(err)
		return
	}
	if ttl == -2 {
		info := getDBUserInfo(deviceId)
		if info == nil {
			return
		}
		saveCacheUser(info)
		return
	}
	if ttl < def.UserInfoExpire {
		db.MainRedis.Do("expire", key, def.UserInfoExpire)
	}
}

func GetUserInfoByDeviceId(deviceId string) *dao.UserInfo {
	checkUserInfoTTL(deviceId)

	logCtx := logrus.WithField("deviceId", deviceId)
	key := def.HSetUserInfo + deviceId
	res, err := redis.Values(db.MainRedis.Do("hGetAll", key))
	if err != nil {
		if err != redis.ErrNil {
			logCtx.WithField("err", err).Error("hgetall user err")
		}
		return nil
	}
	u := new(dao.UserInfo)
	err = redis.ScanStruct(res, u)
	if err != nil {
		logCtx.WithField("err", err).Error("user info scan err")
		return nil
	}
	if u.Uid == 0 {
		return nil
	}

	return u
}

package store

import (
	"github.com/gin-gonic/gin/app/dao"
	"github.com/gin-gonic/gin/app/db"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetTopList(rankType, limit int) []dao.ZsetItem {
	list := make([]dao.ZsetItem, 0)
	key := def.ZSetRankList + strconv.Itoa(rankType)
	reply, err := redis.Values(db.MainRedis.Do("ZREVRANGEBYSCORE", key, "+inf", 0, "WITHSCORES", "LIMIT", 0, limit))
	if err != nil && err != redis.ErrNil {
		logrus.Errorln("[GetTopList] values err", err)
	}
	err = redis.ScanSlice(reply, &list)
	if err != nil {
		logrus.Errorln("[GetTopList] scan err", err)
	}
	return list
}

func AddRankScore(score, uid, rankType int) error {
	key := def.ZSetRankList + strconv.Itoa(rankType)
	_, err := db.MainRedis.Do("ZIncrBy", key, score, uid)
	if err != nil {
		logrus.WithFields(logrus.Fields{"uid": uid, "score": score}).Errorln("AddRankScore err:", err)
		return err
	}
	return nil
}

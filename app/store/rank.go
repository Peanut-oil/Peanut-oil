package store

import (
	"github.com/gin-gonic/gin/app/dao"
	"github.com/gin-gonic/gin/app/db"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetTopList(scoreType, limit int) []dao.ZsetItem {
	list := make([]dao.ZsetItem, 0)
	key := def.ZSetRankList + strconv.Itoa(scoreType)
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

func AddRankScore(scoreTime, scoreSpeed, scoreHeight int, did string) error {
	timeKey := def.ZSetRankList + strconv.Itoa(def.RankTypeTime)
	speedKey := def.ZSetRankList + strconv.Itoa(def.RankTypeSpeed)
	heightKey := def.ZSetRankList + strconv.Itoa(def.RankTypeHeight)

	db.MainRedis.Do("ZIncrBy", timeKey, scoreTime, did)
	db.MainRedis.Do("ZIncrBy", speedKey, scoreSpeed, did)
	db.MainRedis.Do("ZIncrBy", heightKey, scoreHeight, did)
	return nil
}

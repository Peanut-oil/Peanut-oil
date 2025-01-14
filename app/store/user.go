package store

import (
	"database/sql"
	"github.com/gin-gonic/gin/app/dao"
	"github.com/gin-gonic/gin/app/db"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gin-gonic/gin/app/helper"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
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

func PipeGetUserInfo(dids []string) ([]*dao.UserInfo, error) {
	if len(dids) == 0 {
		return []*dao.UserInfo{}, nil
	}
	var commands []db.SendCommand
	pipelineCheckUserInfoTtl(dids)
	for _, did := range dids {
		key := def.HSetUserInfo + did
		commands = append(commands, db.SendCommand{
			CommandName: "HGETALL",
			Args:        []interface{}{key},
		})
	}
	logCtx := logrus.WithField("dids", dids)
	r, err := redis.Values(db.MainRedis.Send(commands))
	if err != nil {
		logCtx.Error("pipeline hgetall user err", err)
		return nil, err
	}
	res := make([]*dao.UserInfo, 0)
	for _, r1 := range r {
		u := new(dao.UserInfo)
		if r2, ok := r1.([]interface{}); ok {
			err = redis.ScanStruct(r2, u)
			if err != nil {
				logCtx.WithField("r2", r2).Error("pipeline hgetall user ScanStruct err", err)
				return nil, err
			}
		}
		res = append(res, u)
	}
	return res, nil
}

func pipelineCheckUserInfoTtl(dids []string) {
	var commands []db.SendCommand
	for _, did := range dids {
		key := def.HSetUserInfo + did
		commands = append(commands, db.SendCommand{
			CommandName: "ttl",
			Args:        []interface{}{key},
		})
	}
	logCtx := logrus.WithField("dids", dids)
	r, err := redis.Ints(db.MainRedis.Send(commands))
	if err != nil {
		logCtx.Error("pipeline ttl user err", err)
		return
	}

	needLoadUidList := make([]string, 0, len(dids)/4)
	expireCommands := make([]db.SendCommand, 0, len(dids)/2)
	for i, ttl := range r {
		if ttl == -2 {
			needLoadUidList = append(needLoadUidList, dids[i])
		} else if ttl < def.DefaultResetExpire {
			expireCommands = append(expireCommands, db.SendCommand{
				CommandName: "Expire",
				Args:        []interface{}{def.HSetUserInfo + dids[i], def.UserInfoExpire},
			})
		}
	}

	if len(needLoadUidList) > 0 {
		userInfos, err := batchGetDBUserInfo(needLoadUidList)
		if err != nil {
			logrus.WithField("uid_list", needLoadUidList).Errorf("[pipelineCheckUserInfoTtl] batch get db user info error:%v", err)
		} else if len(userInfos) > 0 {
			existUidList := make([]int, 0, len(userInfos))
			for _, userInfo := range userInfos {
				existUidList = append(existUidList, userInfo.Uid)
			}
			cacheCommands := make([]db.SendCommand, 0, len(existUidList))
			for i := 0; i < len(userInfos); i++ {
				did := userInfos[i].DeviceId
				fieldsMap, err := helper.Struct2Map(userInfos[i], "redis")
				if err != nil {
					logrus.WithField("did", did).WithField("info", userInfos[i]).Errorf("[pipelineCheckUserInfoTtl] struct to map error:%s", err.Error())
					continue
				}
				infoKey := def.HSetUserInfo + did
				cacheCommands = append(cacheCommands, db.SendCommand{
					CommandName: "HMSet",
					Args:        redis.Args{}.Add(infoKey).AddFlat(fieldsMap),
				})
				expireCommands = append(expireCommands, db.SendCommand{
					CommandName: "Expire",
					Args:        []interface{}{infoKey, def.UserInfoExpire},
				})
			}

			if len(cacheCommands) > 0 {
				_, err = db.MainRedis.Send(cacheCommands)
				if err != nil {
					logrus.WithField("uid_list", existUidList).Errorf("[pipelineCheckUserInfoTtl] pipe cache user info error:%s", err.Error())
				}
			}
		}
	}

	if len(expireCommands) > 0 {
		_, err = db.MainRedis.Send(expireCommands)
		if err != nil {
			logrus.WithField("did_list", dids).Errorf("[pipelineCheckUserInfoTtl] pipe expire user"+
				" level info error:%s", err.Error())
		}
	}
}

func batchGetDBUserInfo(didList []string) (userInfos []dao.UserInfo, err error) {
	oriSql := "select * from " + def.TableUserInfo + " where device_id in (?)"
	selectSql, args, err := sqlx.In(oriSql, didList)
	if err != nil {
		return make([]dao.UserInfo, 0), err
	}
	userInfos = make([]dao.UserInfo, 0, len(didList))
	err = db.MysqlDB.Unsafe().Select(&userInfos, selectSql, args...)
	return
}

func AddUserInfo(info *dao.UserInfo) (int, error) {
	uid, err := insertUserInfo(info)
	if err != nil {
		return 0, err
	}
	info.Uid = uid
	err = saveCacheUser(info)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func insertUserInfo(info *dao.UserInfo) (int, error) {
	sqlStr := "INSERT INTO `" + def.TableUserInfo + "` ("
	// 获取所有field及其value
	fields, values := helper.GetStructFieldsAndValuesExcept(*info, []string{})
	query := sqlStr + strings.Join(fields, ",") + ") values (" + strings.Join(values, ",") + ")"
	res, err := db.MysqlDB.NamedExec(query, info)
	if err != nil {
		logrus.WithField("err", err.Error()).Warn("insert user info err")
		return 0, err
	}
	uid, err := res.LastInsertId()
	if err != nil {
		logrus.WithField("err", err.Error()).Warn("get insert id err")
		return 0, err
	}
	return int(uid), nil
}

func UpdateUserInfo(fields map[string]interface{}, uid int, did string) error {
	set := ""
	for k := range fields {
		set += k + "=:" + k + ","
	}
	set = strings.TrimRight(set, ",")
	fields["uid"] = uid
	fields["update_time"] = int(time.Now().Unix())
	_, err := db.MysqlDB.NamedExec("update user_info set "+set+" where uid=:uid", fields)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uid": uid,
			"set": set,
			"did": did,
		}).Warn(err)
		return err
	}
	checkUserInfoTTL(did)
	key := def.HSetUserInfo + did
	_, err = db.MainRedis.Do("hMSet", redis.Args{}.Add(key).AddFlat(fields)...)
	if err != nil {
		logrus.WithField("err", err).Warn("hMSet err")
		return err
	}
	db.MainRedis.Do("expire", key, def.UserInfoExpire)
	return nil
}

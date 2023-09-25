package service

import (
	"errors"
	"github.com/gin-gonic/gin/app/dao"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gin-gonic/gin/app/helper"
	"github.com/gin-gonic/gin/app/store"
	"github.com/sirupsen/logrus"
)

func UserLoginByDeviceId(deviceId string) (int, error) {
	logrus.WithFields(logrus.Fields{"deviceId": deviceId})
	userInfo := store.GetUserInfoByDeviceId(deviceId)
	// 用户信息为null，直接注册
	if userInfo == nil {
		registerInfo := helper.GenerateUserInfo(deviceId)
		// 更新信息
		uid, err := store.AddUserInfo(registerInfo)
		if err != nil {
			logrus.Errorf("[UserLoginByDeviceId] AddUserInfo err:%s", err.Error())
			return 0, errors.New(def.MsgSystemErr)
		}
		return uid, nil
	}

	return userInfo.Uid, nil
}

func GetRankList(rankType int) ([]*dao.RankUserInfo, error) {
	topList := store.GetTopList(rankType, 10)
	res := make([]*dao.RankUserInfo, 0)
	if len(topList) == 0 {
		return res, errors.New("无排行数据")
	}
	uids := make([]int, len(topList))
	for _, info := range topList {
		uids = append(uids, info.Member)
	}
	userInfos, err := store.PipeGetUserInfo(uids)
	if err != nil || len(uids) != len(userInfos) {
		return res, errors.New("系统异常")
	}

	for index, info := range topList {
		item := &dao.RankUserInfo{
			Uid:      info.Member,
			Score:    info.Score,
			Rank:     index + 1,
			Avatar:   userInfos[index].Avatar,
			NickName: userInfos[index].NickName,
		}
		res = append(res, item)
	}

	return res, nil
}

package service

import (
	"errors"
	"github.com/gin-gonic/gin/app/dao"
	"github.com/gin-gonic/gin/app/helper"
	"github.com/gin-gonic/gin/app/pkg/graceful"
	"github.com/gin-gonic/gin/app/store"
	"github.com/sirupsen/logrus"
	"time"
)

func userLoginByDeviceId(scoreTime, scoreSpeed, scoreHeight int, nickName, avatar, country, deviceId string) (*dao.UserInfo, error, bool) {
	logrus.WithFields(logrus.Fields{"deviceId": deviceId})
	userInfo := store.GetUserInfoByDeviceId(deviceId)
	// 用户信息为null，直接注册
	if userInfo == nil {
		registerInfo := &dao.UserInfo{
			NickName:    nickName,
			Avatar:      avatar,
			Country:     country,
			DeviceId:    deviceId,
			ScoreTime:   scoreTime,
			ScoreSpeed:  scoreSpeed,
			ScoreHeight: scoreHeight,
			CreateTime:  int(time.Now().Unix()),
			UpdateTime:  int(time.Now().Unix()),
		}
		// 更新信息
		_, err := store.AddUserInfo(registerInfo)
		if err != nil {
			logrus.Errorf("[UserLoginByDeviceId] AddUserInfo err:%s", err.Error())
			return &dao.UserInfo{}, err, false
		}
		return registerInfo, nil, true
	} else { // 更新用户信息
		changeFields := helper.GetUpdateUserInfoChangeFields(userInfo, avatar, nickName, country, scoreTime, scoreSpeed, scoreHeight)
		if len(changeFields) > 0 {
			graceful.Go(func() {
				err := store.UpdateUserInfo(changeFields, userInfo.Uid, deviceId)
				if err != nil {
					logrus.Errorf("[UserLoginByDeviceId] UpdateUserInfo err:%s", err.Error())
					return
				}
			})
		}
	}

	return userInfo, nil, false
}

func GetRankList(rankTypeOneClass, rankTypeTwoClass int) ([]*dao.RankUserInfo, error) {
	topList := store.GetTopList(rankTypeOneClass, rankTypeTwoClass, 10)
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

func AddRankScoreWithInfo(scoreTime, scoreSpeed, scoreHeight int, nickName, avatar, country, did string) error {
	// 首先获取用户信息并且更新
	userInfo, err, isFresh := userLoginByDeviceId(scoreTime, scoreSpeed, scoreHeight, nickName, avatar, country, did)
	if err != nil {
		return err
	}

	addScoreTime := scoreTime
	addScoreSpeed := scoreSpeed
	addScoreHeigh := scoreHeight
	if !isFresh {
		addScoreTime += userInfo.ScoreTime
		addScoreSpeed += userInfo.ScoreSpeed
		addScoreHeigh += userInfo.ScoreHeight
	}
	store.AddRankScore(addScoreTime, addScoreSpeed, addScoreHeigh, did)

	return nil
}

package service

import (
	"errors"
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

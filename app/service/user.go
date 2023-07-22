package service

import "github.com/gin-gonic/gin/app/store"

func UserLoginByDeviceId(deviceId string) error {
	userInfo := store.GetUserInfoByDeviceId(deviceId)
	// 用户信息为null，直接注册
	if userInfo == nil {

	}

	return nil
}

package helper

import "github.com/gin-gonic/gin/app/dao"

func GenerateUserInfo(deviceId string) *dao.UserInfo {
	info := new(dao.UserInfo)
	info.NickName = GenerateNickName(deviceId)
	info.Avatar = dao.DefaultAvatar
	info.Gender = dao.DefaultGender
	info.CreateTime = Millisecond() / 1000

	return info
}

func GenerateNickName(deviceId string) string {
	return "genus" + deviceId[35:40]
}

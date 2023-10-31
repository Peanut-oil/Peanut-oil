package helper

import "github.com/gin-gonic/gin/app/dao"

func GenerateUserInfoWithNameAndAvatar(deviceId, nickName, avatar string) *dao.UserInfo {
	info := new(dao.UserInfo)
	info.DeviceId = deviceId
	info.NickName = nickName
	info.Avatar = avatar
	info.Gender = dao.DefaultGender
	info.CreateTime = Millisecond() / 1000
	info.UpdateTime = Millisecond() / 1000

	return info
}

func GenerateUserInfo(deviceId string) *dao.UserInfo {
	info := new(dao.UserInfo)
	info.DeviceId = deviceId
	info.NickName = GenerateNickName(deviceId)
	info.Avatar = dao.DefaultAvatar
	info.Gender = dao.DefaultGender
	info.CreateTime = Millisecond() / 1000
	info.UpdateTime = Millisecond() / 1000

	return info
}

// todo
func GenerateNickName(deviceId string) string {
	return "genus" + deviceId[1:5]
}

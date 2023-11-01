package helper

import "github.com/gin-gonic/gin/app/dao"

func GetUpdateUserInfoChangeFields(info *dao.UserInfo, avatar, nickname, country string, scoreTime, scoreSpeed, scoreHeight int) map[string]interface{} {
	fields := make(map[string]interface{}, 0)
	if avatar != "" && info.Avatar != avatar {
		fields["avatar"] = avatar
	}
	if nickname != "" && info.NickName != nickname {
		fields["nickname"] = nickname
	}
	if country != "" && info.Country != country {
		fields["country"] = country
	}
	if scoreTime != 0 {
		fields["score_time"] = info.ScoreTime + scoreTime
	}
	if scoreSpeed != 0 {
		fields["score_speed"] = info.ScoreSpeed + scoreSpeed
	}
	if scoreHeight != 0 {
		fields["score_height"] = info.ScoreHeight + scoreHeight
	}

	return fields
}

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
	if scoreTime != 0 && scoreTime != info.ScoreTime {
		fields["score_time"] = scoreTime
	}
	if scoreSpeed != 0 && scoreSpeed != info.ScoreSpeed {
		fields["score_speed"] = scoreSpeed
	}
	if scoreHeight != 0 && scoreHeight != info.ScoreHeight {
		fields["score_height"] = scoreHeight
	}

	return fields
}

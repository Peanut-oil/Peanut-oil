package dao

const (
	// todo
	DefaultAvatar = ""
	DefaultGender = GenderMan

	GenderMan   = 1
	GenderWoMan = 2

	UserLogOut = 1
)

type ZsetItem struct {
	Member string `json:"member"`
	Score  int    `json:"score"`
}

type RankUserInfo struct {
	List     []*RankInfo `json:"list"`
	SelfRank *RankInfo   `json:"self_rank"`
}

type RankInfo struct {
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
	Country  string `json:"country"`
	Avatar   string `json:"avatar"`
	NickName string `json:"nickname"`
}

type UserInfo struct {
	DeviceId    string `json:"device_id" redis:"did" db:"device_id"`
	Uid         int    `json:"uid" redis:"uid" db:"uid"`
	NickName    string `json:"nickname" redis:"nickname" db:"nickname"`
	Avatar      string `json:"avatar" redis:"avatar" db:"avatar"`
	CreateTime  int    `json:"create_time" redis:"ct" db:"create_time"`
	UpdateTime  int    `json:"update_time" redis:"ut" db:"update_time"`
	Country     string `json:"country" redis:"country" db:"country"`
	ScoreTime   int    `json:"score_time" redis:"s_time" db:"score_time"`
	ScoreSpeed  int    `json:"score_speed" redis:"s_speed" db:"score_speed"`
	ScoreHeight int    `json:"score_height" redis:"s_height" db:"score_height"`
}

package def

const (
	StringUserLock = "genus:user_lock:"
	HSetUserInfo   = "genus:user_info:"
)

// 排行榜
const (
	ZSetRankList   = "genus:rank_list:"
	RankTypeTime   = 1
	RankTypeSpeed  = 2
	RankTypeHeight = 3
)

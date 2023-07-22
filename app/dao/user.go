package dao

type UserInfo struct {
	DeviceId   string `json:"device_id" redis:"deviceid" db:"device_id"`
	Uid        int    `json:"uid" redis:"uid" db:"uid"`
	NickName   string `json:"nickname" redis:"nickname" db:"nickname"`
	Avatar     string `json:"avatar" redis:"avatar" db:"avatar"`
	Gender     int    `json:"gender" redis:"gender" db:"gender"`
	CreateTime int    `json:"create_time" redis:"create_time" db:"create_time"`
	Coin       int    `json:"coin" redis:"coin" db:"coin"`
	Phone      string `json:"phone" redis:"phone" db:"phone"`
	Logout     bool   `json:"logout" db:"logout" redis:"logout"`
}

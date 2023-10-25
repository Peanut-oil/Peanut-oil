package def

// 配置相关
const (
	ServerPort    = "10087"
	RedisPassWord = "xf233@123"
	RedisAddr     = "localhost:6379"
	MysqlAddr     = "root:qq937634115@tcp(localhost:3306)/elastic_pet?charset=utf8mb4"
	ServerLog     = "/root/serverlog/api"
)

const (
	CodeSucc           = 200
	CodeError          = 500
	CodeUnAuth         = 401
	CodeRequestTooFast = 603

	MsgSystemErr = "system error,please try again later"
	MsgParamErr  = "param err"

	// todo
	JwtEncryptKey = "genuisskxhuisvusicabbcosjxcbsodnojivyf"
	SignKey       = "genuisiscyewvcmihuxwyubxewniciegvc"
)

// 时间相关
const (
	MinSecond          = 60
	HourSecond         = 60 * 60
	DaySecond          = 24 * 60 * 60
	WeekSecond         = DaySecond * 7
	UserInfoExpire     = DaySecond * 7
	DefaultResetExpire = DaySecond * 3

	TimeFormatDay              = "2006-01-02"
	TimeFormatDayNoSep         = "20060102"
	TimeFormatYMD              = "20060102150405"
	TimeFormatDateWord         = "2006年01月02日"
	TimeFormatDotSimpleMD      = "1.2"
	TimeFormatWordMD           = "01月02日"
	TimeFormatActivity         = "01月02日 15:04:05"
	TimeFormatDate             = "2006-01-02 15:04:05"
	TimeFormatMDNoSep          = "0102"
	TimeFormatDateHM           = "2006-01-02 15:04"
	TimeFormatDateHMS          = "15:04:05"
	TimeFormatYear             = "2006"
	TimeFormatMonth            = "01"
	TimeFormatOnlyDay          = "02"
	TimeFormatMinuteFilterYear = "01-02 15:04"

	TimeCalDaySeconds = 86400 // 一天的秒数
)

var (
	LogWhiteMethod  = []string{}
	SignWhiteMethod = []string{}
)

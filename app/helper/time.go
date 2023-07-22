package helper

import (
	"github.com/gin-gonic/gin/app/def"
	"math"
	"strconv"
	"time"
)

func GetNowTime() int {
	return int(time.Now().Unix())
}

func GetNowHour() int {
	hour := time.Now().Hour()
	return hour
}

// GetIntervalDay 获取两个unix之间的间隔天数（精确到天）
func GetIntervalDay(startUnixTime int64, endUnixTime int64) int {
	startFormat := time.Unix(startUnixTime, 0).Format(def.TimeFormatDay)
	endFormat := time.Unix(endUnixTime, 0).Format(def.TimeFormatDay)
	startDay, err := time.ParseInLocation(def.TimeFormatDay, startFormat, time.Local)
	if err != nil {
		return 0
	}
	endDay, err := time.ParseInLocation(def.TimeFormatDay, endFormat, time.Local)
	if err != nil {
		return 0
	}
	return int(endDay.Sub(startDay).Hours() / 24)
}

func DiffNatureDays(t1, t2 int64) int {
	return int(math.Ceil(time.Unix(t1, 0).Sub(time.Unix(t2, 0)).Hours() / 24))
}

func GetOtherDayStartTimeWithToday(t time.Time, days int) time.Time {
	todayTimeStr := t.AddDate(0, 0, 1+days).Format(def.TimeFormatDay)
	todayBDTime, _ := time.ParseInLocation(def.TimeFormatDay, todayTimeStr, time.Local)
	return todayBDTime
}

func GetOtherDayEndTimeWithToday(t time.Time, days int) time.Time {
	todayTimeStr := t.AddDate(0, 0, days).Format(def.TimeFormatDay)
	todayBDTime, _ := time.ParseInLocation(def.TimeFormatDay, todayTimeStr, time.Local)
	return todayBDTime.Add(-time.Second)
}

func OtherDayMillisecond(d int) int64 {
	now := time.Now()
	day := now.AddDate(0, 0, d)
	return day.UnixNano() / 1000000
}

func MillisecondInt64() int64 {
	return time.Now().UnixNano() / 1000000
}

func Millisecond() int {
	return int(time.Now().UnixNano()) / 1000000
}

func GetNow() int64 {
	return time.Now().Unix()
}

func GetToday() string {
	return time.Now().Format(def.TimeFormatDay)
}

func GetTodayFormatNoSep() string {
	return time.Now().Format(def.TimeFormatDayNoSep)
}

func GetTimestampDate(unix int) string {
	return time.Unix(int64(unix), 0).Format(def.TimeFormatDay)
}

func GetTimestampWholeDate(unix int) string {
	if unix == 0 {
		return ""
	}
	return time.Unix(int64(unix), 0).Format(def.TimeFormatDate)
}
func GetYesterday() string {
	return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
}

func GetTodayInt() int {
	now := time.Now()
	d, _ := strconv.Atoi(now.Format("20060102"))
	return d
}

func GetDays(day int) []string {
	days := []string{}
	now := time.Now()
	for i := 0; i < day; i++ {
		d := now.AddDate(0, 0, i)
		days = append(days, d.Format("2006-01-02"))
	}
	return days
}

func GetThisWeekMondy() int {
	now := time.Now()
	wd := int(now.Weekday())
	if wd == 0 {
		wd = 7
	}
	period := 1 - wd
	monday := now.AddDate(0, 0, period)
	currentLocation := now.Location()

	zeroMonday := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, currentLocation)

	return int(zeroMonday.Unix())
}

func GetStartOfTodayUnix() int {
	now := time.Now()
	currentLocation := now.Location()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, currentLocation)
	return int(startOfToday.Unix())
}

func GetEndOfTodayUnix() int {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	currentLocation := now.Location()
	endToday := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, currentLocation)
	return int(endToday.Unix())
}

func GetWeekAgoUnix() int {
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	return int(weekAgo.Unix())
}

func Get3MonthAgoUnix() int {
	now := time.Now()
	return int(now.AddDate(0, -3, 0).Unix())
}

func GetMonthAgoUnix(nMonth int) int {
	now := time.Now()
	return int(now.AddDate(0, -nMonth, 0).Unix())
}

// 获取时间戳n天后的时间戳
func GetOtherDayTime(oldUnix, addDays int) int {
	if oldUnix == 0 {
		now := time.Now()
		day := now.AddDate(0, 0, addDays)
		return int(day.Unix())
	}
	oldDay := time.Unix(int64(oldUnix), 0)
	newDay := oldDay.AddDate(0, 0, addDays)
	return int(newDay.Unix())
}

// 获取x天y小时后的时间戳
func GetOtherDayHourTime(oldUnix, x, y int) int {
	if oldUnix == 0 {
		now := time.Now()
		day := now.AddDate(0, 0, x).Add(time.Duration(y) * time.Hour)
		return int(day.Unix())
	}
	oldDay := time.Unix(int64(oldUnix), 0)
	newDay := oldDay.AddDate(0, 0, x).Add(time.Duration(y) * time.Hour)
	return int(newDay.Unix())
}

// 获取今天剩余时间
func GetTodayExpireTime() int {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	currentLocation := now.Location()
	endToday := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, currentLocation)
	return int(endToday.Unix()) - int(now.Unix())
}

func GetThisWeekExpireTime() int {
	now := time.Now()
	wd := int(now.Weekday())
	if wd == 0 {
		wd = 7
	}
	period := 8 - wd
	nextMonday := now.AddDate(0, 0, period)
	currentLocation := now.Location()
	zeroNextMonday := time.Date(nextMonday.Year(), nextMonday.Month(), nextMonday.Day(), 0, 0, 0, 0, currentLocation)
	return int(zeroNextMonday.Unix()) - int(now.Unix())
}

func GetThisMonthExpireTime() int {
	now := time.Now()
	lastDay := GetLastDayOfMonth(now).Unix() + def.DaySecond
	return int(lastDay) - int(now.Unix())
}

func GetFirstDayOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

func GetLastDayOfMonth(d time.Time) time.Time {
	return GetFirstDayOfMonth(d).AddDate(0, 1, -1)
}

func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取昨天开始到结束的时间戳
func GetYesterdayTimeUnix() (int, int) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	currentLocation := now.Location()

	zeroYestoday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, currentLocation)
	endYestoday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, currentLocation)

	return int(zeroYestoday.Unix()), int(endYestoday.Unix())
}

func GetYesterdayDate() string {
	return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
}

func GetTomorrowDate() string {
	return time.Now().AddDate(0, 0, 1).Format(def.TimeFormatDay)
}

func GetTimeSub(old int, new int) float64 {
	tm := time.Unix(int64(old), 0)
	tm2 := time.Unix(int64(new), 0)
	a := tm2.Sub(tm)
	return a.Hours()
}

func UnixToFormatString(unix int) string {
	tm := time.Unix(int64(unix), 0)
	return tm.Format("2006-01-02")
}

func GetSubDays(oldUnix int) int {
	t := time.Unix(int64(oldUnix), 0)
	hours := time.Now().Sub(t).Hours()
	return int(hours) / 24
}

func FormatUnixTime(unix int, layout string) string {
	return time.Unix(int64(unix), 0).Format(layout)
}

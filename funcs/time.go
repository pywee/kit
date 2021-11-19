package funcs

import (
	"strconv"
	"strings"
	"time"
)

// TimeSub 将当前服务器时间与传入的时间戳进行对比，检查时间差是否在 n 秒以内
func TimeSub(t time.Time, n float64) bool {
	return time.Since(t).Seconds() > n
}

// GetTimeDate 取得指定时间戳当天初始时间
func GetTimeDate(ts int64) time.Time {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	b := time.Unix(ts, 0)
	b = b.In(cstSh)
	year, month, day := b.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
}

// GetTimeDateAddDay
// 取得指定时间戳当天初始时间 支持时间加减
func GetTimeDateAddDay(ts int64, addDay int) time.Time {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	b := time.Unix(ts, 0)
	b = b.In(cstSh)
	year, month, day := b.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, addDay)
}

// GetLastDaysUnix 获取前某段时间的0点时间戳
// startDay: -1; endDay: -30
// 前1天到前30天的时间段 不包括今天
func GetLastDaysUnix(startDay, endDay int) string {
	var (
		str         string
		timeNowUnix = time.Now().Unix()
	)
	for i := startDay; i > endDay; i-- {
		str += strconv.FormatInt(GetTimeDateAddDay(timeNowUnix, i).Unix(), 10) + ","
	}
	return strings.TrimRight(str, ",")
}

// GetFirstDateOfWeek 获取本周日的时间戳
func GetFirstDateOfWeek() int64 {
	now := time.Now()
	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, int(time.Sunday-now.Weekday()))
	return weekStartDate.Unix()
}

// GetThisMonthTime 获取本月第一天的时间
func GetThisMonthTime() time.Time {
	year, month, _ := time.Now().Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

// GetNextMonthTime 获取下个月第一天的时间
func GetNextMonthTime() time.Time {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return thisMonth.AddDate(0, 1, 0)
}

// func timeFormat() {
// 	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
// 	a := int64(1622449109)
// 	b := time.Unix(a, 0)
// 	b = b.In(cstSh)
// 	year, month, day := b.Date()
// 	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
// 	fmt.Println(t)
// }

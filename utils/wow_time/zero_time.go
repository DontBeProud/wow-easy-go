package wow_time

import "time"

// GetTomorrowZeroTime 获取明天零点的时间对象
func GetTomorrowZeroTime() time.Time {
	return GetZeroTimeByDateOffset(1)
}

// GetYesterdayZeroTime 获取昨天零点的时间对象
func GetYesterdayZeroTime() time.Time {
	return GetZeroTimeByDateOffset(-1)
}

// GetZeroTimeByDateOffset 根据日期偏移获取指定日期零点的时间对象
func GetZeroTimeByDateOffset(off int) time.Time {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t.AddDate(0, 0, off)
}

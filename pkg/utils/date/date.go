package date

import (
	"time"

	"github.com/pkg/errors"
)

const YMD_HIS = "2006-01-02T15:04:05"
const Y_M_D_H_I_S = "2006-01-02 15:04:05"
const YYYY_MM_DD = "2006-01-02"

func Format(t *time.Time) string {
	if t != nil {
		return t.Local().Format(Y_M_D_H_I_S)
	}
	return ""
}

// Unix 时间戳
func Unix(t *time.Time) int64 {
	if t == nil {
		temp := time.Now()
		t = &temp
	}
	return t.Unix()
}

// UnixMilli 毫秒时间戳
func UnixMilli(t *time.Time) int64 {
	if t == nil {
		temp := time.Now()
		t = &temp
	}
	return t.UnixNano() / 1e6
}

// 获取当前的时间 - 字符串
func GetCurrentDate() string {
	return time.Now().Format(Y_M_D_H_I_S)
}

func GetCurrentDay() string {
	return time.Now().Format(YYYY_MM_DD)
}

// 获取当前时间戳 - Unix时间戳
func GetCurrentUnix() int64 {
	return time.Now().Unix()
}

// 获取当前时间戳 - 毫秒级时间戳
func GetCurrentMilliUnix() int64 {
	return time.Now().UnixNano() / 1e6
}

// 获取当前时间戳 - 纳秒级时间戳
func GetCurrentNanoUnix() int64 {
	return time.Now().UnixNano()
}

// 获取指定时间戳 - 毫秒级时间戳
func GetMilliUnix(t, layout string) int64 {
	result := GetLocalTime(t, layout)
	return result.UnixNano() / 1e6
}

func GetLocalTime(t, layout string) *time.Time {
	if t == "" {
		return nil
	}
	result, _ := time.ParseInLocation(layout, t, time.Local)
	// log.Printf("GetLocalTime result: %v, t: %s", result, t)
	return &result
}

func GetTimeOrNow(t, layout string) time.Time {
	if t == "" {
		result := time.Now().Local()
		return result
	}
	result, err := time.ParseInLocation(layout, t, time.Local)
	if err != nil {
		panic(errors.WithStack(err))
	}
	return result
}

// GetZeroTimeUTCText 将本地当天零时转为UTC时间
func GetZeroTimeUTCText(d time.Time) string {
	zero := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	return zero.UTC().Format(time.RFC3339)
}

func GetTodayBegin(ti *time.Time) int64 {
	timeStr := ti.Format(YYYY_MM_DD)
	t, err := time.ParseInLocation(Y_M_D_H_I_S, timeStr+" 00:00:00", time.Local)
	if err != nil {
		panic(errors.WithStack(err))
	}
	return t.Unix()
}

// 获取今天的最后结束时间
func GetTodayEnd(ti *time.Time) int64 {
	timeStr := ti.Format(YYYY_MM_DD)
	t, err := time.ParseInLocation(Y_M_D_H_I_S, timeStr+" 23:59:59", time.Local)
	if err != nil {
		panic(errors.WithStack(err))
	}
	return t.Unix()
}

// 格式化传入的时间
func GetTimeFormat(t int64) string {
	return time.Unix(t, 0).Format(Y_M_D_H_I_S)
}

func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

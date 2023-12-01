//  Copyright (c) 2021. The EFF Team Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  See the License for the specific language governing permissions and
//  limitations under the License.

package ztime

import (
	"fmt"
	"strings"
	"time"

	"github.com/ergoapi/util/common"
	"github.com/ergoapi/util/exstr"
	"github.com/golang-module/carbon/v2"
)

// NowUnixString 当前时间时间戳
func NowUnixString() string {
	return exstr.Int642Str(time.Now().Unix())
}

// NowUnix 当前时间戳
func NowUnix() int64 {
	return time.Now().Unix()
}

func NowUTC() time.Time {
	return time.Now().UTC()
}

func UTCify(now time.Time) time.Time {
	if now.IsZero() {
		val := time.Now().UTC()
		return val
	}
	return now.UTC()
}

func Localify(now time.Time) time.Time {
	if now.IsZero() {
		val := time.Now().Local()
		return val
	}
	return now.Local()
}

// TimeParse time parse
func TimeParse(layout, t string) (time.Time, error) {
	return time.Parse(layout, t)
}

// NowFormat 当前时间format
func NowFormat() string {
	return time.Now().Format(common.DefaultTimeLayout)
}

func NextDay(d int, layout string) string {
	return time.Now().AddDate(0, 0, d).Format(layout)
}

func CustomNextDay(day string, d int, layout string) string {
	dayunit := TimeToUninxv2(layout, day)
	return time.Unix(dayunit, 0).AddDate(0, 0, d).Format(layout)
}

// Today0hourUnix 今天0时时间戳
func Today0hourUnix() int64 {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, -1)
	return t1.Unix()
}

// BeforeNowUnix 历史时间戳
func BeforeNowUnix(old string) (oldunix int64) {
	return time.Now().Unix() - exstr.Str2Int64(old)
}

// UnixInt642String unix转化为字符串
func UnixInt642String(t int64) string {
	return time.Unix(t, 0).Format(common.DefaultTimeLayout)
}

// UnixString2String unix转化为字符串
func UnixString2String(t string) string {
	return time.Unix(exstr.Str2Int64(t), 0).Format(common.DefaultTimeLayout)
}

// UnixNanoInt642String unix转化为字符串
func UnixNanoInt642String(t int64) string {
	return time.Unix(0, t).Format(common.DefaultTimeLayout)
}

// UnixNanoString2String unix转化为字符串
func UnixNanoString2String(t string) string {
	return time.Unix(0, exstr.Str2Int64(t)).Format(common.DefaultTimeLayout)
}

// GetTime 获取时间
func GetTime(layout string) string {
	return time.Now().Format(layout)
}

// GetTodayMin 获取今天时间分钟
func GetTodayMin() string {
	return time.Now().Format("200601021504")
}

// GetTodayHour 获取今天时间小时
func GetTodayHour() string {
	return time.Now().Format("2006010215")
}

// GetToday 获取今天时间
func GetToday() string {
	return time.Now().Format("20060102")
}

// GetYesterday 昨天
func GetYesterday(layout ...string) string {
	if len(layout) == 0 {
		layout = append(layout, common.DefaultTimeLayout)
	}
	return time.Now().AddDate(0, 0, -1).Format(layout[0])
}

// GetTomorrow 明天
func GetTomorrow(layout ...string) string {
	if len(layout) == 0 {
		layout = append(layout, common.DefaultTimeLayout)
	}
	return time.Now().AddDate(0, 0, 1).Format(layout[0])
}

// GetNowTimeByLayout layout time
func GetNowTimeByLayout(layout string) string {
	return time.Now().Format(layout)
}

// GetTodaySingleHour 获取今天具体小时
func GetTodaySingleHour() string {
	today := GetToday()
	todayhour := GetTodayHour()
	return strings.ReplaceAll(todayhour, today, "")
}

// GetMonth 获取当前月份
func GetMonth() string {
	return time.Now().Format("200601")
}

// GetShortMonth 获取当前月份
func GetShortMonth() string {
	res := time.Now().Format("2006-01")
	return strings.Split(res, "-")[1]
}

// GetYear 获取当前年份
func GetYear() string {
	return time.Now().Format("2006")
}

// GetWeekFristDayUnix 时间
func GetWeekFristDayUnix() int64 {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStart.Unix()
}

// GetWeekLastDayUnix 时间
func GetWeekLastDayUnix() int64 {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekNextStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset+7)
	return weekNextStart.Unix() - 1
}

// GetWeekDayUnix 时间
func GetWeekDayUnix() (int64, int64) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekNextStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset+7)
	return weekStart.Unix(), weekNextStart.Unix() - 1
}

// GetWeekDayUnixString 时间
func GetWeekDayUnixString() (string, string) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekNextStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset+7)
	return weekStart.Format(common.DefaultTimeLayout), weekNextStart.Format(common.DefaultTimeLayout)
}

// NowAddUnix2Int64 当前时间戳 add
func NowAddUnix2Int64(key string, value time.Duration) int64 {
	lkey := strings.ToLower(key)
	if lkey == "m" || lkey == "minute" {
		return time.Now().Add(time.Minute * value).Unix()
	}
	if lkey == "d" || lkey == "day" {
		return time.Now().Add(time.Hour * 24 * value).Unix()
	}
	if lkey == "w" || lkey == "week" {
		return time.Now().Add(time.Hour * 24 * 7 * value).Unix()
	}
	return time.Now().Add(time.Hour * value).Unix()
}

// NowAddUnix2Str 当前时间戳 add
func NowAddUnix2Str(key string, value time.Duration) string {
	lkey := strings.ToLower(key)
	if lkey == "m" || lkey == "minute" {
		return time.Now().Add(time.Minute * value).String()
	}
	if lkey == "d" || lkey == "day" {
		return time.Now().Add(time.Hour * 24 * value).String()
	}
	if lkey == "w" || lkey == "week" {
		return time.Now().Add(time.Hour * 24 * 7 * value).String()
	}
	return time.Now().Add(time.Hour * value).String()
}

// GetMonthDayNum 获取任已一年月的天数
func GetMonthDayNum(year, month string) int {
	switch month {
	case "01", "03", "05", "07", "08", "10", "12":
		return 31
	case "04", "06", "09", "11":
		return 30
	default:
		if IsLeapYear(exstr.Str2Int(year)) {
			return 29
		}
		return 28
	}
}

// 判断是否为闰年
func IsLeapYear(year int) bool { //y == 2000, 2004
	//判断是否为闰年
	if year%4 == 0 && year%100 != 0 || year%400 == 0 {
		return true
	}
	return false
}

// 时间转时间戳
func TimeToUninx(t string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tt, _ := time.ParseInLocation(common.DefaultTimeLayout, t, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	return tt.Unix()
}

// 时间转时间戳
func TimeToUninxv2(layout, value string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tt, _ := time.ParseInLocation(layout, value, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	return tt.Unix()
}

// 时间转时间戳
func TimeFormatToUninx(layout, value string) (int64, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tt, err := time.ParseInLocation(layout, value, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	return tt.Unix(), err
}

// 某月开始和结束时间戳
func GetMonthStartEndUnix(year, month string) (int64, int64) {
	if exstr.Str2Int(month) < 10 {
		month = fmt.Sprintf("0%v", month)
	}
	st := fmt.Sprintf("%v-%v-01 00:00:00", year, month)
	et := fmt.Sprintf("%v-%v-%v 23:59:59", year, month, GetMonthDayNum(year, month))
	return TimeToUninx(st), TimeToUninx(et)
}

// GetStartQuarter 季度开始时间
func GetStartQuarter(t string) string {
	return carbon.Parse(t).StartOfQuarter().ToDateTimeString()
}

// GetEndQuarter 季度结束时间
func GetEndQuarter(t string) string {
	return carbon.Parse(t).EndOfQuarter().ToDateTimeString()
}

// GetStartMonth 月度开始时间
func GetStartMonth(t string) string {
	return carbon.Parse(t).StartOfMonth().ToDateTimeString()
}

// GetEndMonth 月度结束时间
func GetEndMonth(t string) string {
	return carbon.Parse(t).EndOfMonth().ToDateTimeString()
}

// GetStartWeek 星期开始时间
func GetStartWeek(t string, startMonday bool) string {
	day := carbon.Sunday
	if startMonday {
		day = carbon.Monday
	}
	return carbon.Parse(t).SetWeekStartsAt(day).StartOfWeek().ToDateTimeString()
}

// GetEndWeek 星期结束时间
func GetEndWeek(t string, startMonday bool) string {
	day := carbon.Sunday
	if startMonday {
		day = carbon.Monday
	}
	return carbon.Parse(t).SetWeekStartsAt(day).EndOfWeek().ToDateTimeString()
}

func AddDays(t string, day int) string {
	return carbon.Parse(t).AddDays(day).ToDateTimeString()
}

func SubDays(t string, day int) string {
	return carbon.Parse(t).SubDays(day).ToDateTimeString()
}

func AddHours(t string, hour int) string {
	return carbon.Parse(t).AddHours(hour).ToDateTimeString()
}

func SubHours(t string, hour int) string {
	return carbon.Parse(t).SubHours(hour).ToDateTimeString()
}

// AddDuration add Duration like "2.5h","2h30m"
func AddDuration(t, duration string) string {
	return carbon.Parse(t).AddDuration(duration).ToDateTimeString()
}

func SubDuration(t, duration string) string {
	return carbon.Parse(t).SubDuration(duration).ToDateTimeString()
}

func DiffInDays(t1, t2 string) int64 {
	return carbon.Parse(t1).DiffInDays(carbon.Parse(t2))
}

func DiffInHours(t1, t2 string) int64 {
	return carbon.Parse(t1).DiffInHours(carbon.Parse(t2))
}

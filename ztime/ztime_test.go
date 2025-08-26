// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package ztime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestYesterday(t *testing.T) {
	// 测试昨天相关函数
	yesterdayNow := YesterdayNow()
	assert.NotEmpty(t, yesterdayNow)
	assert.Contains(t, yesterdayNow, " ")

	yesterdayDate := YesterdayDate()
	assert.NotEmpty(t, yesterdayDate)
	assert.Equal(t, 10, len(yesterdayDate)) // YYYY-MM-DD

	yesterdayTS := YesterdayTimestamp()
	assert.Greater(t, yesterdayTS, int64(0))
	assert.Less(t, yesterdayTS, time.Now().Unix())
}

func TestNow(t *testing.T) {
	// 测试当前时间相关函数
	now := Now()
	assert.NotEmpty(t, now)
	assert.Contains(t, now, " ")

	nowDate := NowDate()
	assert.NotEmpty(t, nowDate)
	assert.Equal(t, 10, len(nowDate)) // YYYY-MM-DD

	nowTS := NowTimestamp()
	assert.Greater(t, nowTS, int64(0))
}

func TestTomorrow(t *testing.T) {
	// 测试明天相关函数
	tomorrowNow := TomorrowNow()
	assert.NotEmpty(t, tomorrowNow)
	assert.Contains(t, tomorrowNow, " ")

	tomorrowDate := TomorrowDate()
	assert.NotEmpty(t, tomorrowDate)
	assert.Equal(t, 10, len(tomorrowDate)) // YYYY-MM-DD

	tomorrowTS := TomorrowTimestamp()
	assert.Greater(t, tomorrowTS, time.Now().Unix())
}

func TestMinute(t *testing.T) {
	// 测试分钟操作
	// 正数增加
	after5Min := Minute(5)
	assert.NotEmpty(t, after5Min)

	after5MinTS := MinuteTimestamp(5)
	assert.Greater(t, after5MinTS, time.Now().Unix())

	// 负数减少
	before5Min := Minute(-5)
	assert.NotEmpty(t, before5Min)

	before5MinTS := MinuteTimestamp(-5)
	assert.Less(t, before5MinTS, time.Now().Unix())
}

func TestHour(t *testing.T) {
	// 测试小时操作
	// 正数增加
	after2Hour := Hour(2)
	assert.NotEmpty(t, after2Hour)

	after2HourTS := HourTimestamp(2)
	assert.Greater(t, after2HourTS, time.Now().Unix())

	after2HourDate := HourDate(2)
	assert.NotEmpty(t, after2HourDate)

	// 负数减少
	before2Hour := Hour(-2)
	assert.NotEmpty(t, before2Hour)

	before2HourTS := HourTimestamp(-2)
	assert.Less(t, before2HourTS, time.Now().Unix())

	before2HourDate := HourDate(-2)
	assert.NotEmpty(t, before2HourDate)
}

func TestDay(t *testing.T) {
	// 测试天数操作
	// 正数增加
	after7Days := Day(7)
	assert.NotEmpty(t, after7Days)

	after7DaysDate := DayDate(7)
	assert.NotEmpty(t, after7DaysDate)

	// 负数减少
	before7Days := Day(-7)
	assert.NotEmpty(t, before7Days)

	before7DaysDate := DayDate(-7)
	assert.NotEmpty(t, before7DaysDate)
}

func TestMonth(t *testing.T) {
	// 测试月份操作
	// 正数增加
	after2Months := Month(2)
	assert.NotEmpty(t, after2Months)

	after2MonthsDate := MonthDate(2)
	assert.NotEmpty(t, after2MonthsDate)

	// 负数减少
	before2Months := Month(-2)
	assert.NotEmpty(t, before2Months)

	before2MonthsDate := MonthDate(-2)
	assert.NotEmpty(t, before2MonthsDate)
}

func TestDaysIn(t *testing.T) {
	// 测试天数计算
	daysInCurrentYear := DaysInYear()
	assert.Greater(t, daysInCurrentYear, 364)
	assert.LessOrEqual(t, daysInCurrentYear, 366)

	daysIn2024 := DaysInYear("2024-01-01 12:00:00")
	assert.Equal(t, 366, daysIn2024) // 2024是闰年

	daysIn2023 := DaysInYear("2023-01-01 12:00:00")
	assert.Equal(t, 365, daysIn2023)

	daysInCurrentMonth := DaysInMonth()
	assert.Greater(t, daysInCurrentMonth, 27)
	assert.LessOrEqual(t, daysInCurrentMonth, 31)

	daysInFeb2024 := DaysInMonth("2024-02-15 12:00:00")
	assert.Equal(t, 29, daysInFeb2024)

	daysInFeb2023 := DaysInMonth("2023-02-15 12:00:00")
	assert.Equal(t, 28, daysInFeb2023)
}

func TestAge(t *testing.T) {
	// 测试年龄计算
	age := Age("2000-01-01 12:00:00")
	assert.Greater(t, age, 20)
}

func TestSeason(t *testing.T) {
	// 测试季节
	season := Season()
	assert.NotEmpty(t, season)
	assert.Contains(t, []string{"Spring", "Summer", "Autumn", "Winter"}, season)
}

func TestConstellation(t *testing.T) {
	// 测试星座
	currentConstellation := Constellation()
	assert.NotEmpty(t, currentConstellation)

	// 测试特定日期的星座
	aries := Constellation("2024-04-15 12:00:00") // 白羊座
	assert.Equal(t, "Aries", aries)

	leo := Constellation("2024-08-15 12:00:00") // 狮子座
	assert.Equal(t, "Leo", leo)
}

func TestWeekAndDay(t *testing.T) {
	// 测试周和天相关
	weekOfYear := WeekOfYear()
	assert.Greater(t, weekOfYear, 0)
	assert.LessOrEqual(t, weekOfYear, 53)

	dayOfYear := DayOfYear()
	assert.Greater(t, dayOfYear, 0)
	assert.LessOrEqual(t, dayOfYear, 366)

	dayOfWeek := DayOfWeek()
	assert.GreaterOrEqual(t, dayOfWeek, 1)
	assert.LessOrEqual(t, dayOfWeek, 7)
}

func TestYearStartEnd(t *testing.T) {
	// 测试年份开始结束
	startStr, endStr := YearStartEnd()
	assert.NotEmpty(t, startStr)
	assert.NotEmpty(t, endStr)
	assert.Contains(t, startStr, "-01-01 00:00:00")
	assert.Contains(t, endStr, "-12-31 23:59:59")

	startTS, endTS := YearStartEndTimestamp()
	assert.Greater(t, startTS, int64(0))
	assert.Greater(t, endTS, startTS)

	startDate, endDate := YearStartEndDate()
	assert.Contains(t, startDate, "-01-01")
	assert.Contains(t, endDate, "-12-31")

	// 测试特定年份
	start2024, end2024 := YearStartEnd("2024-06-15 12:00:00")
	assert.Contains(t, start2024, "2024-01-01")
	assert.Contains(t, end2024, "2024-12-31")
}

func TestQuarterStartEnd(t *testing.T) {
	// 测试季度开始结束
	startStr, endStr := QuarterStartEnd()
	assert.NotEmpty(t, startStr)
	assert.NotEmpty(t, endStr)

	startTS, endTS := QuarterStartEndTimestamp()
	assert.Greater(t, startTS, int64(0))
	assert.Greater(t, endTS, startTS)

	startDate, endDate := QuarterStartEndDate()
	assert.NotEmpty(t, startDate)
	assert.NotEmpty(t, endDate)

	// 测试Q2
	startQ2, endQ2 := QuarterStartEnd("2024-05-15 12:00:00")
	assert.Contains(t, startQ2, "2024-04-01")
	assert.Contains(t, endQ2, "2024-06-30")
}

func TestMonthStartEnd(t *testing.T) {
	// 测试月份开始结束
	startStr, endStr := MonthStartEnd()
	assert.NotEmpty(t, startStr)
	assert.NotEmpty(t, endStr)

	startTS, endTS := MonthStartEndTimestamp()
	assert.Greater(t, startTS, int64(0))
	assert.Greater(t, endTS, startTS)

	startDate, endDate := MonthStartEndDate()
	assert.NotEmpty(t, startDate)
	assert.NotEmpty(t, endDate)

	// 测试特定月份
	startFeb, endFeb := MonthStartEnd("2024-02-15 12:00:00")
	assert.Contains(t, startFeb, "2024-02-01")
	assert.Contains(t, endFeb, "2024-02-29")
}

func TestWeekStartEnd(t *testing.T) {
	// 测试周开始结束
	startStr, endStr := WeekStartEnd()
	assert.NotEmpty(t, startStr)
	assert.NotEmpty(t, endStr)

	startTS, endTS := WeekStartEndTimestamp()
	assert.Greater(t, startTS, int64(0))
	assert.Greater(t, endTS, startTS)

	startDate, endDate := WeekStartEndDate()
	assert.NotEmpty(t, startDate)
	assert.NotEmpty(t, endDate)
}

func TestDayStartEnd(t *testing.T) {
	// 测试日开始结束
	startStr, endStr := DayStartEnd()
	assert.NotEmpty(t, startStr)
	assert.NotEmpty(t, endStr)
	assert.Contains(t, startStr, " 00:00:00")
	assert.Contains(t, endStr, " 23:59:59")

	startTS, endTS := DayStartEndTimestamp()
	assert.Greater(t, startTS, int64(0))
	assert.Greater(t, endTS, startTS)

	startDate, endDate := DayStartEndDate()
	assert.NotEmpty(t, startDate)
	assert.NotEmpty(t, endDate)
	assert.Equal(t, startDate, endDate) // 同一天

	// 测试特定日期
	start0101, end0101 := DayStartEnd("2024-01-01 15:30:00")
	assert.Equal(t, "2024-01-01 00:00:00", start0101)
	assert.Equal(t, "2024-01-01 23:59:59", end0101)
}

func TestHourStartEnd(t *testing.T) {
	// 测试小时开始结束
	startStr, endStr := HourStartEnd()
	assert.NotEmpty(t, startStr)
	assert.NotEmpty(t, endStr)
	assert.Contains(t, startStr, ":00:00")
	assert.Contains(t, endStr, ":59:59")

	startTS, endTS := HourStartEndTimestamp()
	assert.Greater(t, startTS, int64(0))
	assert.Greater(t, endTS, startTS)

	// 测试特定时间
	start15, end15 := HourStartEnd("2024-01-01 15:30:45")
	assert.Equal(t, "2024-01-01 15:00:00", start15)
	assert.Equal(t, "2024-01-01 15:59:59", end15)
}

func TestLunar(t *testing.T) {
	// 测试农历
	lunar := ToLunar()
	assert.NotEmpty(t, lunar)

	lunarDate := ToLunarDate()
	assert.NotEmpty(t, lunarDate)

	// 测试特定日期
	lunar20240101 := ToLunar("2024-01-01")
	assert.Contains(t, lunar20240101, "2023") // 2024年1月1日是农历2023年

	lunarDate20240210 := ToLunarDate("2024-02-10")
	assert.Contains(t, lunarDate20240210, "二零二四年正月初一") // 2024年春节

	// 测试生肖
	animal := LunarAnimal()
	assert.NotEmpty(t, animal)

	dragon := LunarAnimal("2024-02-10") // 龙年
	assert.Equal(t, "龙", dragon)

	rabbit := LunarAnimal("2023-02-10") // 兔年
	assert.Equal(t, "兔", rabbit)
}

func TestParseTimestamp(t *testing.T) {
	// 测试时间戳解析
	ts := int64(1704067200) // 2024-01-01 00:00:00
	parsed := ParseTimestamp(ts)
	assert.Contains(t, parsed, "2024-01-01")

	parsedDate := ParseTimestampDate(ts)
	assert.Equal(t, "2024-01-01", parsedDate)

	// 测试毫秒时间戳
	tsMilli := ts * 1000
	parsedMilli := ParseTimestampMilli(tsMilli)
	assert.Contains(t, parsedMilli, "2024-01-01")

	parsedMilliDate := ParseTimestampMilliDate(tsMilli)
	assert.Equal(t, "2024-01-01", parsedMilliDate)
}

func TestParseString(t *testing.T) {
	// 测试字符串解析
	parsed := ParseString("2024-01-01 15:30:45")
	assert.Equal(t, "2024-01-01 15:30:45", parsed)

	parsedDate := ParseStringDate("2024-01-01 15:30:45")
	assert.Equal(t, "2024-01-01", parsedDate)

	// 测试其他格式
	parsed2 := ParseString("2024/01/01")
	assert.Contains(t, parsed2, "2024-01-01")

	parsedDate2 := ParseStringDate("20240101")
	assert.Contains(t, parsedDate2, "2024-01-01")
}

func TestStartEndFunctionsWithParam(t *testing.T) {
	// 测试所有 StartEnd 函数传入参数的情况
	testDate := "2024-06-15 12:30:45"

	// YearStartEndTimestamp with param
	yearStartTS, yearEndTS := YearStartEndTimestamp(testDate)
	assert.Greater(t, yearStartTS, int64(0))
	assert.Greater(t, yearEndTS, yearStartTS)

	// YearStartEndDate with param
	yearStartDate, yearEndDate := YearStartEndDate(testDate)
	assert.Equal(t, "2024-01-01", yearStartDate)
	assert.Equal(t, "2024-12-31", yearEndDate)

	// QuarterStartEndTimestamp with param
	qStartTS, qEndTS := QuarterStartEndTimestamp(testDate)
	assert.Greater(t, qStartTS, int64(0))
	assert.Greater(t, qEndTS, qStartTS)

	// QuarterStartEndDate with param
	qStartDate, qEndDate := QuarterStartEndDate(testDate)
	assert.Equal(t, "2024-04-01", qStartDate)
	assert.Equal(t, "2024-06-30", qEndDate)

	// MonthStartEndTimestamp with param
	mStartTS, mEndTS := MonthStartEndTimestamp(testDate)
	assert.Greater(t, mStartTS, int64(0))
	assert.Greater(t, mEndTS, mStartTS)

	// MonthStartEndDate with param
	mStartDate, mEndDate := MonthStartEndDate(testDate)
	assert.Equal(t, "2024-06-01", mStartDate)
	assert.Equal(t, "2024-06-30", mEndDate)

	// WeekStartEnd with param
	wStart, wEnd := WeekStartEnd(testDate)
	assert.Contains(t, wStart, "2024-06-10") // 2024-06-15是周六，周一是06-10
	assert.Contains(t, wEnd, "2024-06-16")   // 周日是06-16

	// WeekStartEndTimestamp with param
	wStartTS, wEndTS := WeekStartEndTimestamp(testDate)
	assert.Greater(t, wStartTS, int64(0))
	assert.Greater(t, wEndTS, wStartTS)

	// WeekStartEndDate with param
	wStartDate, wEndDate := WeekStartEndDate(testDate)
	assert.Equal(t, "2024-06-10", wStartDate)
	assert.Equal(t, "2024-06-16", wEndDate)

	// DayStartEndTimestamp with param
	dStartTS, dEndTS := DayStartEndTimestamp(testDate)
	assert.Greater(t, dStartTS, int64(0))
	assert.Greater(t, dEndTS, dStartTS)

	// DayStartEndDate with param
	dStartDate, dEndDate := DayStartEndDate(testDate)
	assert.Equal(t, "2024-06-15", dStartDate)
	assert.Equal(t, "2024-06-15", dEndDate)

	// HourStartEndTimestamp with param
	hStartTS, hEndTS := HourStartEndTimestamp(testDate)
	assert.Greater(t, hStartTS, int64(0))
	assert.Greater(t, hEndTS, hStartTS)
}

func BenchmarkNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Now()
	}
}

func BenchmarkNowTimestamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NowTimestamp()
	}
}

func BenchmarkParseTimestamp(b *testing.B) {
	ts := int64(1704067200)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ParseTimestamp(ts)
	}
}

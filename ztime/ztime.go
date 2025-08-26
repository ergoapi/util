package ztime

import (
	"github.com/dromara/carbon/v2"
)

func init() {
	carbon.SetTimezone("Asia/Shanghai")
}

// YesterdayNow 昨天此刻
func YesterdayNow() string {
	return carbon.Yesterday().ToDateTimeString()
}

// YesterdayDate 昨天日期
func YesterdayDate() string {
	return carbon.Yesterday().ToDateString()
}

// YesterdayTimestamp 昨天时间戳秒
func YesterdayTimestamp() int64 {
	return carbon.Yesterday().Timestamp()
}

// Now 当前时间
func Now() string {
	return carbon.Now().ToDateTimeString()
}

// NowDate 今天日期
func NowDate() string {
	return carbon.Now().ToDateString()
}

// NowTimestamp 当前时间戳秒
func NowTimestamp() int64 {
	return carbon.Now().Timestamp()
}

// TomorrowNow 明天此刻
func TomorrowNow() string {
	return carbon.Tomorrow().ToDateTimeString()
}

// TomorrowDate 明天日期
func TomorrowDate() string {
	return carbon.Tomorrow().ToDateString()
}

// TomorrowTimestamp 明天时间戳秒
func TomorrowTimestamp() int64 {
	return carbon.Tomorrow().Timestamp()
}

// Minute 当前时间minutes分钟数后
func Minute(minutes int) string {
	if minutes >= 0 {
		return carbon.Now().AddMinutes(minutes).ToDateTimeString()
	}
	// 负数则减去minutes分钟
	minutes = -minutes
	return carbon.Now().SubMinutes(minutes).ToDateTimeString()
}

// MinuteTimestamp 当前时间minutes分钟数后时间戳秒
func MinuteTimestamp(minutes int) int64 {
	if minutes >= 0 {
		return carbon.Now().AddMinutes(minutes).Timestamp()
	}
	// 负数则减去minutes分钟
	minutes = -minutes
	return carbon.Now().SubMinutes(minutes).Timestamp()
}

// Hour 当前时间hours小时数后
func Hour(hours int) string {
	if hours >= 0 {
		return carbon.Now().AddHours(hours).ToDateTimeString()
	}
	// 负数则减去hours小时
	hours = -hours
	return carbon.Now().SubHours(hours).ToDateTimeString()
}

// HourTimestamp 当前时间hours小时数后时间戳秒
func HourTimestamp(hours int) int64 {
	if hours >= 0 {
		return carbon.Now().AddHours(hours).Timestamp()
	}
	// 负数则减去hours小时
	hours = -hours
	return carbon.Now().SubHours(hours).Timestamp()
}

// HourDate 当前时间hours小时数后日期
func HourDate(hours int) string {
	if hours >= 0 {
		return carbon.Now().AddHours(hours).ToDateString()
	}
	// 负数则减去hours小时
	hours = -hours
	return carbon.Now().SubHours(hours).ToDateString()
}

// Day 当前时间day天数后
func Day(day int) string {
	if day >= 0 {
		return carbon.Now().AddDays(day).ToDateTimeString()
	}
	// 负数则减去days天
	day = -day
	return carbon.Now().SubDays(day).ToDateTimeString()
}

// DayDate 当前时间day天数后日期
func DayDate(day int) string {
	if day >= 0 {
		return carbon.Now().AddDays(day).ToDateString()
	}
	// 负数则减去days天
	day = -day
	return carbon.Now().SubDays(day).ToDateString()
}

// Month 当前时间month月数后
func Month(month int) string {
	if month >= 0 {
		return carbon.Now().AddMonthsNoOverflow(month).ToDateTimeString()
	}
	// 负数则减去months月
	month = -month
	return carbon.Now().SubMonthsNoOverflow(month).ToDateTimeString()
}

// MonthDate 当前时间month月数后日期
func MonthDate(month int) string {
	if month >= 0 {
		return carbon.Now().AddMonthsNoOverflow(month).ToDateString()
	}
	// 负数则减去months月
	month = -month
	return carbon.Now().SubMonthsNoOverflow(month).ToDateString()
}

// DaysInYear 获取指定时间所在年份的天数, 如果未指定时间, 则获取当前时间所在年份的天数
func DaysInYear(t ...string) int {
	if len(t) == 0 {
		return carbon.Now().DaysInYear()
	}
	// t 格式为 2019-08-05 13:14:15
	return carbon.Parse(t[0]).DaysInYear()
}

// DaysInMonth 获取指定时间所在月份的天数, 如果未指定时间, 则获取当前时间所在月份的天数
func DaysInMonth(t ...string) int {
	if len(t) == 0 {
		return carbon.Now().DaysInMonth()
	}
	// t 格式为 2019-08-05 13:14:15
	return carbon.Parse(t[0]).DaysInMonth()
}

// Age 获取年龄
func Age(t string) int {
	return carbon.Parse(t).Age()
}

// Season 获取当前季节
func Season() string {
	return carbon.Now().SetLocale("en").Season()
}

// Constellation 获取星座
func Constellation(t ...string) string {
	if len(t) == 0 {
		return carbon.Now().SetLocale("en").Constellation()
	}
	// t 格式为 2019-08-05 13:14:15
	return carbon.Parse(t[0]).SetLocale("en").Constellation()
}

// WeekOfYear 获取当前时间所在周数
func WeekOfYear() int {
	return carbon.Now().WeekOfYear()
}

// DayOfYear 获取当前时间所在年的第几天
func DayOfYear() int {
	return carbon.Now().DayOfYear()
}

// DayOfWeek 获取当前时间所在周的第几天
func DayOfWeek() int {
	return carbon.Now().DayOfWeek()
}

// YearStartEnd 获取年份开始和结束时间
func YearStartEnd(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfYear().ToDateTimeString(), c.EndOfYear().ToDateTimeString()
}

// YearStartEndTimestamp 获取年份开始和结束时间戳
func YearStartEndTimestamp(t ...string) (int64, int64) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfYear().Timestamp(), c.EndOfYear().Timestamp()
}

// YearStartEndDate 获取年份开始和结束日期
func YearStartEndDate(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfYear().ToDateString(), c.EndOfYear().ToDateString()
}

// QuarterStartEnd 获取当前时间所在季度的开始和结束时间
func QuarterStartEnd(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfQuarter().ToDateTimeString(), c.EndOfQuarter().ToDateTimeString()
}

// QuarterStartEndTimestamp 获取当前时间所在季度的开始和结束时间戳
func QuarterStartEndTimestamp(t ...string) (int64, int64) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfQuarter().Timestamp(), c.EndOfQuarter().Timestamp()
}

// QuarterStartEndDate 获取当前时间所在季度的开始和结束日期
func QuarterStartEndDate(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfQuarter().ToDateString(), c.EndOfQuarter().ToDateString()
}

// MonthStartEnd 获取月份开始和结束时间
func MonthStartEnd(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfMonth().ToDateTimeString(), c.EndOfMonth().ToDateTimeString()
}

// MonthStartEndTimestamp 获取月份开始和结束时间戳
func MonthStartEndTimestamp(t ...string) (int64, int64) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfMonth().Timestamp(), c.EndOfMonth().Timestamp()
}

// MonthStartEndDate 获取月份开始和结束日期
func MonthStartEndDate(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfMonth().ToDateString(), c.EndOfMonth().ToDateString()
}

// WeekStartEnd 获取周开始和结束时间
func WeekStartEnd(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfWeek().ToDateTimeString(), c.EndOfWeek().ToDateTimeString()
}

// WeekStartEndTimestamp 获取周开始和结束时间戳
func WeekStartEndTimestamp(t ...string) (int64, int64) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfWeek().Timestamp(), c.EndOfWeek().Timestamp()
}

// WeekStartEndDate 获取周开始和结束日期
func WeekStartEndDate(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfWeek().ToDateString(), c.EndOfWeek().ToDateString()
}

// DayStartEnd 获取日开始和结束时间
func DayStartEnd(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfDay().ToDateTimeString(), c.EndOfDay().ToDateTimeString()
}

// DayStartEndTimestamp 获取日开始和结束时间戳
func DayStartEndTimestamp(t ...string) (int64, int64) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfDay().Timestamp(), c.EndOfDay().Timestamp()
}

// DayStartEndDate 获取日开始和结束日期
func DayStartEndDate(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfDay().ToDateString(), c.EndOfDay().ToDateString()
}

// HourStartEnd 获取小时开始和结束时间
func HourStartEnd(t ...string) (string, string) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfHour().ToDateTimeString(), c.EndOfHour().ToDateTimeString()
}

// HourStartEndTimestamp 获取小时开始和结束时间戳
func HourStartEndTimestamp(t ...string) (int64, int64) {
	var c *carbon.Carbon
	if len(t) == 0 {
		c = carbon.Now()
	} else {
		c = carbon.Parse(t[0])
	}
	return c.StartOfHour().Timestamp(), c.EndOfHour().Timestamp()
}

// ToLunar 转换为农历日期 2020-06-16
func ToLunar(t ...string) string {
	if len(t) == 0 {
		return carbon.Now().Lunar().String()
	}
	// t 格式为 2020-08-05
	return carbon.Parse(t[0]).Lunar().String()
}

// ToLunarDate 转换为农历日期 二零二零年六月十六
func ToLunarDate(t ...string) string {
	if len(t) == 0 {
		return carbon.Now().Lunar().ToDateString()
	}
	// t 格式为 2020-08-05
	return carbon.Parse(t[0]).Lunar().ToDateString()
}

// LunarAnimal 获取农历生肖
func LunarAnimal(t ...string) string {
	if len(t) == 0 {
		return carbon.Now().Lunar().Animal()
	}
	// t 格式为 2020-08-05
	return carbon.Parse(t[0]).Lunar().Animal()
}

// ParseTimestamp 解析时间戳为时间字符串
func ParseTimestamp(t int64) string {
	return carbon.CreateFromTimestamp(t).ToDateTimeString()
}

// ParseTimestampDate 解析时间戳为日期字符串
func ParseTimestampDate(t int64) string {
	return carbon.CreateFromTimestamp(t).ToDateString()
}

// ParseTimestampMilli 解析毫秒时间戳为时间字符串
func ParseTimestampMilli(t int64) string {
	return carbon.CreateFromTimestampMilli(t).ToDateTimeString()
}

// ParseTimestampMilliDate 解析毫秒时间戳为日期字符串
func ParseTimestampMilliDate(t int64) string {
	return carbon.CreateFromTimestampMilli(t).ToDateString()
}

// ParseString 解析时间字符串为时间字符串
func ParseString(t string) string {
	return carbon.Parse(t).ToDateTimeString()
}

// ParseStringDate 解析时间字符串为日期字符串
func ParseStringDate(t string) string {
	return carbon.Parse(t).ToDateString()
}

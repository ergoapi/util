package ztime

import (
	"strings"
	"testing"

	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/assert"
)

func TestRelativeOffsets(t *testing.T) {
	// Minute
	assert.Equal(t, carbon.Now().AddMinutes(5).ToDateTimeString(), Minute(5))
	assert.Equal(t, carbon.Now().SubMinutes(7).ToDateTimeString(), Minute(-7))
	assert.Equal(t, carbon.Now().AddMinutes(3).Timestamp(), MinuteTimestamp(3))

	// Hour
	assert.Equal(t, carbon.Now().AddHours(2).ToDateTimeString(), Hour(2))
	assert.Equal(t, carbon.Now().SubHours(2).ToDateTimeString(), Hour(-2))
	assert.Equal(t, carbon.Now().AddHours(4).Timestamp(), HourTimestamp(4))
	assert.Equal(t, carbon.Now().AddHours(1).ToDateString(), HourDate(1))

	// Day
	assert.Equal(t, carbon.Now().AddDays(1).ToDateTimeString(), Day(1))
	assert.Equal(t, carbon.Now().SubDays(1).ToDateTimeString(), Day(-1))
	assert.Equal(t, carbon.Now().AddDays(2).ToDateString(), DayDate(2))

	// Month (no overflow)
	assert.Equal(t, carbon.Now().AddMonthsNoOverflow(1).ToDateTimeString(), Month(1))
	assert.Equal(t, carbon.Now().SubMonthsNoOverflow(1).ToDateTimeString(), Month(-1))
	assert.Equal(t, carbon.Now().AddMonthsNoOverflow(2).ToDateString(), MonthDate(2))
}

func TestCountsAndPositions(t *testing.T) {
	// Year/Month days with explicit time
	assert.Equal(t, 366, DaysInYear("2020-08-05 13:14:15"))
	assert.Equal(t, 29, DaysInMonth("2024-02-10 12:00:00"))

	// Compare to carbon for now (non-deterministic but consistent with implementation)
	assert.Equal(t, carbon.Now().DaysInYear(), DaysInYear())
	assert.Equal(t, carbon.Now().DaysInMonth(), DaysInMonth())

	// Week/Day positions
	assert.Equal(t, carbon.Now().WeekOfYear(), WeekOfYear())
	assert.Equal(t, carbon.Now().DayOfYear(), DayOfYear())
	assert.Equal(t, carbon.Now().DayOfWeek(), DayOfWeek())
}

func TestLunarCountsAndData(t *testing.T) {
	// Test lunar days in year
	assert.Greater(t, LunarDaysInYear("2020-01-01"), 350)
	assert.Less(t, LunarDaysInYear("2020-01-01"), 390)
	assert.Greater(t, LunarDaysInYear(), 350) // Current year

	// Test lunar days in month
	assert.Greater(t, LunarDaysInMonth("2020-01-01"), 25)
	assert.Less(t, LunarDaysInMonth("2020-01-01"), 32)
	assert.Greater(t, LunarDaysInMonth(), 25) // Current month

	// Test lunar without parameters (default to current time)
	assert.NotEmpty(t, ToLunar())
	assert.NotEmpty(t, ToLunarDate())
	assert.NotEmpty(t, LunarAnimal())
}

func TestAgeSeasonConstellation(t *testing.T) {
	// Age (relative to now, compare with carbon for stability)
	birth := "1990-06-15 00:00:00"
	assert.Equal(t, carbon.Parse(birth).Age(), Age(birth))

	// Season & Constellation (locale en)
	assert.Equal(t, carbon.Now().SetLocale("en").Season(), Season())
	some := "2019-08-05 13:14:15"
	assert.Equal(t, carbon.Parse(some).SetLocale("en").Constellation(), Constellation(some))
	assert.Equal(t, carbon.Now().SetLocale("en").Constellation(), Constellation())
}

func TestStartEndRanges(t *testing.T) {
	base := "2023-05-10 11:22:33"
	c := carbon.Parse(base)

	ys, ye := YearStartEnd(base)
	assert.Equal(t, c.StartOfYear().ToDateTimeString(), ys)
	assert.Equal(t, c.EndOfYear().ToDateTimeString(), ye)
	yst, yet := YearStartEndTimestamp(base)
	assert.Equal(t, c.StartOfYear().Timestamp(), yst)
	assert.Equal(t, c.EndOfYear().Timestamp(), yet)
	ysd, yed := YearStartEndDate(base)
	assert.Equal(t, c.StartOfYear().ToDateString(), ysd)
	assert.Equal(t, c.EndOfYear().ToDateString(), yed)

	qs, qe := QuarterStartEnd(base)
	assert.Equal(t, c.StartOfQuarter().ToDateTimeString(), qs)
	assert.Equal(t, c.EndOfQuarter().ToDateTimeString(), qe)
	qst, qet := QuarterStartEndTimestamp(base)
	assert.Equal(t, c.StartOfQuarter().Timestamp(), qst)
	assert.Equal(t, c.EndOfQuarter().Timestamp(), qet)
	qsd, qed := QuarterStartEndDate(base)
	assert.Equal(t, c.StartOfQuarter().ToDateString(), qsd)
	assert.Equal(t, c.EndOfQuarter().ToDateString(), qed)

	ms, me := MonthStartEnd(base)
	assert.Equal(t, c.StartOfMonth().ToDateTimeString(), ms)
	assert.Equal(t, c.EndOfMonth().ToDateTimeString(), me)
	mst, met := MonthStartEndTimestamp(base)
	assert.Equal(t, c.StartOfMonth().Timestamp(), mst)
	assert.Equal(t, c.EndOfMonth().Timestamp(), met)
	msd, med := MonthStartEndDate(base)
	assert.Equal(t, c.StartOfMonth().ToDateString(), msd)
	assert.Equal(t, c.EndOfMonth().ToDateString(), med)

	ws, we := WeekStartEnd(base)
	assert.Equal(t, c.StartOfWeek().ToDateTimeString(), ws)
	assert.Equal(t, c.EndOfWeek().ToDateTimeString(), we)
	wst, wet := WeekStartEndTimestamp(base)
	assert.Equal(t, c.StartOfWeek().Timestamp(), wst)
	assert.Equal(t, c.EndOfWeek().Timestamp(), wet)
	wsd, wed := WeekStartEndDate(base)
	assert.Equal(t, c.StartOfWeek().ToDateString(), wsd)
	assert.Equal(t, c.EndOfWeek().ToDateString(), wed)

	ds, de := DayStartEnd(base)
	assert.Equal(t, c.StartOfDay().ToDateTimeString(), ds)
	assert.Equal(t, c.EndOfDay().ToDateTimeString(), de)
	dst, det := DayStartEndTimestamp(base)
	assert.Equal(t, c.StartOfDay().Timestamp(), dst)
	assert.Equal(t, c.EndOfDay().Timestamp(), det)
	dsd, ded := DayStartEndDate(base)
	assert.Equal(t, c.StartOfDay().ToDateString(), dsd)
	assert.Equal(t, c.EndOfDay().ToDateString(), ded)

	hs, he := HourStartEnd(base)
	assert.Equal(t, c.StartOfHour().ToDateTimeString(), hs)
	assert.Equal(t, c.EndOfHour().ToDateTimeString(), he)
	hst, het := HourStartEndTimestamp(base)
	assert.Equal(t, c.StartOfHour().Timestamp(), hst)
	assert.Equal(t, c.EndOfHour().Timestamp(), het)
}

func TestBasicTimeFunctions(t *testing.T) {
	// Yesterday functions
	assert.Equal(t, carbon.Yesterday().ToDateTimeString(), YesterdayNow())
	assert.Equal(t, carbon.Yesterday().ToDateString(), YesterdayDate())
	assert.Equal(t, carbon.Yesterday().Timestamp(), YesterdayTimestamp())

	// Now functions
	assert.Equal(t, carbon.Now().ToDateTimeString(), Now())
	assert.Equal(t, carbon.Now().ToDateString(), NowDate())
	assert.Equal(t, carbon.Now().Timestamp(), NowTimestamp())

	// Tomorrow functions
	assert.Equal(t, carbon.Tomorrow().ToDateTimeString(), TomorrowNow())
	assert.Equal(t, carbon.Tomorrow().ToDateString(), TomorrowDate())
	assert.Equal(t, carbon.Tomorrow().Timestamp(), TomorrowTimestamp())
}

func TestLunarAndParsing(t *testing.T) {
	d := "2020-08-05"
	assert.Equal(t, carbon.Parse(d).Lunar().String(), ToLunar(d))
	assert.Equal(t, carbon.Parse(d).Lunar().ToDateString(), ToLunarDate(d))
	assert.Equal(t, carbon.Parse(d).Lunar().Animal(), LunarAnimal(d))

	// Timestamp parsing
	c := carbon.Parse("2024-01-02 03:04:05")
	ts := c.Timestamp()
	tms := c.TimestampMilli()
	assert.Equal(t, c.ToDateTimeString(), ParseTimestamp(ts))
	assert.Equal(t, c.ToDateString(), ParseTimestampDate(ts))
	assert.Equal(t, c.ToDateTimeString(), ParseTimestampMilli(tms))
	assert.Equal(t, c.ToDateString(), ParseTimestampMilliDate(tms))

	// String parsing
	s := "2024-01-02 03:04:05"
	assert.Equal(t, carbon.Parse(s).ToDateTimeString(), ParseString(s))
	assert.Equal(t, carbon.Parse(s).ToDateString(), ParseStringDate(s))
}

func TestTraditionalCulture(t *testing.T) {
	// Test current shichen (时辰)
	shichen := NowShiChen()
	assert.NotEmpty(t, shichen)
	// Shichen should contain one of the 12 traditional time periods
	// 时辰 format is like "甲子时" or similar, so we check if it contains time characters
	found := false
	timeChars := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	for _, char := range timeChars {
		if strings.Contains(shichen, char) {
			found = true
			break
		}
	}
	assert.True(t, found, "Shichen should contain a valid time character: %s", shichen)

	// Test festivals for Spring Festival 2024 (2024-02-10)
	festivals := DayFestivals("2024-02-10")
	assert.Contains(t, festivals, "春节")

	// Test festivals with no parameters (current date)
	currentFestivals := DayFestivals()
	assert.NotNil(t, currentFestivals)

	// Test YiGi (宜忌) for a specific date
	yigi := DayYiGi("2024-01-01")
	assert.Equal(t, "2024-01-01", yigi.Date)
	assert.NotNil(t, yigi.Yi)
	assert.NotNil(t, yigi.Gi)

	// Test YiGi with no parameters (current date)
	currentYiGi := DayYiGi()
	assert.NotEmpty(t, currentYiGi.Date)
	assert.NotNil(t, currentYiGi.Yi)
	assert.NotNil(t, currentYiGi.Gi)
}

func TestBoundaryAndEdgeCases(t *testing.T) {
	// Test leap year
	assert.Equal(t, 366, DaysInYear("2024-01-01"))
	assert.Equal(t, 365, DaysInYear("2023-01-01"))

	// Test February in leap year
	assert.Equal(t, 29, DaysInMonth("2024-02-01"))
	assert.Equal(t, 28, DaysInMonth("2023-02-01"))

	// Test zero values
	assert.Equal(t, carbon.Now().ToDateTimeString(), Minute(0))
	assert.Equal(t, carbon.Now().ToDateTimeString(), Hour(0))
	assert.Equal(t, carbon.Now().ToDateTimeString(), Day(0))
	assert.Equal(t, carbon.Now().ToDateTimeString(), Month(0))

	// Test large values
	result := Minute(1440) // 24 hours in minutes
	expected := carbon.Now().AddMinutes(1440).ToDateTimeString()
	assert.Equal(t, expected, result)

	result = Hour(24)
	expected = carbon.Now().AddHours(24).ToDateTimeString()
	assert.Equal(t, expected, result)

	// Test negative values
	assert.Equal(t, carbon.Now().SubMinutes(60).ToDateTimeString(), Minute(-60))
	assert.Equal(t, carbon.Now().SubHours(12).ToDateTimeString(), Hour(-12))
	assert.Equal(t, carbon.Now().SubDays(7).ToDateTimeString(), Day(-7))

	// Test extreme timestamp values
	ts := int64(0)                                             // Unix epoch
	assert.Equal(t, "1970-01-01 08:00:00", ParseTimestamp(ts)) // Asia/Shanghai timezone

	tms := int64(0) // Millisecond epoch
	assert.Equal(t, "1970-01-01 08:00:00", ParseTimestampMilli(tms))

	// Test invalid date strings (should not panic, carbon handles gracefully)
	assert.NotPanics(t, func() {
		ParseString("invalid-date")
	})

	// Test month overflow behavior
	// AddMonthsNoOverflow should handle overflow scenarios gracefully
	result = Month(1) // This tests current behavior
	assert.NotEmpty(t, result)

	// Test end of year/month boundaries
	dec31 := "2023-12-31 23:59:59"
	yearStart, yearEnd := YearStartEnd(dec31)
	assert.Equal(t, "2023-01-01 00:00:00", yearStart)
	assert.Equal(t, "2023-12-31 23:59:59", yearEnd)
}

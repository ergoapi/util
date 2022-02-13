package ztime

import (
	"time"

	"github.com/ergoapi/util/common"
)

// https://github.com/golang-module/carbon 参考

type ZTime struct {
	time time.Time
	loc  *time.Location
}

func NewZTime() ZTime {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return ZTime{loc: loc}
}

// Time2ZTime time.Time 转化成Ztime
func Time2ZTime(tt time.Time) ZTime {
	zt := NewZTime()
	zt.time = tt
	return zt
}

func (zt ZTime) ZTime2Time() time.Time {
	return zt.time.In(zt.loc)
}

func (zt ZTime) Now() ZTime {
	zt.time = time.Now().In(zt.loc)
	return zt
}

// Now 当前时间
func Now() ZTime {
	return NewZTime().Now()
}

func (zt ZTime) Tomorrow() ZTime {
	if zt.IsZero() {
		zt.time = time.Now().In(zt.loc).AddDate(0, 0, 1)
	} else {
		zt.time = zt.time.In(zt.loc).AddDate(0, 0, 1)
	}
	return zt
}

// Tomorrow
func Tomorrow() ZTime {
	return NewZTime().Tomorrow()
}

func (zt ZTime) Yesterday() ZTime {
	if zt.IsZero() {
		zt.time = time.Now().In(zt.loc).AddDate(0, 0, -1)
	} else {
		zt.time = zt.time.In(zt.loc).AddDate(0, 0, -1)
	}
	return zt
}

// Yesterday 昨天
func Yesterday() ZTime {
	return NewZTime().Yesterday()
}

// IsZero 是否是零值时间
func (zt ZTime) IsZero() bool {
	return zt.time.IsZero()
}

// 是否是无效时间
func (zt ZTime) IsInvalid() bool {
	if zt.IsZero() {
		return true
	}
	return false
}

// 是否是闰年
func (zt ZTime) IsLeapYear() bool {
	if zt.IsInvalid() {
		return false
	}
	year := zt.time.Year()
	if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
		return true
	}
	return false
}

func (zt ZTime) Weekday() time.Weekday {
	return zt.time.Weekday()
}

// IsSaturday 是否是周六
func (zt ZTime) IsSaturday() bool {
	if zt.IsInvalid() {
		return false
	}
	return zt.time.In(zt.loc).Weekday() == time.Saturday
}

// IsSunday 是否是周日
func (zt ZTime) IsSunday() bool {
	if zt.IsInvalid() {
		return false
	}
	return zt.time.In(zt.loc).Weekday() == time.Sunday
}

// IsWeekday 是否是工作日
func (zt ZTime) IsWeekday() bool {
	if zt.IsInvalid() {
		return false
	}
	return !zt.IsSaturday() && !zt.IsSunday()
}

// IsWeekend reports whether is weekend.
// 是否是周末
func (zt ZTime) IsWeekend() bool {
	if zt.IsInvalid() {
		return false
	}
	return zt.IsSaturday() || zt.IsSunday()
}

func (zt ZTime) NeedWork() bool {
	today := zt.DateTimeLayout()
	h := HolidayGet(today)
	return h.NeedWork
}

// NeedWork 需要工作
func NeedWork() bool {
	today := NewZTime().Now().DateTimeLayout()
	h := HolidayGet(today)
	return h.NeedWork
}

// DefaultTimeLayout 输出日期字符串
func (zt ZTime) DefaultTimeLayout() string {
	if zt.IsInvalid() {
		return ""
	}
	return zt.time.In(zt.loc).Format(common.DefaultTimeLayout)
}

// DateTimeLayout 输出日期字符串
func (zt ZTime) DateTimeLayout() string {
	if zt.IsInvalid() {
		return ""
	}
	return zt.time.In(zt.loc).Format(common.DateTimeLayout)
}

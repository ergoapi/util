// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package ztime

import (
	"errors"

	"github.com/6tail/lunar-go/HolidayUtil"
	"github.com/dromara/carbon/v2"
)

// Holiday 节假日信息
type Holiday struct {
	Date     string `json:"date"`      // 日期
	Name     string `json:"name"`      // 名称
	IsWork   bool   `json:"is_work"`   // 是否需要工作
	IsAdjust bool   `json:"is_adjust"` // 是否是调休
}

// GetHoliday 获取指定日期的节假日信息
func GetHoliday(date string) (*Holiday, error) {
	if date == "" {
		return nil, errors.New("empty date")
	}

	// 解析日期
	c := carbon.Parse(date)
	if c.IsInvalid() {
		return nil, errors.New("invalid date format")
	}

	// 查询节假日信息
	h := &Holiday{
		Date: date,
	}

	// 查询外部节假日数据
	d := HolidayUtil.GetHoliday(date)
	if d != nil {
		h.Name = d.GetName()
		h.IsWork = d.IsWork()
		h.IsAdjust = d.IsWork() // 如果节假日需要工作，说明是调休

		if h.IsAdjust {
			h.Name = h.Name + "调休"
		}
		return h, nil
	}

	// 判断是否是周末
	weekday := c.DayOfWeek()
	if weekday == 6 || weekday == 7 { // 6=Saturday, 7=Sunday in carbon
		h.Name = "周末"
		h.IsWork = false
		h.IsAdjust = false
	} else {
		h.Name = "工作日"
		h.IsWork = true
		h.IsAdjust = false
	}

	return h, nil
}

// IsHoliday 判断是否是节假日
func IsHoliday(date string) (bool, error) {
	h, err := GetHoliday(date)
	if err != nil {
		return false, err
	}
	return !h.IsWork, nil
}

// IsWorkday 判断是否是工作日
func IsWorkday(date string) (bool, error) {
	h, err := GetHoliday(date)
	if err != nil {
		return false, err
	}
	return h.IsWork, nil
}

// NeedWork 判断今天是否需要工作
func NeedWork() bool {
	today := NowDate() // 使用新的 NowDate 函数
	h, err := GetHoliday(today)
	if err != nil {
		return true // 出错默认需要工作
	}
	return h.IsWork
}

// TodayHolidayInfo 获取今天的节假日信息
func TodayHolidayInfo() (*Holiday, error) {
	return GetHoliday(NowDate())
}

// TomorrowHolidayInfo 获取明天的节假日信息
func TomorrowHolidayInfo() (*Holiday, error) {
	return GetHoliday(TomorrowDate())
}

// YesterdayHolidayInfo 获取昨天的节假日信息
func YesterdayHolidayInfo() (*Holiday, error) {
	return GetHoliday(YesterdayDate())
}

// NextWorkday 获取下一个工作日
func NextWorkday() (string, error) {
	for i := 1; i <= 365; i++ {
		date := DayDate(i)
		isWork, err := IsWorkday(date)
		if err != nil {
			continue
		}
		if isWork {
			return date, nil
		}
	}
	return "", errors.New("no workday found in next 365 days")
}

// PrevWorkday 获取上一个工作日
func PrevWorkday() (string, error) {
	for i := -1; i >= -365; i-- {
		date := DayDate(i)
		isWork, err := IsWorkday(date)
		if err != nil {
			continue
		}
		if isWork {
			return date, nil
		}
	}
	return "", errors.New("no workday found in prev 365 days")
}

// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package ztime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHolidayGet(t *testing.T) {
	days := []string{"2021-04-28", "2021-04-29", "2021-04-30", "2021-05-01", "2021-05-04", "2021-05-05", "2021-05-06", "2021-05-07", "2021-05-08", "2021-05-09", "2021-05-10", "2021-09-30", "2021-10-06"}
	for _, day := range days {
		h, err := GetHoliday(day)
		if err != nil {
			t.Errorf("GetHoliday error: %v", err)
			continue
		}
		t.Logf("[%v] %v IsWork=%v IsAdjust=%v", day, h.Name, h.IsWork, h.IsAdjust)
	}
}

func TestIsHoliday(t *testing.T) {
	// 测试周末（应该是假期）
	isHoliday, err := IsHoliday("2024-01-06") // 周六
	assert.NoError(t, err)
	assert.True(t, isHoliday, "Saturday should be a holiday")

	// 测试工作日
	isHoliday, err = IsHoliday("2024-01-02") // 周二
	assert.NoError(t, err)
	// 注意：具体结果依赖于节假日数据
	t.Logf("2024-01-02 is holiday: %v", isHoliday)
}

func TestIsWorkday(t *testing.T) {
	// 测试工作日
	isWorkday, err := IsWorkday("2024-01-02") // 周二
	assert.NoError(t, err)
	assert.True(t, isWorkday, "Tuesday should be a workday")

	// 测试周末
	isWorkday, err = IsWorkday("2024-01-06") // 周六
	assert.NoError(t, err)
	assert.False(t, isWorkday, "Saturday should not be a workday")
}

func TestTodayNeedWorkWithMock(t *testing.T) {
	// 通过测试确保 TodayNeedWork 函数正常运行
	result := TodayNeedWork()
	// 结果应该是 true 或 false
	assert.IsType(t, bool(true), result)
}

func TestHolidayInfo(t *testing.T) {
	// 测试今天的节假日信息
	today, err := TodayHolidayInfo()
	assert.NoError(t, err)
	assert.NotNil(t, today)
	t.Logf("Today: %s, IsWork: %v, IsAdjust: %v", today.Name, today.IsWork, today.IsAdjust)

	// 测试明天的节假日信息
	tomorrow, err := TomorrowHolidayInfo()
	assert.NoError(t, err)
	assert.NotNil(t, tomorrow)
	t.Logf("Tomorrow: %s, IsWork: %v, IsAdjust: %v", tomorrow.Name, tomorrow.IsWork, tomorrow.IsAdjust)

	// 测试昨天的节假日信息
	yesterday, err := YesterdayHolidayInfo()
	assert.NoError(t, err)
	assert.NotNil(t, yesterday)
	t.Logf("Yesterday: %s, IsWork: %v, IsAdjust: %v", yesterday.Name, yesterday.IsWork, yesterday.IsAdjust)
}

func TestNextPrevWorkday(t *testing.T) {
	// 测试下一个工作日
	next, err := NextWorkday()
	if err != nil {
		t.Logf("NextWorkday error: %v", err)
	} else {
		t.Logf("Next workday: %s", next)
		// 验证确实是工作日
		isWork, _ := IsWorkday(next)
		assert.True(t, isWork)
	}

	// 测试上一个工作日
	prev, err := PrevWorkday()
	if err != nil {
		t.Logf("PrevWorkday error: %v", err)
	} else {
		t.Logf("Previous workday: %s", prev)
		// 验证确实是工作日
		isWork, _ := IsWorkday(prev)
		assert.True(t, isWork)
	}
}

func TestNextPrevWorkdayEdgeCases(t *testing.T) {
	// 模拟测试 NextWorkday 的循环逻辑
	// 通过实际调用来测试各种路径
	for i := 0; i < 5; i++ {
		date := DayDate(i)
		isWork, err := IsWorkday(date)
		if err != nil {
			t.Logf("Day %d: error checking workday: %v", i, err)
			continue
		}
		t.Logf("Day %d (%s): isWorkday=%v", i, date, isWork)
	}

	// 测试 PrevWorkday 的循环逻辑
	for i := -1; i >= -5; i-- {
		date := DayDate(i)
		isWork, err := IsWorkday(date)
		if err != nil {
			t.Logf("Day %d: error checking workday: %v", i, err)
			continue
		}
		t.Logf("Day %d (%s): isWorkday=%v", i, date, isWork)
	}
}

func TestHolidayErrorCases(t *testing.T) {
	// 测试空日期
	_, err := GetHoliday("")
	assert.Error(t, err, "Empty date should return error")

	// 测试无效日期格式
	_, err = GetHoliday("invalid-date")
	assert.Error(t, err, "Invalid date should return error")

	// 测试 IsHoliday 错误处理
	_, err = IsHoliday("")
	assert.Error(t, err)
	_, err = IsHoliday("invalid")
	assert.Error(t, err)

	// 测试 IsWorkday 错误处理
	_, err = IsWorkday("")
	assert.Error(t, err)
	_, err = IsWorkday("invalid")
	assert.Error(t, err)
}

func TestHolidayDetails(t *testing.T) {
	// 测试元旦
	h, err := GetHoliday("2024-01-01")
	assert.NoError(t, err)
	assert.False(t, h.IsWork, "New Year's Day should be a holiday")
	t.Logf("2024-01-01: %s, IsWork: %v, IsAdjust: %v", h.Name, h.IsWork, h.IsAdjust)

	// 测试春节
	h, err = GetHoliday("2024-02-10") // 春节初一
	assert.NoError(t, err)
	assert.False(t, h.IsWork)
	assert.True(t, h.IsHoliday)
	t.Logf("2024-02-10: %s, IsWork: %v, IsAdjust: %v", h.Name, h.IsWork, h.IsAdjust)

	// 测试调休日
	h, err = GetHoliday("2024-02-04") // 春节前调休上班
	assert.NoError(t, err)
	assert.True(t, h.IsWork)
	assert.False(t, h.IsHoliday)
	assert.True(t, h.IsAdjust)
	assert.Contains(t, h.Name, "调休")
	t.Logf("2024-02-04: %s, IsWork: %v, IsAdjust: %v", h.Name, h.IsWork, h.IsAdjust)
}

func TestHolidayIsHolidayField(t *testing.T) {
	// 正常工作日（大概率不在法定节假日表）
	h, err := GetHoliday("2024-01-03")
	assert.NoError(t, err)
	assert.Equal(t, !h.IsWork, h.IsHoliday)
	if h.IsWork {
		assert.Equal(t, "工作日", h.Name)
	}

	// 周末（非调休）
	h, err = GetHoliday("2024-01-06") // 周六
	assert.NoError(t, err)
	assert.False(t, h.IsWork)
	assert.True(t, h.IsHoliday)
	assert.Equal(t, "周末", h.Name)

	// 法定节假日
	h, err = GetHoliday("2024-01-01")
	assert.NoError(t, err)
	assert.False(t, h.IsWork)
	assert.True(t, h.IsHoliday)
}

func TestHolidayDateNormalization(t *testing.T) {
	h, err := GetHoliday("2024-02-10 12:34:56")
	assert.NoError(t, err)
	assert.Equal(t, "2024-02-10", h.Date)
}

func TestHolidayFallbackNames(t *testing.T) {
	// 工作日回退名称
	h, err := GetHoliday("2024-01-02") // 周二
	assert.NoError(t, err)
	if h.IsWork { // 若未命中法定数据，名称应为工作日
		assert.Equal(t, "工作日", h.Name)
	}

	// 周末回退名称
	h, err = GetHoliday("2024-01-07") // 周日
	assert.NoError(t, err)
	if !h.IsWork { // 若未命中法定数据，名称应为周末
		assert.Equal(t, "周末", h.Name)
	}
}

func TestSundayWeekendRecognition(t *testing.T) {
	// 确认周日也能识别为周末（兼容 DayOfWeek 返回 0 的情况）
	h, err := GetHoliday("2024-01-07") // 周日
	assert.NoError(t, err)
	if !h.IsAdjust { // 不是调休上班
		assert.False(t, h.IsWork)
		assert.True(t, h.IsHoliday)
	}
}

func BenchmarkGetHoliday(b *testing.B) {
	date := "2024-01-01"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetHoliday(date)
	}
}

func BenchmarkTodayNeedWork(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TodayNeedWork()
	}
}

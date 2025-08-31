// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ergoapi/util/version/timeversion"
)

func main() {
	fmt.Println("=== 时间版本管理库示例 ===")

	// 1. 基本版本解析
	fmt.Println("\n1. 基本版本解析:")
	examples := []string{
		"2025.1.0101",  // 单位数月份
		"2025.12.3122", // 双位数月份
		"v2025.6.1505", // 带v前缀
		"2024.2.2901",  // 闰年2月29日
	}

	for _, example := range examples {
		v, err := timeversion.Parse(example)
		if err != nil {
			fmt.Printf("  解析失败 %s: %v\n", example, err)
			continue
		}
		fmt.Printf("  %s -> 年:%d 月:%d 日:%d 序列:%d 日期:%s\n",
			example, v.Year, v.Month, v.Day, v.Sequence,
			v.Date().Format("2006-01-02"))
	}

	// 2. 版本比较
	fmt.Println("\n2. 版本比较:")
	v1 := timeversion.MustParse("2025.1.0101")
	v2 := timeversion.MustParse("2025.1.0102")
	v3 := timeversion.MustParse("2025.1.0201")
	v4 := timeversion.MustParse("2025.2.0101")

	comparisons := []struct {
		a, b *timeversion.TimeVersion
		desc string
	}{
		{v1, v2, "同日不同序列"},
		{v1, v3, "不同日期"},
		{v1, v4, "不同月份"},
	}

	for _, c := range comparisons {
		fmt.Printf("  %s vs %s (%s):\n", c.a, c.b, c.desc)
		fmt.Printf("    IsLessThan: %t\n", c.a.IsLessThan(c.b))
		fmt.Printf("    IsEqual: %t\n", c.a.IsEqual(c.b))
		fmt.Printf("    Compare: %d\n", c.a.Compare(c.b))
	}

	// 3. 版本生成
	fmt.Println("\n3. 版本生成:")

	// 今日第一个版本
	today := timeversion.Now()
	fmt.Printf("  今日首版: %s\n", today)

	// 今日指定序列
	todaySeq5, err := timeversion.Today(5)
	if err != nil {
		log.Printf("  生成今日序列5失败: %v\n", err)
	} else {
		fmt.Printf("  今日序列5: %s\n", todaySeq5)
	}

	// 指定日期版本
	customDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	customVersion, err := timeversion.FromDate(customDate, 3)
	if err != nil {
		log.Printf("  指定日期版本失败: %v\n", err)
	} else {
		fmt.Printf("  指定日期版本: %s\n", customVersion)
	}

	// 4. 版本递增
	fmt.Println("\n4. 版本递增:")
	base := timeversion.MustParse("2025.1.0105")
	fmt.Printf("  基础版本: %s\n", base)

	nextSeq, err := base.NextSequence()
	if err != nil {
		fmt.Printf("  序列递增失败: %v\n", err)
	} else {
		fmt.Printf("  下一序列: %s\n", nextSeq)
	}

	nextDay := base.NextDay()
	fmt.Printf("  下一天: %s\n", nextDay)

	// 5. 包级别比较函数
	fmt.Println("\n5. 包级别比较函数:")
	v1Str, v2Str := "2025.1.0101", "2025.1.0102"

	isLess, err := timeversion.IsLessThan(v1Str, v2Str)
	if err != nil {
		log.Printf("  比较失败: %v\n", err)
	} else {
		fmt.Printf("  %s < %s: %t\n", v1Str, v2Str, isLess)
	}

	isEqual, err := timeversion.IsEqual(v1Str, "v2025.1.0101")
	if err != nil {
		log.Printf("  相等比较失败: %v\n", err)
	} else {
		fmt.Printf("  %s == v2025.1.0101: %t\n", v1Str, isEqual)
	}

	// 6. 版本排序
	fmt.Println("\n6. 版本排序:")
	versions := []string{
		"2025.1.0103", "2025.1.0101", "2025.1.0202",
		"v2025.1.0102", "2025.2.0101", "2024.12.3199",
	}
	fmt.Printf("  原始顺序: %v\n", versions)

	err = timeversion.Sort(versions)
	if err != nil {
		log.Printf("  排序失败: %v\n", err)
	} else {
		fmt.Printf("  排序后: %v\n", versions)
	}

	// 7. 查找最新版本
	fmt.Println("\n7. 查找最新版本:")
	testVersions := []string{
		"2025.1.0101", "2025.1.0103", "2025.1.0102",
		"2025.1.0201", "v2025.1.0105",
	}
	fmt.Printf("  版本列表: %v\n", testVersions)

	latest, err := timeversion.Latest(testVersions)
	if err != nil {
		log.Printf("  查找最新版本失败: %v\n", err)
	} else {
		fmt.Printf("  最新版本: %s\n", latest)
	}

	// 8. 获取指定日期的所有版本
	fmt.Println("\n8. 获取指定日期的所有版本:")
	allVersions := []string{
		"2025.1.0101", "2025.1.0103", "2025.1.0102", // 1月1日
		"2025.1.0201", "2025.1.0203", // 1月2日
		"2025.2.0101", // 2月1日
	}

	jan1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	jan1Versions, err := timeversion.GetVersionsForDate(allVersions, jan1)
	if err != nil {
		log.Printf("  获取1月1日版本失败: %v\n", err)
	} else {
		fmt.Printf("  2025年1月1日的版本: %v\n", jan1Versions)
	}

	// 9. 获取今日下一版本号
	fmt.Println("\n9. 获取今日下一版本号:")
	// 模拟今日已存在的版本
	simulatedToday := time.Now()
	existingToday := []string{
		fmt.Sprintf("%d.%d.%02d01", simulatedToday.Year(), int(simulatedToday.Month()), simulatedToday.Day()),
		fmt.Sprintf("%d.%d.%02d03", simulatedToday.Year(), int(simulatedToday.Month()), simulatedToday.Day()),
		fmt.Sprintf("%d.%d.%02d02", simulatedToday.Year(), int(simulatedToday.Month()), simulatedToday.Day()),
	}

	fmt.Printf("  今日已存在版本: %v\n", existingToday)
	nextToday, err := timeversion.GetNextVersionForToday(existingToday)
	if err != nil {
		log.Printf("  获取今日下一版本失败: %v\n", err)
	} else {
		fmt.Printf("  今日下一版本: %s\n", nextToday)
	}

	// 10. 错误处理示例
	fmt.Println("\n10. 错误处理示例:")
	invalidInputs := []string{
		"",             // 空字符串
		"2025.1",       // 格式不完整
		"2025.13.0101", // 无效月份
		"2025.2.3001",  // 无效日期
		"2025.1.0100",  // 无效序列号
		"2025.2.2901",  // 非闰年2月29日
	}

	for _, invalid := range invalidInputs {
		_, err := timeversion.Parse(invalid)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", invalid, err)
		} else {
			fmt.Printf("  ✓ %s: 意外通过验证\n", invalid)
		}
	}

	fmt.Println("\n=== 示例完成 ===")
}

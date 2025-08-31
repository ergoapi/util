// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/ergoapi/util/version/semver"
)

func main() {
	fmt.Println("=== 新版本库示例 ===")

	// 基本版本比较
	fmt.Println("\n1. 基本版本比较:")
	result, err := semver.IsLessThan("1.0.0", "1.0.1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1.0.0 < 1.0.1: %t\n", result)

	// 使用Version对象
	fmt.Println("\n2. 使用Version对象:")
	v1 := semver.MustParse("1.2.3-alpha.1")
	v2 := semver.MustParse("v1.2.3")

	fmt.Printf("%s vs %s\n", v1, v2)
	fmt.Printf("  IsLessThan: %t\n", v1.IsLessThan(v2))
	fmt.Printf("  IsEqual: %t\n", v1.IsEqual(v2))
	fmt.Printf("  Compare: %d\n", v1.Compare(v2))

	// 版本递增
	fmt.Println("\n3. 版本递增:")
	base := semver.MustParse("1.2.3-alpha.1+build.123")
	fmt.Printf("Base version: %s\n", base)
	fmt.Printf("  Next major: %s\n", base.IncrementMajor())
	fmt.Printf("  Next minor: %s\n", base.IncrementMinor())
	fmt.Printf("  Next patch: %s\n", base.IncrementPatch())

	// 版本属性
	fmt.Println("\n4. 版本属性:")
	complex := semver.MustParse("v2.1.3-alpha.1+build.123")
	fmt.Printf("Version: %s\n", complex)
	fmt.Printf("  Major: %d\n", complex.Major())
	fmt.Printf("  Minor: %d\n", complex.Minor())
	fmt.Printf("  Patch: %d\n", complex.Patch())
	fmt.Printf("  Pre-release: %v\n", complex.Pre())
	fmt.Printf("  Build: %v\n", complex.Build())

	// 版本排序
	fmt.Println("\n5. 版本排序:")
	versions := []string{"2.0.0", "1.0.0", "1.5.0", "v1.2.0", "2.0.0-alpha"}
	fmt.Printf("原始: %v\n", versions)

	err = semver.Sort(versions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("排序后: %v\n", versions)

	// 查找最新版本
	fmt.Println("\n6. 查找最新版本:")
	testVersions := []string{"1.0.0", "2.0.0", "1.5.0", "v1.2.0"}
	latest, err := semver.Latest(testVersions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("版本列表: %v\n", testVersions)
	fmt.Printf("最新版本: %s\n", latest)

	// 错误处理示例
	fmt.Println("\n7. 错误处理:")
	_, err = semver.Parse("invalid.version")
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
	}
}

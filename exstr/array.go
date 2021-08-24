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

package exstr

// DuplicateInt64ElementInt64 数组去重
func DuplicateInt64ElementInt64(addrs []int64) []int64 {
	result := make([]int64, 0, len(addrs))
	temp := map[int64]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// DuplicateIntElementInt 数组去重
func DuplicateIntElementInt(addrs []int) []int {
	result := make([]int, 0, len(addrs))
	temp := map[int]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// DuplicateStrElement 字符串去重
func DuplicateStrElement(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	temp := map[string]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// StringArrayContains 字符串数组是否包含某字符串
func StringArrayContains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

// Int64ArrayContains 数组是否包含某字符串
func Int64ArrayContains(addrs []int64, i int64) bool {
	for _, s := range addrs {
		if s == i {
			return true
		}
	}
	return false
}

// IntArrayContains 数组是否包含某字符串
func IntArrayContains(addrs []int, i int) bool {
	for _, s := range addrs {
		if s == i {
			return true
		}
	}
	return false
}

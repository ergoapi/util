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

import (
	"strings"

	"github.com/ergoapi/util/common"
)

// Blacklist
func Blacklist(s string) bool {
	if strings.Contains(s, "<") {
		return true
	}

	if strings.Contains(s, ">") {
		return true
	}

	if strings.Contains(s, "&") {
		return true
	}

	if strings.Contains(s, "'") {
		return true
	}

	if strings.Contains(s, "\"") {
		return true
	}

	if strings.Contains(s, "file://") {
		return true
	}

	if strings.Contains(s, "../") {
		return true
	}

	if strings.Contains(s, "%") {
		return true
	}

	if strings.Contains(s, "=") {
		return true
	}

	if strings.Contains(s, "--") {
		return true
	}

	return false
}

// KubeBlacklist
func KubeBlacklist(s string, extlist ...string) bool {
	if strings.HasPrefix(s, "kube-") {
		return true
	}

	for _, i := range extlist {
		if strings.Contains(s, i) {
			return true
		}
	}
	return false
}

// ToLower converts ascii string to lower-case
func ToLower(b string) string {
	res := make([]byte, len(b))
	copy(res, b)
	for i := 0; i < len(res); i++ {
		res[i] = common.ToLowerTable[res[i]]
	}

	return UnsafeString(res)
}

// ToUpper converts ascii string to upper-case
func ToUpper(b string) string {
	res := make([]byte, len(b))
	copy(res, b)
	for i := 0; i < len(res); i++ {
		res[i] = common.ToUpperTable[res[i]]
	}

	return UnsafeString(res)
}

// TrimLeft is the equivalent of strings.TrimLeft
func TrimLeft(s string, cutset byte) string {
	lenStr, start := len(s), 0
	for start < lenStr && s[start] == cutset {
		start++
	}
	return s[start:]
}

// Trim is the equivalent of strings.Trim
func Trim(s string, cutset byte) string {
	i, j := 0, len(s)-1
	for ; i <= j; i++ {
		if s[i] != cutset {
			break
		}
	}
	for ; i < j; j-- {
		if s[j] != cutset {
			break
		}
	}

	return s[i : j+1]
}

// TrimRight is the equivalent of strings.TrimRight
func TrimRight(s string, cutset byte) string {
	lenStr := len(s)
	for lenStr > 0 && s[lenStr-1] == cutset {
		lenStr--
	}
	return s[:lenStr]
}

// EqualFold tests ascii strings for equality case-insensitively
func EqualFold(b, s string) bool {
	if len(b) != len(s) {
		return false
	}
	for i := len(b) - 1; i >= 0; i-- {
		if common.ToUpperTable[b[i]] != common.ToUpperTable[s[i]] {
			return false
		}
	}
	return true
}

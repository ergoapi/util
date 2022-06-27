//  Copyright (c) 2020. The EFF Team Authors.
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
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

// Str2Int string to int
func Str2Int(s string) int {
	si, _ := strconv.Atoi(s)
	return si
}

// Str2Int32 string to int32
func Str2Int32(s string) int32 {
	si32, _ := strconv.ParseInt(s, 10, 32)
	return int32(si32)
}

// Str2Int64 string to int64
func Str2Int64(s string) int64 {
	si64, _ := strconv.ParseInt(s, 10, 64)
	return si64
}

// Str2UInt64 string to uint64
func Str2UInt64(s string) uint64 {
	si64, _ := strconv.ParseUint(s, 10, 64)
	return si64
}

// Str2Float64 string to float64
func Str2Float64(s string) float64 {
	sf, _ := strconv.ParseFloat(s, 64)
	return sf
}

// Str2Byte string to byte
func Str2Byte(s string) []byte {
	return []byte(s)
}

// Int642Str int64 to string
func Int642Str(i int64) string {
	return strconv.FormatInt(i, 10)
}

// UInt642Str uint64 to string
func UInt642Str(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// Int2Str int to string
func Int2Str(i int) string {
	return strconv.Itoa(i)
}

// Str2Bytes converts string to byte slice without a memory allocation.
func Str2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// Bytes2Str converts byte slice to string without a memory allocation.
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func IsEmptyLine(str string) bool {
	re := regexp.MustCompile(`^\s*$`)

	return re.MatchString(str)
}

func TrimWS(str string) string {
	return strings.Trim(str, "\n\t")
}

func TrimSpaceWS(str string) string {
	return strings.TrimRight(str, " \n\t")
}

func RemoveSliceEmpty(list []string) (fList []string) {
	for i := range list {
		if strings.TrimSpace(list[i]) != "" {
			fList = append(fList, list[i])
		}
	}
	return
}

func SplitRemoveEmpty(s, sep string) []string {
	data := strings.Split(s, sep)
	return RemoveSliceEmpty(data)
}

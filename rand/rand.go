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
package rand

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var digits = []rune("0123456789")

const size = 62

func RandLetters(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(size)]
	}

	return string(b)
}

func RandDigits(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = digits[rand.Intn(10)]
	}

	return string(b)
}

// Rand 随机数
func Rand() int {
	rand.Seed(time.Now().Unix())
	return rand.Int()
}

// NumRand 随机数
func NumRand(num int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(num)
}

// Rand 随机数
func Rand64() int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63()
}

// NumRand 随机数
func NumRand64(num int64) int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63n(num)
}

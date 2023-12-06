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

package exhash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
)

// GenSha256 生成sha256
func GenSha256(code string) string {
	s := sha256.New()
	s.Write([]byte(code))
	return hex.EncodeToString(s.Sum(nil))
}

// GenSha512 生成sha512
func GenSha512(code string) string {
	s := sha512.New()
	s.Write([]byte(code))
	return hex.EncodeToString(s.Sum(nil))
}

// GenSha1 生成sha1
func GenSha1(code string) string {
	s := sha1.New()
	s.Write([]byte(code))
	return hex.EncodeToString(s.Sum(nil))
}

// StringToNumber hashes a given string to a number
func StringToNumber(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

// String hashes a given string
func String(s string) string {
	hash := sha256.New()
	_, _ = io.WriteString(hash, s)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Hex(s string, length int) string {
	h := md5.Sum([]byte(s))
	d := hex.EncodeToString(h[:])
	return d[:length]
}

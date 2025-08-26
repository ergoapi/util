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
	"encoding/hex"
)

// MD5 计算MD5哈希值
// Deprecated: MD5已不安全，仅用于向后兼容。新代码应使用SHA256或更安全的哈希算法
func MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

// Hex 生成指定长度的MD5十六进制字符串
// Deprecated: MD5已不安全，仅用于向后兼容。新代码应使用SHA256或更安全的哈希算法
func Hex(s string, length int) string {
	h := md5.Sum([]byte(s))
	d := hex.EncodeToString(h[:])
	return d[:length]
}

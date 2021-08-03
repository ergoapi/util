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
	b32 "encoding/base32"
	b64 "encoding/base64"
)

// B64EnCode base64加密
func B64EnCode(code string) string {
	return b64.StdEncoding.EncodeToString([]byte(code))
}

// B64Decode base64解密
func B64Decode(code string) (string, error) {
	ds, err := b64.StdEncoding.DecodeString(code)
	return string(ds), err
}

// B32EnCode base32加密
func B32EnCode(code string) string {
	return b32.StdEncoding.EncodeToString([]byte(code))
}

// B32Decode base32解密
func B32Decode(code string) (string, error) {
	ds, err := b32.StdEncoding.DecodeString(code)
	return string(ds), err
}

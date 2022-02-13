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
	"bytes"
	b32 "encoding/base32"
	b64 "encoding/base64"
	"math/big"

	"github.com/ergoapi/util/common"
	"github.com/ergoapi/util/exstr"
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

// B58EnCode base58加密
func B58EnCode(code string) string {
	var src []byte
	intBytes := big.NewInt(0).SetBytes(exstr.Str2Bytes(code))
	int0, int58 := big.NewInt(0), big.NewInt(58)
	for intBytes.Cmp(big.NewInt(0)) > 0 {
		intBytes.DivMod(intBytes, int58, int0)
		src = append(src, exstr.Str2Bytes(common.Base58table)[int0.Int64()])
	}
	return string(reverseBytes(src))
}

// B58Decode base58解密
func B58Decode(code string) (string, error) {
	int0 := big.NewInt(0)
	for _, val := range exstr.Str2Bytes(code) {
		index := bytes.IndexByte([]byte(common.Base58table), val)
		int0.Mul(int0, big.NewInt(58))
		int0.Add(int0, big.NewInt(int64(index)))
	}
	return int0.String(), nil
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

func reverseBytes(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-1-i] = b[len(b)-1-i], b[i]
	}
	return b
}

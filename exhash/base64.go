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

	"github.com/cockroachdb/errors"
)

// B64Encode base64编码
func B64Encode(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

// B64EnCode base64加密 (deprecated: use B64Encode instead)
// Deprecated: Use B64Encode instead
func B64EnCode(code string) string {
	return B64Encode(code)
}

// B64Decode base64解码
func B64Decode(data string) (string, error) {
	if data == "" {
		return "", errors.New("empty input")
	}
	ds, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.Wrap(err, "base64 decode failed")
	}
	return string(ds), nil
}

// B64EncodeBytes base64编码字节数组
func B64EncodeBytes(data []byte) string {
	return b64.StdEncoding.EncodeToString(data)
}

// B64DecodeBytes base64解码到字节数组
func B64DecodeBytes(data string) ([]byte, error) {
	if data == "" {
		return nil, errors.New("empty input")
	}
	return b64.StdEncoding.DecodeString(data)
}

// B58Encode base58编码
func B58Encode(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// 将字节数组转换为大整数
	intBytes := big.NewInt(0).SetBytes(data)

	// 计算前导零的数量
	leadingZeros := 0
	for _, b := range data {
		if b == 0 {
			leadingZeros++
		} else {
			break
		}
	}

	// 转换为base58
	var result []byte
	int0, int58 := big.NewInt(0), big.NewInt(58)

	for intBytes.Cmp(big.NewInt(0)) > 0 {
		intBytes.DivMod(intBytes, int58, int0)
		result = append(result, common.Base58table[int0.Int64()])
	}

	// 添加前导1（代表前导零）
	for i := 0; i < leadingZeros; i++ {
		result = append(result, '1')
	}

	// 反转结果
	return string(reverseBytes(result))
}

// B58EnCode base58加密 (deprecated: use B58Encode instead)
// Deprecated: Use B58Encode instead
func B58EnCode(code string) string {
	return B58Encode([]byte(code))
}

// B58Decode base58解码
func B58Decode(data string) ([]byte, error) {
	if data == "" {
		return nil, errors.New("empty input")
	}

	// 计算前导1的数量（代表前导零）
	leadingOnes := 0
	for _, c := range data {
		if c == '1' {
			leadingOnes++
		} else {
			break
		}
	}

	// 转换为大整数
	int0 := big.NewInt(0)
	for _, c := range []byte(data) {
		index := bytes.IndexByte([]byte(common.Base58table), c)
		if index == -1 {
			return nil, errors.Errorf("invalid base58 character: %c", c)
		}
		int0.Mul(int0, big.NewInt(58))
		int0.Add(int0, big.NewInt(int64(index)))
	}

	// 转换为字节数组
	result := int0.Bytes()

	// 添加前导零
	for i := 0; i < leadingOnes; i++ {
		result = append([]byte{0}, result...)
	}

	return result, nil
}

// B58DecodeString base58解码到字符串
func B58DecodeString(data string) (string, error) {
	bytes, err := B58Decode(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// B32Encode base32编码
func B32Encode(data string) string {
	return b32.StdEncoding.EncodeToString([]byte(data))
}

// B32EnCode base32加密 (deprecated: use B32Encode instead)
// Deprecated: Use B32Encode instead
func B32EnCode(code string) string {
	return B32Encode(code)
}

// B32Decode base32解码
func B32Decode(data string) (string, error) {
	if data == "" {
		return "", errors.New("empty input")
	}
	ds, err := b32.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.Wrap(err, "base32 decode failed")
	}
	return string(ds), nil
}

func reverseBytes(b []byte) []byte {
	result := make([]byte, len(b))
	copy(result, b)
	for i := 0; i < len(result)/2; i++ {
		result[i], result[len(result)-1-i] = result[len(result)-1-i], result[i]
	}
	return result
}

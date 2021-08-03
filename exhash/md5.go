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
	"fmt"
	"time"
)

// MD5 md5
func MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

// CryptoPass crypto password use salt
func CryptoPass(salt, raw string) string {
	//salt, err := models.ConfigsGet("salt")
	//if err != nil {
	//	return "", fmt.Errorf("query salt from mysql fail: %v", err)
	//}
	return MD5(salt + "<-*Uk30^96eY*->" + raw)
}

func GenUUIDForUser(username string) string {
	return MD5(username + fmt.Sprint(time.Now().UnixNano()))
}

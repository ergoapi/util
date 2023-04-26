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

package exid

import (
	"fmt"
	"strings"
	"time"

	"github.com/ergoapi/util/exhash"
	"github.com/google/uuid"
)

// GenUUID 生成新的uuid
func GenUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}

// CheckUUID 检查uuid是否合法
func CheckUUID(uid string) bool {
	_, err := uuid.Parse(uid)
	return err != nil
}

// Deprecated: use HashUID instead
func GenUID(username string) string {
	return exhash.MD5(username + fmt.Sprint(time.Now().UnixNano()))
}

// HashUID 生成新的uid
func HashUID(username string, prefix string) string {
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "-") {
		prefix = fmt.Sprintf("%s-", prefix)
	}
	return fmt.Sprintf("%s%s", prefix, exhash.MD5(username+fmt.Sprint(time.Now().UnixNano())))
}

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

package zos

import "github.com/google/uuid"

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
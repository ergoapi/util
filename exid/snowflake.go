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
	"github.com/bwmarrin/snowflake"
)

var (
	snowflakeNode *snowflake.Node
)

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	snowflakeNode = node
}

func GenSnowflakeID() int64 {
	return snowflakeNode.Generate().Int64()
}

func GenSnowflakeIDStr() string {
	return snowflakeNode.Generate().String()
}

// ParseID 解析雪花ID
func ParseID(id int64) map[string]any {
	snowflakeID := snowflake.ParseInt64(id)
	return map[string]any{
		"time":     snowflakeID.Time(),
		"node":     snowflakeID.Node(),
		"workerID": snowflakeID.Node() % 31,
		"step":     snowflakeID.Step(),
	}
}

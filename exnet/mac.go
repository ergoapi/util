//  Copyright (c) 2022. The EFF Team Authors.
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

package exnet

import (
	"fmt"
	"strings"
)

const macChars = "0123456789abcdef"
const macCapChars = "ABCDEF"

func FormatMacAddr(macAddr string) string {
	buf := make([]byte, 12)
	bufIdx := 0
	for i := 0; i < len(macAddr) && bufIdx < len(buf); i += 1 {
		c := macAddr[i]
		if strings.IndexByte(macChars, c) >= 0 {
			buf[bufIdx] = c
			bufIdx += 1
		} else if strings.IndexByte(macCapChars, c) >= 0 {
			buf[bufIdx] = c - 'A' + 'a'
			bufIdx += 1
		}
	}
	if len(buf) == bufIdx {
		return fmt.Sprintf("%c%c:%c%c:%c%c:%c%c:%c%c:%c%c", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5],
			buf[6], buf[7], buf[8], buf[9], buf[10], buf[11])
	} else {
		return ""
	}
}

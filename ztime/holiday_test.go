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

package ztime

import (
	"testing"
)

func TestHolidayGet(t *testing.T) {
	days := []string{"2021-04-28", "2021-04-29", "2021-04-30", "2021-05-01", "2021-05-04", "2021-05-05", "2021-05-06", "2021-05-07", "2021-05-08", "2021-05-09", "2021-05-10", "2021-09-30", "2021-10-06"}
	for _, day := range days {
		h := HolidayGet(day)
		hnext := HolidayGet(CustomNextDay(day, 1, "2006-01-02"))
		// log.Println(h.Day, h.Name, h.NeedWork)
		// log.Println(hnext.Day, hnext.Name, hnext.NeedWork)
		if !h.NeedWork {
			// 节假日&双休日
			t.Logf("[%v] %v 全天禁止上线, 紧急上线请联系平台保障部SRE团队", day, h.Name)
			continue
		}
		if !hnext.NeedWork && hnext.Name != "双休日" {
			// 节假日提前1天封板
			t.Logf("[%v]  %v 封板全天禁止上线, 紧急上线请联系平台保障部SRE团队", day, hnext.Name)
			continue
		}
		t.Logf("[%v]  %v 业务高峰 紧急上线请联系平台保障部SRE团队", day, h.Name)
	}
}

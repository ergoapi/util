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

import "github.com/6tail/lunar-go/HolidayUtil"

type Holiday struct {
	Day       string `json:"day"`
	IsTiaoxiu bool   `json:"is_tiaoxiu"`
	Name      string `json:"name"`
	NeedWork  bool   `json:"need_work"`
}

func HolidayGet(day string) Holiday {
	var h Holiday
	d := HolidayUtil.GetHoliday(day)
	h.Day = day
	if d == nil {
		t, _ := TimeParse("2006-01-02", day)
		week := int(t.Weekday())
		if week == 0 || week == 7 || week == 6 {
			h.NeedWork = false
			h.Name = "双休日"
			h.IsTiaoxiu = false
			return h
		}
		h.NeedWork = true
		h.Name = "工作日"
		h.IsTiaoxiu = false
		return h
	}
	h.IsTiaoxiu = d.IsWork()
	h.Name = d.GetName()
	if h.IsTiaoxiu {
		h.Name = d.GetName() + "调休"
	}
	h.NeedWork = d.IsWork()
	return h
}

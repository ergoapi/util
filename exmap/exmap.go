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

package exmap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Struct2Map ...
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// Struct2Json2Map ...
func Struct2Json2Map(obj interface{}) (result map[string]interface{}, err error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonBytes, &result)
	return
}

// Map2String ...
func Map2String(data []string) (result string) {
	if len(data) <= 0 {
		return
	}
	for _, v := range data {
		if strings.Contains(v, "\"") {
			result += v
		} else {
			result += "\"" + v + "\""
		}
		result += ","
	}
	result = strings.Trim(result, ",")
	return
}

func MapString2String(labels map[string]string) string {
	result := make([]string, 0)
	for k, v := range labels {
		result = append(result, fmt.Sprintf("%s=%s", k, v))

	}

	return strings.Join(result, ",")
}

// merge label
// the new map will overwrite the old one.
// e.g. new: {"foo": "newbar"} old: {"foo": "bar"} will return {"foo": "newbar"}
func MergeLabels(old map[string]string, new map[string]string) map[string]string {
	if new == nil {
		return old
	}

	if old == nil {
		old = make(map[string]string)
	}

	for key, value := range new {
		old[key] = value
	}
	return old
}
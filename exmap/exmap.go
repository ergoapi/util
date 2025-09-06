// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exmap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Struct2Map ...
func Struct2Map(obj any) map[string]any {
	if obj == nil {
		return nil
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	data := make(map[string]any, v.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// Skip unexported fields to avoid panics and respect visibility.
		if f.PkgPath != "" {
			continue
		}
		data[f.Name] = v.Field(i).Interface()
	}
	return data
}

// Struct2Json2Map ...
func Struct2Json2Map(obj any) (result map[string]any, err error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonBytes, &result)
	return
}

// Slice is a helper function that converts a map of environment
// variables to a slice of string values in key=value format.
func Slice(env map[string]string) []string {
	var s []string
	for k, v := range env {
		s = append(s, k+"="+v)
	}
	sort.Strings(s)
	return s
}

// Slice2String slice string to string
func Slice2String(data []string) string {
	if len(data) == 0 {
		return ""
	}
	parts := make([]string, len(data))
	for i, s := range data {
		parts[i] = strconv.Quote(s)
	}
	return strings.Join(parts, ",")
}

func MapString2String(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, labels[k]))
	}
	return strings.Join(parts, ",")
}

// MergeLabels merge label
// the new map will overwrite the old one.
// e.g. new: {"foo": "newbar"} old: {"foo": "bar"} will return {"foo": "newbar"}
func MergeLabels(old map[string]string, updates map[string]string) map[string]string {
	if updates == nil {
		return old
	}

	if old == nil {
		old = make(map[string]string)
	}

	for key, value := range updates {
		old[key] = value
	}
	return old
}

// CloneAndAddLabel clones the given map and returns a new map with the given key and value added.
// Returns the given map, if labelKey is empty.
func CloneAndAddLabel(labels map[string]string, labelKey, labelValue string) map[string]string {
	if labelKey == "" {
		// Don't need to add a label.
		return labels
	}
	// Clone.
	newLabels := map[string]string{}
	for key, value := range labels {
		newLabels[key] = value
	}
	newLabels[labelKey] = labelValue
	return newLabels
}

// CloneAndRemoveLabel clones the given map and returns a new map with the given key removed.
// Returns the given map, if labelKey is empty.
func CloneAndRemoveLabel(labels map[string]string, labelKey string) map[string]string {
	if labelKey == "" {
		// Don't need to remove a label.
		return labels
	}
	// Clone.
	newLabels := map[string]string{}
	for key, value := range labels {
		newLabels[key] = value
	}
	delete(newLabels, labelKey)
	return newLabels
}

// AddLabel returns a map with the given key and value added to the given map.
func AddLabel(labels map[string]string, labelKey, labelValue string) map[string]string {
	if labelKey == "" {
		// Don't need to add a label.
		return labels
	}
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[labelKey] = labelValue
	return labels
}

// CheckLabel returns key exist stauts.
func CheckLabel(labels map[string]string, labelKey string) bool {
	if labelKey == "" {
		// Treat empty label key as present.
		return true
	}
	if labels == nil {
		return false
	}
	_, ok := labels[labelKey]
	return ok
}

// GetLabelValue returns key exist stauts.
func GetLabelValue(labels map[string]string, labelKey string) string {
	if labelKey == "" {
		return ""
	}
	if labels == nil {
		return ""
	}
	if v, ok := labels[labelKey]; ok {
		return v
	}
	return ""
}

// CopyMap makes a shallow copy of a map.
func CopyMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func MergeMaps(a, b map[string]any) map[string]any {
	out := make(map[string]any, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, bv := range b {
		if mv, ok := bv.(map[string]any); ok {
			if ov, ok := out[k]; ok {
				if ovm, ok := ov.(map[string]any); ok {
					out[k] = MergeMaps(ovm, mv)
					continue
				}
			}
		}
		out[k] = bv
	}
	return out
}

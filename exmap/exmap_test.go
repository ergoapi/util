// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exmap

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestCloneAndAddLabel(t *testing.T) {
	labels := map[string]string{
		"foo1": "bar1",
		"foo2": "bar2",
		"foo3": "bar3",
	}

	cases := []struct {
		labels     map[string]string
		labelKey   string
		labelValue string
		want       map[string]string
	}{
		{
			labels: labels,
			want:   labels,
		},
		{
			labels:     labels,
			labelKey:   "foo4",
			labelValue: "42",
			want: map[string]string{
				"foo1": "bar1",
				"foo2": "bar2",
				"foo3": "bar3",
				"foo4": "42",
			},
		},
	}

	for _, tc := range cases {
		got := CloneAndAddLabel(tc.labels, tc.labelKey, tc.labelValue)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("[Add] got %v, want %v", got, tc.want)
		}
		// now test the inverse.
		gotrm := CloneAndRemoveLabel(got, tc.labelKey)
		if !reflect.DeepEqual(gotrm, tc.labels) {
			t.Errorf("[RM] got %v, want %v", gotrm, tc.labels)
		}
	}
}

func TestAddLabel(t *testing.T) {
	labels := map[string]string{
		"foo1": "bar1",
		"foo2": "bar2",
		"foo3": "bar3",
	}

	cases := []struct {
		labels     map[string]string
		labelKey   string
		labelValue string
		want       map[string]string
	}{
		{
			labels: labels,
			want:   labels,
		},
		{
			labels:     labels,
			labelKey:   "foo4",
			labelValue: "food",
			want: map[string]string{
				"foo1": "bar1",
				"foo2": "bar2",
				"foo3": "bar3",
				"foo4": "food",
			},
		},
		{
			labels:     nil,
			labelKey:   "foo4",
			labelValue: "food",
			want: map[string]string{
				"foo4": "food",
			},
		},
	}

	for _, tc := range cases {
		got := AddLabel(tc.labels, tc.labelKey, tc.labelValue)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("got %v, want %v", got, tc.want)
		}
	}
}

func TestCheckLabel(t *testing.T) {
	labels := map[string]string{
		"foo1": "bar1",
		"foo2": "bar2",
		"foo3": "bar3",
	}

	cases := []struct {
		labels   map[string]string
		labelKey string
		want     bool
	}{
		{
			labels: labels,
			want:   true,
		},
		{
			labels:   labels,
			labelKey: "foo3",
			want:     true,
		},
		{
			labels:   labels,
			labelKey: "foo4",
			want:     false,
		},
		{
			labels:   nil,
			labelKey: "foo4",
			want:     false,
		},
	}

	for _, tc := range cases {
		got := CheckLabel(tc.labels, tc.labelKey)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("got %v, want %v", got, tc.want)
		}
	}
}

func TestMergeValues(t *testing.T) {
	nestedMap := map[string]any{
		"foo": "bar",
		"baz": map[string]string{
			"cool": "stuff",
		},
	}
	anotherNestedMap := map[string]any{
		"foo": "bar",
		"baz": map[string]string{
			"cool":    "things",
			"awesome": "stuff",
		},
	}
	flatMap := map[string]any{
		"foo": "bar",
		"baz": "stuff",
	}
	anotherFlatMap := map[string]any{
		"testing": "fun",
	}

	testMap := MergeMaps(flatMap, nestedMap)
	equal := reflect.DeepEqual(testMap, nestedMap)
	if !equal {
		t.Errorf("Expected a nested map to overwrite a flat value. Expected: %v, got %v", nestedMap, testMap)
	}

	testMap = MergeMaps(nestedMap, flatMap)
	equal = reflect.DeepEqual(testMap, flatMap)
	if !equal {
		t.Errorf("Expected a flat value to overwrite a map. Expected: %v, got %v", flatMap, testMap)
	}

	testMap = MergeMaps(nestedMap, anotherNestedMap)
	equal = reflect.DeepEqual(testMap, anotherNestedMap)
	if !equal {
		t.Errorf("Expected a nested map to overwrite another nested map. Expected: %v, got %v", anotherNestedMap, testMap)
	}

	testMap = MergeMaps(anotherFlatMap, anotherNestedMap)
	expectedMap := map[string]any{
		"testing": "fun",
		"foo":     "bar",
		"baz": map[string]string{
			"cool":    "things",
			"awesome": "stuff",
		},
	}
	equal = reflect.DeepEqual(testMap, expectedMap)
	if !equal {
		t.Errorf("Expected a map with different keys to merge properly with another map. Expected: %v, got %v",
			expectedMap, testMap)
	}
}

func TestStruct2Map(t *testing.T) {
	type inner struct{ X int }
	type S struct {
		A     string
		B     int
		c     string // unexported, should be skipped
		Inner inner
	}
	s := S{A: "hello", B: 2, c: "secret", Inner: inner{X: 9}}

	got := Struct2Map(s)
	if got == nil {
		t.Fatalf("expected non-nil map")
	}
	if v, ok := got["A"]; !ok || v != "hello" {
		t.Errorf("expected A=hello, got %v (ok=%v)", v, ok)
	}
	if _, ok := got["c"]; ok {
		t.Errorf("did not expect unexported field 'c' in map")
	}

	// Pointer input should behave the same.
	gotPtr := Struct2Map(&s)
	if !reflect.DeepEqual(got, gotPtr) {
		t.Errorf("pointer and value results differ: %v vs %v", got, gotPtr)
	}

	// Nil and non-struct should return nil
	if res := Struct2Map(nil); res != nil {
		t.Errorf("expected nil for nil input, got %v", res)
	}
	if res := Struct2Map(123); res != nil {
		t.Errorf("expected nil for non-struct input, got %v", res)
	}
}

func TestSlice2String(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{[]string{"a", "b", "c"}, strings.Join([]string{strconv.Quote("a"), strconv.Quote("b"), strconv.Quote("c")}, ",")},
		{[]string{"a\"b", "x,y", ""}, strings.Join([]string{strconv.Quote("a\"b"), strconv.Quote("x,y"), strconv.Quote("")}, ",")},
		{nil, ""},
		{[]string{}, ""},
	}
	for _, tc := range cases {
		got := Slice2String(tc.in)
		if got != tc.want {
			t.Errorf("Slice2String(%v) got %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestMapString2StringOrder(t *testing.T) {
	m := map[string]string{"b": "2", "a": "1", "c": "3"}
	got := MapString2String(m)
	want := "a=1,b=2,c=3"
	if got != want {
		t.Errorf("MapString2String order deterministic: got %q, want %q", got, want)
	}
}

func TestSyncMapLen(t *testing.T) {
	var nilMap *sync.Map
	if got := SyncMapLen(nilMap); got != 0 {
		t.Errorf("nil map length = %d, want 0", got)
	}
	m := &sync.Map{}
	m.Store("a", 1)
	m.Store("b", 2)
	if got := SyncMapLen(m); got != 2 {
		t.Errorf("length = %d, want 2", got)
	}
}

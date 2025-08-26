// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

// Package exerror provides utilities.
package exerror

import "fmt"

type ErgoError struct {
	Message string
}

func (ee *ErgoError) Error() string {
	return ee.Message
}

func (ee *ErgoError) String() string {
	return ee.Message
}

func Bomb(format string, args ...any) {
	panic(ErgoError{Message: fmt.Sprintf(format, args...)})
}

func Dangerous(v any) {
	if v == nil {
		return
	}

	switch t := v.(type) {
	case string:
		if t != "" {
			panic(ErgoError{Message: t})
		}
	case error:
		panic(ErgoError{Message: t.Error()})
	}
}

func Boka(value string, v any) {
	if v == nil {
		return
	}
	Bomb(value)
}

// CheckAndExit check & exit
func CheckAndExit(err error) {
	if err != nil {
		panic(err)
	}
}

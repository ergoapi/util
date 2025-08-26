// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exstr

import "unsafe"

// UnsafeString returns a string pointer without allocation
//
//nolint:gosec // unsafe is used for better performance here
func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

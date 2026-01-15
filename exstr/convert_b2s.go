// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exstr

import "unsafe"

// UnsafeString returns a string pointer without allocation.
//
// WARNING: This function uses unsafe operations for performance optimization.
// The returned string shares memory with the input byte slice.
//
// IMPORTANT CONSTRAINTS:
//   - The input byte slice MUST NOT be modified after calling this function,
//     as it would violate Go's string immutability guarantee.
//   - The returned string is only valid as long as the byte slice is not modified or garbage collected.
//   - Use this only when you have strict control over the byte slice lifecycle.
//   - For general use cases, prefer the standard string(b) conversion.
//
//nolint:gosec // unsafe is used for better performance here
func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exmap

import "sync"

func SyncMapLen(m *sync.Map) int {
	if m == nil {
		return 0
	}
	length := 0
	m.Range(func(key, value any) bool {
		length++
		return true
	})
	return length
}

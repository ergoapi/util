// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package zos

import (
	"runtime"
	"runtime/debug"
)

// ForceGC 强制手动GC垃圾回收(阻塞).
func ForceGC() {
	runtime.GC()
	debug.FreeOSMemory()
}

// TriggerGC 触发GC(非阻塞).
func TriggerGC() {
	go func() {
		ForceGC()
	}()
}

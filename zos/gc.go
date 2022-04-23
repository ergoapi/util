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

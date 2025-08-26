// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package formatter

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// findLibraryCallerInternal 内部共享的查找逻辑
func findLibraryCallerInternal(libraryPathPrefix string) *runtime.Frame {
	// 从调用堆栈的第4帧开始查找（跳过本函数、findLibraryCaller、Format函数）
	pcs := make([]uintptr, 50)
	n := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		// 检查是否是库内的调用
		if strings.Contains(frame.File, libraryPathPrefix) {
			return &runtime.Frame{
				PC:       frame.PC,
				Func:     frame.Func,
				Function: frame.Function,
				File:     frame.File,
				Line:     frame.Line,
			}
		}
		if !more {
			break
		}
	}
	return nil
}

// FilteredTextFormatter 是一个自定义的logrus.TextFormatter，它可以过滤调用堆栈
// 只显示特定库路径前缀的调用信息
type FilteredTextFormatter struct {
	logrus.TextFormatter
	// 要保留的库路径前缀，例如 "github.com/ergoapi/util"
	LibraryPathPrefix string
}

// Format 重写了logrus.TextFormatter的Format方法
// 在调用原始Format方法之前，它会检查并过滤调用者信息
func (f *FilteredTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 如果启用了调用者报告且有调用者信息
	if entry.Caller != nil {
		// 检查调用者文件路径是否包含指定的库路径前缀
		if !strings.Contains(entry.Caller.File, f.LibraryPathPrefix) {
			// 如果不包含，尝试查找堆栈中第一个匹配的调用者
			if caller := f.findLibraryCaller(); caller != nil {
				// 替换为库内的调用者信息
				entry.Caller = caller
			}
		}
	}

	// 调用原始的Format方法
	return f.TextFormatter.Format(entry)
}

// findLibraryCaller 在调用堆栈中查找第一个匹配库路径前缀的调用者
func (f *FilteredTextFormatter) findLibraryCaller() *runtime.Frame {
	return findLibraryCallerInternal(f.LibraryPathPrefix)
}

// FilteredJSONFormatter 是一个自定义的logrus.JSONFormatter，它可以过滤调用堆栈
// 只显示特定库路径前缀的调用信息
type FilteredJSONFormatter struct {
	logrus.JSONFormatter
	// 要保留的库路径前缀，例如 "github.com/ergoapi/util"
	LibraryPathPrefix string
}

// Format 重写了logrus.JSONFormatter的Format方法
// 在调用原始Format方法之前，它会检查并过滤调用者信息
func (f *FilteredJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 如果启用了调用者报告且有调用者信息
	if entry.Caller != nil {
		// 检查调用者文件路径是否包含指定的库路径前缀
		if !strings.Contains(entry.Caller.File, f.LibraryPathPrefix) {
			// 如果不包含，尝试查找堆栈中第一个匹配的调用者
			if caller := f.findLibraryCaller(); caller != nil {
				// 替换为库内的调用者信息
				entry.Caller = caller
			}
		}
	}

	// 调用原始的Format方法
	return f.JSONFormatter.Format(entry)
}

// findLibraryCaller 在调用堆栈中查找第一个匹配库路径前缀的调用者
func (f *FilteredJSONFormatter) findLibraryCaller() *runtime.Frame {
	return findLibraryCallerInternal(f.LibraryPathPrefix)
}

// NewFilteredTextFormatter 创建一个新的FilteredTextFormatter实例
func NewFilteredTextFormatter(libraryPathPrefix string) *FilteredTextFormatter {
	return &FilteredTextFormatter{
		TextFormatter:     logrus.TextFormatter{},
		LibraryPathPrefix: libraryPathPrefix,
	}
}

// NewFilteredJSONFormatter 创建一个新的FilteredJSONFormatter实例
func NewFilteredJSONFormatter(libraryPathPrefix string) *FilteredJSONFormatter {
	return &FilteredJSONFormatter{
		JSONFormatter:     logrus.JSONFormatter{},
		LibraryPathPrefix: libraryPathPrefix,
	}
}

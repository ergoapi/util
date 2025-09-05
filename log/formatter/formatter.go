// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package formatter

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// FilteredTextFormatter 是一个自定义的logrus.TextFormatter，它可以过滤调用堆栈
// 只显示特定库路径前缀的调用信息
type FilteredTextFormatter struct {
	logrus.TextFormatter
	// 要保留的库路径前缀，例如 "github.com/ergoapi/util"
	LibraryPathPrefix string
}

// Format 重写了logrus.TextFormatter的Format方法
// 如果调用者信息来自库内（包含LibraryPathPrefix），则将其设置为nil以隐藏
func (f *FilteredTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 如果启用了调用者报告且有调用者信息
	if entry.Caller != nil {
		// 检查调用者文件路径是否包含指定的库路径前缀
		if strings.Contains(entry.Caller.File, f.LibraryPathPrefix) {
			// 如果包含库路径前缀，将调用者信息设置为nil以隐藏
			entry.Caller = nil
		}
	}

	// 调用原始的Format方法
	return f.TextFormatter.Format(entry)
}

// FilteredJSONFormatter 是一个自定义的logrus.JSONFormatter，它可以过滤调用堆栈
// 只显示特定库路径前缀的调用信息
type FilteredJSONFormatter struct {
	logrus.JSONFormatter
	// 要保留的库路径前缀，例如 "github.com/ergoapi/util"
	LibraryPathPrefix string
}

// Format 重写了logrus.JSONFormatter的Format方法
// 如果调用者信息来自库内（包含LibraryPathPrefix），则将其设置为nil以隐藏
func (f *FilteredJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 如果启用了调用者报告且有调用者信息
	if entry.Caller != nil {
		// 检查调用者文件路径是否包含指定的库路径前缀
		if strings.Contains(entry.Caller.File, f.LibraryPathPrefix) {
			// 如果包含库路径前缀，将调用者信息设置为nil以隐藏
			entry.Caller = nil
		}
	}

	// 调用原始的Format方法
	return f.JSONFormatter.Format(entry)
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

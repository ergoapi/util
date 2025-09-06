// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package formatter

import (
	"maps"
	"strings"

	"github.com/sirupsen/logrus"
)

// defaultLibraryPath 默认的库路径
const defaultLibraryPath = "github.com/ergoapi/util"

// FilteredTextFormatter 是一个自定义的logrus.TextFormatter，它可以过滤调用堆栈
// 只显示特定库路径前缀的调用信息
type FilteredTextFormatter struct {
	logrus.TextFormatter
	// 要过滤的库路径前缀列表，例如 ["github.com/ergoapi/util", "github.com/mycompany"]
	LibraryPathPrefixes []string
}

// Format 重写了logrus.TextFormatter的Format方法
// 如果调用者信息来自库内（包含LibraryPathPrefixes中的任何路径），则将其设置为nil以隐藏
func (f *FilteredTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 如果启用了调用者报告且有调用者信息
	if entry.Caller != nil {
		// 检查调用者文件路径是否包含任何指定的库路径前缀
		for _, prefix := range f.LibraryPathPrefixes {
			if strings.Contains(entry.Caller.File, prefix) {
				// 克隆 entry，避免修改原始 entry 影响其他 Hook/Formatter
				cloned := *entry
				// 深拷贝 Data map，防止并发修改
				cloned.Data = maps.Clone(entry.Data)
				cloned.Caller = nil
				return f.TextFormatter.Format(&cloned)
			}
		}
	}

	// 调用原始的Format方法
	return f.TextFormatter.Format(entry)
}

// FilteredJSONFormatter 是一个自定义的logrus.JSONFormatter，它可以过滤调用堆栈
// 只显示特定库路径前缀的调用信息
type FilteredJSONFormatter struct {
	logrus.JSONFormatter
	// 要过滤的库路径前缀列表，例如 ["github.com/ergoapi/util", "github.com/mycompany"]
	LibraryPathPrefixes []string
}

// Format 重写了logrus.JSONFormatter的Format方法
// 如果调用者信息来自库内（包含LibraryPathPrefixes中的任何路径），则将其设置为nil以隐藏
func (f *FilteredJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 如果启用了调用者报告且有调用者信息
	if entry.Caller != nil {
		// 检查调用者文件路径是否包含任何指定的库路径前缀
		for _, prefix := range f.LibraryPathPrefixes {
			if strings.Contains(entry.Caller.File, prefix) {
				// 克隆 entry，避免副作用
				cloned := *entry
				// 深拷贝 Data map，防止并发修改
				cloned.Data = maps.Clone(entry.Data)
				cloned.Caller = nil
				return f.JSONFormatter.Format(&cloned)
			}
		}
	}

	// 调用原始的Format方法
	return f.JSONFormatter.Format(entry)
}

// NewFilteredTextFormatter 创建一个新的FilteredTextFormatter实例
// 如果没有传入additionalPrefixes，将使用默认的"github.com/ergoapi/util"
// 如果传入了additionalPrefixes，将与默认值合并
func NewFilteredTextFormatter(additionalPrefixes ...string) *FilteredTextFormatter {
	prefixes := []string{defaultLibraryPath}
	prefixes = append(prefixes, additionalPrefixes...)
	return &FilteredTextFormatter{
		TextFormatter: logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
		LibraryPathPrefixes: prefixes,
	}
}

// NewFilteredJSONFormatter 创建一个新的FilteredJSONFormatter实例
// 如果没有传入additionalPrefixes，将使用默认的"github.com/ergoapi/util"
// 如果传入了additionalPrefixes，将与默认值合并
func NewFilteredJSONFormatter(additionalPrefixes ...string) *FilteredJSONFormatter {
	prefixes := []string{defaultLibraryPath}
	prefixes = append(prefixes, additionalPrefixes...)
	return &FilteredJSONFormatter{
		JSONFormatter: logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
		LibraryPathPrefixes: prefixes,
	}
}

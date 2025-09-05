// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package main

import (
	"os"

	"github.com/ergoapi/util/log/formatter"
	"github.com/sirupsen/logrus"
)

func init() {
	// 创建一个 FilteredTextFormatter
	// 当日志调用来自 github.com/ergoapi/util 内部时，会隐藏调用者信息（func 和 file 字段）
	// 当日志调用来自其他包时，正常显示调用者信息
	textFormatter := formatter.NewFilteredTextFormatter("github.com/ergoapi/util")

	// 可选：配置格式化器的其他选项
	textFormatter.TimestampFormat = "2006-01-02 15:04:05"
	textFormatter.FullTimestamp = true

	// 设置 logrus 使用我们的自定义 formatter
	logrus.SetFormatter(textFormatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	// 重要：必须启用调用者报告功能，否则 formatter 不会接收到调用者信息
	logrus.SetReportCaller(true)
}

func main() {
	// 示例1：直接从 main 函数调用
	// 这将显示调用者信息，因为调用来自应用代码而非 util 库内部
	logrus.Info("应用层日志：显示调用者信息")

	// 示例2：带字段的日志
	logrus.WithFields(logrus.Fields{
		"user_id": "12345",
		"action":  "login",
	}).Info("用户登录事件")

	// 示例3：模拟库内部调用（实际上仍是应用代码）
	simulateLibraryCall()

	// 示例4：使用 JSON 格式
	useJSONFormatter()

	// 示例5：展示实际效果对比
	demonstrateFiltering()
}

// simulateLibraryCall 模拟库函数调用
// 注意：这个函数实际上不在 github.com/ergoapi/util 包内，所以仍会显示调用者信息
func simulateLibraryCall() {
	logrus.Info("模拟库调用：仍显示调用者信息（因为不在 util 包内）")
}

// useJSONFormatter 展示 FilteredJSONFormatter 的用法
func useJSONFormatter() {
	// 创建一个 FilteredJSONFormatter
	jsonFormatter := formatter.NewFilteredJSONFormatter("github.com/ergoapi/util")
	jsonFormatter.TimestampFormat = "2006-01-02T15:04:05.000Z"
	jsonFormatter.PrettyPrint = false // 紧凑的 JSON 输出

	// 临时切换到 JSON 格式
	originalFormatter := logrus.StandardLogger().Formatter
	logrus.SetFormatter(jsonFormatter)

	// 输出 JSON 格式的日志
	logrus.WithFields(logrus.Fields{
		"trace_id": "abc-123-def",
		"span_id":  "span-456",
		"service":  "user-service",
	}).Info("JSON格式日志：显示调用者信息")

	// 恢复原来的 formatter
	logrus.SetFormatter(originalFormatter)
}

// demonstrateFiltering 展示过滤效果的对比
func demonstrateFiltering() {
	logrus.Info("\n=== 演示调用者信息过滤 ===")

	// 不使用过滤的标准 TextFormatter
	standardFormatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	}
	logrus.SetFormatter(standardFormatter)
	logrus.Info("标准格式化器：总是显示调用者信息")

	// 使用过滤的 FilteredTextFormatter
	filteredFormatter := formatter.NewFilteredTextFormatter("github.com/ergoapi/util")
	filteredFormatter.FullTimestamp = true
	filteredFormatter.TimestampFormat = "15:04:05"
	logrus.SetFormatter(filteredFormatter)
	logrus.Info("过滤格式化器：根据调用来源决定是否显示")

	// 实际使用场景说明
	logrus.Info(`
使用场景说明：
1. 当你在使用 gorm 的 glog logger 时，所有来自 glog 包的日志将不显示 func 和 file 字段
2. 当你的应用代码直接调用 logrus 时，会正常显示调用者信息
3. 这样可以避免日志中充斥着库内部的调用路径，只关注应用层的调用位置`)
}

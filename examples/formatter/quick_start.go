// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

// Package main 展示如何在实际项目中配置 FilteredFormatter 来隐藏库内部的调用者信息
package main

import (
	"github.com/ergoapi/util/log/formatter"
	"github.com/sirupsen/logrus"
)

// InitLogger 初始化日志配置（通常在 main 函数或 init 函数中调用）
func InitLogger() {
	// 方案1：使用 JSON 格式（适合生产环境，便于日志收集系统解析）
	if false { // 改为 true 使用 JSON 格式
		jsonFormatter := formatter.NewFilteredJSONFormatter("github.com/ergoapi/util")
		jsonFormatter.TimestampFormat = "2006-01-02T15:04:05.000Z"
		jsonFormatter.PrettyPrint = false // 生产环境使用紧凑格式
		logrus.SetFormatter(jsonFormatter)
	} else {
		// 方案2：使用文本格式（适合开发环境，便于人工阅读）
		textFormatter := formatter.NewFilteredTextFormatter("github.com/ergoapi/util")
		textFormatter.FullTimestamp = true
		textFormatter.TimestampFormat = "2006-01-02 15:04:05"
		textFormatter.DisableColors = false // 开发环境可以启用颜色
		logrus.SetFormatter(textFormatter)
	}

	// 必须启用调用者报告，否则 formatter 无法获取调用者信息
	logrus.SetReportCaller(true)

	// 设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	// 初始化日志配置
	InitLogger()

	// 示例：应用代码的日志（会显示调用者信息）
	logrus.Info("应用启动")
	logrus.WithFields(logrus.Fields{
		"version": "1.0.0",
		"env":     "development",
	}).Info("配置加载完成")

	// 当你使用 gorm 配合 glog.GLogger 时：
	// db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	//     Logger: &glog.GLogger{
	//         LogLevel:      logger.Info,
	//         SlowThreshold: 200 * time.Millisecond,
	//     },
	// })
	//
	// 执行数据库操作时的日志输出：
	// - SQL 查询日志来自 glog 包，不会显示 func 和 file 字段
	// - 你的应用代码日志会正常显示调用位置

	// 实际效果对比
	logrus.Info("这条日志会显示调用位置: quick_start.go:XX")

	// 如果这是从 github.com/ergoapi/util 内部调用的
	// logrus.Info("这条日志不会显示 func 和 file 字段")

	logrus.Info(`
使用建议：
1. 在应用启动时调用 InitLogger() 配置日志格式
2. 生产环境使用 JSON 格式，便于 ELK、Fluentd 等日志系统收集
3. 开发环境使用文本格式，便于调试和阅读
4. 可以根据需要调整 LibraryPathPrefix 参数来过滤不同的库`)
}

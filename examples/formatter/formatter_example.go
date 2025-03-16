package main

import (
	"os"

	"github.com/ergoapi/util/log/formatter"

	"github.com/sirupsen/logrus"
)

func init() {
	// 创建一个过滤的TextFormatter，只显示ergoapi/util库内的调用
	textFormatter := formatter.NewFilteredTextFormatter("github.com/ergoapi/util")

	// 设置logrus使用我们的自定义formatter
	logrus.SetFormatter(textFormatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	// 启用调用者报告功能
	logrus.SetReportCaller(true)
}

func main() {
	// 直接调用日志函数
	logrus.Info("这是直接从main函数调用的日志")

	// 通过库函数调用日志
	logFromLibrary()

	// 使用JSON格式输出
	useJSONFormatter()
}

// logFromLibrary 模拟从库函数中调用日志
func logFromLibrary() {
	logrus.Info("这是从库函数中调用的日志")
}

// useJSONFormatter 展示如何使用FilteredJSONFormatter
func useJSONFormatter() {
	// 创建一个过滤的JSONFormatter
	jsonFormatter := formatter.NewFilteredJSONFormatter("github.com/ergoapi/util")

	// 临时切换到JSON格式
	originalFormatter := logrus.StandardLogger().Formatter
	logrus.SetFormatter(jsonFormatter)

	// 输出JSON格式的日志
	logrus.Info("这是JSON格式的日志")

	// 恢复原来的formatter
	logrus.SetFormatter(originalFormatter)
}

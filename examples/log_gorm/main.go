// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package main

import (
	"context"
	"errors"
	"time"

	"github.com/ergoapi/util/exctx"
	"github.com/ergoapi/util/log/formatter"
	"github.com/ergoapi/util/log/glog"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 示例模型
type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:100"`
	Age  int
}

func main() {
	// 配置 logrus 使用 FilteredJSONFormatter
	// 这将隐藏来自 github.com/ergoapi/util 内部的调用者信息
	jsonFormatter := formatter.NewFilteredJSONFormatter("github.com/ergoapi/util")
	jsonFormatter.TimestampFormat = "2006-01-02T15:04:05.000Z"
	logrus.SetFormatter(jsonFormatter)
	logrus.SetReportCaller(true) // 启用调用者报告
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("=== 开始 GORM + GLogger 示例 ===")

	// 示例1：基本的 GLogger 使用
	basicGLoggerExample()

	// 示例2：与 GORM 集成使用
	gormIntegrationExample()

	// 示例3：使用带 trace 上下文的日志
	traceContextExample()

	logrus.Info("=== 示例完成 ===")
}

// basicGLoggerExample 展示基本的 GLogger 使用
func basicGLoggerExample() {
	logrus.Info("\n--- 基本 GLogger 使用 ---")

	ctx := context.Background()
	glogger := &glog.DefaultGLogger

	// 测试不同级别的日志
	// 注意：这些日志来自 glog 包内部，调用者信息会被隐藏
	glogger.Info(ctx, "信息日志", "user_id", 12345)
	glogger.Warn(ctx, "警告日志", "threshold", 0.8)
	glogger.Error(ctx, "错误日志", "error_code", "ERR_001")

	// 测试 Trace 功能（通常用于数据库操作）
	startTime := time.Now()
	mockSQLFunc := func() (string, int64) {
		return "SELECT * FROM users WHERE id = ?", 1
	}

	// 正常查询
	glogger.Trace(ctx, startTime, mockSQLFunc, nil)

	// 模拟慢查询
	slowStart := time.Now().Add(-300 * time.Millisecond) // 模拟已经运行了300ms
	glogger.Trace(ctx, slowStart, mockSQLFunc, nil)

	// 模拟错误查询
	glogger.Trace(ctx, startTime, mockSQLFunc, errors.New("connection refused"))
}

// gormIntegrationExample 展示与 GORM 的集成
func gormIntegrationExample() {
	logrus.Info("\n--- GORM 集成示例 ---")

	// 创建自定义的 GLogger 实例
	customLogger := &glog.GLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 100 * time.Millisecond, // 设置慢查询阈值为100ms
	}

	// 初始化 GORM 数据库连接
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: customLogger, // 使用我们的自定义 logger
	})
	if err != nil {
		logrus.WithError(err).Fatal("无法连接数据库")
	}

	// 自动迁移
	if err := db.AutoMigrate(&User{}); err != nil {
		logrus.WithError(err).Fatal("自动迁移失败")
	}

	// 创建测试数据
	users := []User{
		{Name: "张三", Age: 25},
		{Name: "李四", Age: 30},
		{Name: "王五", Age: 35},
	}

	// 批量插入（这些 SQL 日志将通过 GLogger 输出，且不显示 glog 包的调用者信息）
	ctx := context.Background()
	db.WithContext(ctx).Create(&users)

	// 查询操作
	var user User
	db.WithContext(ctx).First(&user, 1)
	logrus.Infof("查询到用户: %+v", user)

	// 更新操作
	db.WithContext(ctx).Model(&user).Update("age", 26)

	// 模拟慢查询（使用 sleep 来模拟）
	db.WithContext(ctx).Exec("SELECT * FROM users WHERE 1=1") // 这会被记录为普通查询

	// 错误查询示例
	var nonExistentUser User
	err = db.WithContext(ctx).First(&nonExistentUser, 999).Error
	if err != nil {
		logrus.Info("预期的错误（记录未找到）已被 GLogger 记录")
	}
}

// traceContextExample 展示使用 trace 上下文
func traceContextExample() {
	logrus.Info("\n--- Trace 上下文示例 ---")

	// 创建带 trace 信息的上下文
	// 创建带有 trace 信息的 context
	trace := exctx.NewTrace()
	trace.TraceID = "trace-123"
	trace.SpanID = "span-456"
	trace.CSpanID = "child-789"
	ctx := exctx.SetTraceContext(context.Background(), trace)

	// 使用带 trace 的 logger
	glogger := &glog.GLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}

	// 设置为 logrus 的默认 logger 接口
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: glogger,
	})
	if err != nil {
		logrus.WithError(err).Fatal("无法连接数据库")
	}

	// 使用带 trace 的上下文执行查询
	db.AutoMigrate(&User{})

	user := User{Name: "测试用户", Age: 20}
	db.WithContext(ctx).Create(&user)

	// 日志将包含 traceID、spanID 和 childSpanID
	logrus.WithContext(ctx).Info("用户创建完成")

	// 执行查询
	var foundUser User
	db.WithContext(ctx).First(&foundUser, user.ID)

	// 输出说明
	logrus.Info(`
说明：
1. 当使用 FilteredJSONFormatter 或 FilteredTextFormatter 时
2. 所有来自 github.com/ergoapi/util/log/glog 的日志不会显示 func 和 file 字段
3. 但来自应用代码（如本文件）的日志会正常显示调用位置
4. 这样可以专注于应用层的调用链，避免库内部路径的干扰`)
}

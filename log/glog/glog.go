// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package glog

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ergoapi/util/exctx"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var DefaultGLogger = GLogger{
	LogLevel:      logger.Info,
	SlowThreshold: 200 * time.Millisecond,
}

// getFilteredFileWithLineNum 获取过滤后的文件和行号信息
// 跳过包含 "github.com/ergoapi/util" 的调用栈帧，返回应用代码的文件和行号
func getFilteredFileWithLineNum() string {
	// 使用与 gorm utils.FileWithLineNum() 相似的逻辑
	pcs := [13]uintptr{}
	// 从第3帧开始，跳过本函数和调用者
	length := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:length])

	for i := 0; i < length; i++ {
		frame, _ := frames.Next()
		// 跳过包含 "github.com/ergoapi/util" 的调用栈帧
		// 同时跳过 gorm 内部文件、Go 标准库和生成的文件
		if (!strings.Contains(frame.File, "github.com/ergoapi/util") &&
			!strings.Contains(frame.File, "gorm.io/") &&
			!strings.HasSuffix(frame.File, "_test.go")) &&
			!strings.HasSuffix(frame.File, ".gen.go") &&
			frame.File != "" {
			return string(strconv.AppendInt(append([]byte(frame.File), ':'), int64(frame.Line), 10))
		}
	}

	// 如果没有找到合适的调用栈帧，回退到原始方法
	return utils.FileWithLineNum()
}

type GLogger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func (mgl *GLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	mgl.LogLevel = logLevel
	return mgl
}

// logWithLevel 是一个辅助函数，用于处理通用的日志格式化和输出逻辑
func (mgl *GLogger) logWithLevel(ctx context.Context, level logrus.Level, message string, values ...any) {
	trace := exctx.GetTraceContext(ctx)
	msg := fmt.Sprintf("message=%+v||values=%+v", message, fmt.Sprint(values...))
	msg = strings.Trim(fmt.Sprintf("%q", msg), "\"")

	entry := logrus.WithFields(logrus.Fields{
		"traceID":     trace.TraceID,
		"SpanID":      trace.SpanID,
		"childSpanID": trace.CSpanID,
		"Tag":         "gorm",
	})

	switch level {
	case logrus.InfoLevel:
		entry.Infof(msg)
	case logrus.WarnLevel:
		entry.Warn(msg)
	case logrus.ErrorLevel:
		entry.Error(msg)
	}
}

func (mgl *GLogger) Info(ctx context.Context, message string, values ...any) {
	mgl.logWithLevel(ctx, logrus.InfoLevel, message, values...)
}

func (mgl *GLogger) Warn(ctx context.Context, message string, values ...any) {
	mgl.logWithLevel(ctx, logrus.WarnLevel, message, values...)
}

func (mgl *GLogger) Error(ctx context.Context, message string, values ...any) {
	mgl.logWithLevel(ctx, logrus.ErrorLevel, message, values...)
}

// createTraceFields 创建通用的 trace 字段
func createTraceFields(trace *exctx.TraceContext, begin time.Time, elapsed time.Duration, sql string, rows int64) logrus.Fields {
	fields := logrus.Fields{
		"traceID":         trace.TraceID,
		"SpanID":          trace.SpanID,
		"childSpanID":     trace.CSpanID,
		"Tag":             "gorm",
		"FileWithLineNum": getFilteredFileWithLineNum(),
		"current_time":    begin.Format(time.RFC3339),
		"proc_time":       elapsed.Milliseconds(),
		"sql":             sql,
	}

	if rows == -1 {
		fields["rows"] = "-"
	} else {
		fields["rows"] = rows
	}

	return fields
}

func (mgl *GLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 静默模式直接返回
	if mgl.LogLevel <= 0 {
		return
	}

	trace := exctx.GetTraceContext(ctx)
	elapsed := time.Since(begin)

	// 根据不同情况记录日志
	switch {
	case err != nil && mgl.LogLevel >= logger.Error:
		sql, rows := fc()
		fields := createTraceFields(trace, begin, elapsed, sql, rows)

		// 记录未找到的错误用 Warn 级别
		if rows == -1 || rows == 0 || err == gorm.ErrRecordNotFound {
			logrus.WithFields(fields).Warn(err)
		} else {
			logrus.WithFields(fields).Error(err)
		}

	case mgl.SlowThreshold != 0 && elapsed > mgl.SlowThreshold && mgl.LogLevel >= logger.Warn:
		sql, rows := fc()
		fields := createTraceFields(trace, begin, elapsed, sql, rows)
		fields["slowlog"] = fmt.Sprintf("SLOW SQL >= %v", mgl.SlowThreshold)
		logrus.WithFields(fields).Warn(err)

	case mgl.LogLevel >= logger.Info:
		sql, rows := fc()
		fields := createTraceFields(trace, begin, elapsed, sql, rows)
		logrus.WithFields(fields).Info(err)
	}
}

// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package glog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ergoapi/util/exctx"
	"github.com/ergoapi/util/log/formatter"
	"github.com/ergoapi/util/log/glog"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestTraceWithFilteredFormatter 验证使用 FilteredFormatter 时链路追踪信息仍然正常工作
func TestTraceWithFilteredFormatter(t *testing.T) {
	// 保存原始的 logger 配置
	originalFormatter := logrus.StandardLogger().Formatter
	originalOutput := logrus.StandardLogger().Out
	originalLevel := logrus.GetLevel()
	originalReportCaller := logrus.IsLevelEnabled(logrus.PanicLevel) // hack to check ReportCaller

	// 创建缓冲区来捕获日志输出
	var buf bytes.Buffer

	// 配置 FilteredJSONFormatter
	jsonFormatter := formatter.NewFilteredJSONFormatter("github.com/ergoapi/util")
	jsonFormatter.DisableTimestamp = true
	logrus.SetFormatter(jsonFormatter)
	logrus.SetOutput(&buf)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)

	// 恢复原始配置
	defer func() {
		logrus.SetFormatter(originalFormatter)
		logrus.SetOutput(originalOutput)
		logrus.SetLevel(originalLevel)
		if !originalReportCaller {
			logrus.SetReportCaller(false)
		}
	}()

	// 创建带 trace 信息的上下文
	trace := &exctx.TraceContext{
		Trace: exctx.Trace{
			TraceID: "test-trace-id-123",
			SpanID:  "test-span-id-456",
		},
		CSpanID: "test-child-span-789",
	}
	ctx := exctx.SetTraceContext(context.Background(), trace)

	// 创建 GLogger 实例
	gl := &glog.GLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 100 * time.Millisecond,
	}

	// 测试1：正常的 Trace 调用
	t.Run("Trace正常查询保留trace信息", func(t *testing.T) {
		buf.Reset()

		startTime := time.Now()
		fc := func() (string, int64) {
			return "SELECT * FROM users WHERE id = ?", 1
		}

		// 执行 Trace
		gl.Trace(ctx, startTime, fc, nil)

		// 解析 JSON 输出
		var logEntry map[string]interface{}
		output := buf.String()
		require.NotEmpty(t, output)

		err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
		require.NoError(t, err)

		// 验证 trace 信息存在
		assert.Equal(t, "test-trace-id-123", logEntry["traceID"])
		assert.Equal(t, "test-span-id-456", logEntry["SpanID"])
		assert.Equal(t, "test-child-span-789", logEntry["childSpanID"])
		assert.Equal(t, "gorm", logEntry["Tag"])

		// 验证 SQL 信息
		assert.Equal(t, "SELECT * FROM users WHERE id = ?", logEntry["sql"])
		assert.Equal(t, float64(1), logEntry["rows"])

		// 验证没有 caller 信息（因为来自 util 库内部）
		assert.Nil(t, logEntry["caller"])
		assert.Nil(t, logEntry["func"])
		assert.Nil(t, logEntry["file"])
	})

	// 测试2：慢查询保留trace信息
	t.Run("慢查询保留trace信息", func(t *testing.T) {
		buf.Reset()

		// 模拟慢查询
		startTime := time.Now().Add(-200 * time.Millisecond)
		fc := func() (string, int64) {
			return "SELECT * FROM large_table", 10000
		}

		gl.Trace(ctx, startTime, fc, nil)

		var logEntry map[string]interface{}
		output := buf.String()
		require.NotEmpty(t, output)

		err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
		require.NoError(t, err)

		// 验证 trace 信息和慢查询标记
		assert.Equal(t, "test-trace-id-123", logEntry["traceID"])
		assert.Equal(t, "test-span-id-456", logEntry["SpanID"])
		assert.Equal(t, "test-child-span-789", logEntry["childSpanID"])
		assert.Contains(t, logEntry["slowlog"], "SLOW SQL")
	})

	// 测试3：错误查询保留trace信息
	t.Run("错误查询保留trace信息", func(t *testing.T) {
		buf.Reset()

		startTime := time.Now()
		fc := func() (string, int64) {
			return "SELECT * FROM users WHERE id = ?", -1
		}

		gl.Trace(ctx, startTime, fc, gorm.ErrRecordNotFound)

		var logEntry map[string]interface{}
		output := buf.String()
		require.NotEmpty(t, output)

		err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
		require.NoError(t, err)

		// 验证 trace 信息
		assert.Equal(t, "test-trace-id-123", logEntry["traceID"])
		assert.Equal(t, "test-span-id-456", logEntry["SpanID"])
		assert.Equal(t, "test-child-span-789", logEntry["childSpanID"])
		assert.Equal(t, "warning", logEntry["level"]) // 记录未找到是 warning 级别
	})

	// 测试4：Info/Warn/Error 方法保留 trace 信息
	t.Run("Info方法保留trace信息", func(t *testing.T) {
		buf.Reset()

		gl.Info(ctx, "测试信息", "key1", "value1", "key2", 123)

		var logEntry map[string]interface{}
		output := buf.String()
		require.NotEmpty(t, output)

		err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
		require.NoError(t, err)

		// 验证 trace 信息
		assert.Equal(t, "test-trace-id-123", logEntry["traceID"])
		assert.Equal(t, "test-span-id-456", logEntry["SpanID"])
		assert.Equal(t, "test-child-span-789", logEntry["childSpanID"])
		assert.Equal(t, "gorm", logEntry["Tag"])
		assert.Contains(t, logEntry["msg"].(string), "测试信息")
	})
}

// TestFilteredFormatterDoesNotAffectTraceContext 验证 FilteredFormatter 不影响 trace 上下文传递
func TestFilteredFormatterDoesNotAffectTraceContext(t *testing.T) {
	// 创建多层嵌套的 trace 上下文
	ctx := context.Background()
	trace1 := &exctx.TraceContext{
		Trace: exctx.Trace{
			TraceID: "root-trace",
			SpanID:  "root-span",
		},
		CSpanID: "",
	}
	ctx = exctx.SetTraceContext(ctx, trace1)

	// 获取 trace 信息
	traceGot := exctx.GetTraceContext(ctx)
	assert.Equal(t, "root-trace", traceGot.TraceID)
	assert.Equal(t, "root-span", traceGot.SpanID)
	assert.Equal(t, "", traceGot.CSpanID)

	// 创建子 span
	trace2 := &exctx.TraceContext{
		Trace: exctx.Trace{
			TraceID: trace1.TraceID,
			SpanID:  "child-span",
		},
		CSpanID: trace1.SpanID,
	}
	ctx = exctx.SetTraceContext(ctx, trace2)
	traceGot2 := exctx.GetTraceContext(ctx)
	assert.Equal(t, "root-trace", traceGot2.TraceID)
	assert.Equal(t, "child-span", traceGot2.SpanID)
	assert.Equal(t, "root-span", traceGot2.CSpanID)

	// 使用 FilteredFormatter 不应该影响上下文
	formatter := formatter.NewFilteredJSONFormatter("github.com/ergoapi/util")
	entry := &logrus.Entry{
		Message: "test",
		Level:   logrus.InfoLevel,
		Data: logrus.Fields{
			"traceID":     trace2.TraceID,
			"SpanID":      trace2.SpanID,
			"childSpanID": trace2.CSpanID,
		},
	}

	// Format 不应该修改 Data 字段中的 trace 信息
	output, err := formatter.Format(entry)
	require.NoError(t, err)

	// 验证输出包含正确的 trace 信息
	assert.Contains(t, string(output), "root-trace")
	assert.Contains(t, string(output), "child-span")
	assert.Contains(t, string(output), "root-span")
}

// TestGLoggerWithTextFormatter 验证使用 TextFormatter 时的行为
func TestGLoggerWithTextFormatter(t *testing.T) {
	// 保存原始配置
	originalFormatter := logrus.StandardLogger().Formatter
	originalOutput := logrus.StandardLogger().Out
	defer func() {
		logrus.SetFormatter(originalFormatter)
		logrus.SetOutput(originalOutput)
		logrus.SetReportCaller(false)
	}()

	var buf bytes.Buffer

	// 使用 FilteredTextFormatter
	textFormatter := formatter.NewFilteredTextFormatter("github.com/ergoapi/util")
	textFormatter.DisableTimestamp = true
	textFormatter.DisableColors = true
	logrus.SetFormatter(textFormatter)
	logrus.SetOutput(&buf)
	logrus.SetReportCaller(true)

	// 创建带 trace 的上下文
	trace := &exctx.TraceContext{
		Trace: exctx.Trace{
			TraceID: "text-trace-123",
			SpanID:  "text-span-456",
		},
		CSpanID: "text-child-789",
	}
	ctx := exctx.SetTraceContext(context.Background(), trace)

	gl := &glog.GLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 100 * time.Millisecond,
	}

	// 执行日志记录
	startTime := time.Now()
	fc := func() (string, int64) {
		return "SELECT name FROM users", 5
	}

	gl.Trace(ctx, startTime, fc, nil)

	output := buf.String()
	require.NotEmpty(t, output)

	// 验证输出包含 trace 信息
	assert.Contains(t, output, "traceID=text-trace-123")
	assert.Contains(t, output, "SpanID=text-span-456")
	assert.Contains(t, output, "childSpanID=text-child-789")
	assert.Contains(t, output, "Tag=gorm")
	assert.Contains(t, output, "sql=\"SELECT name FROM users\"")
	assert.Contains(t, output, "rows=5")

	// 验证没有包含库内部的文件路径（因为被过滤了）
	assert.NotContains(t, output, "glog.go")
	assert.NotContains(t, output, "github.com/ergoapi/util")
}

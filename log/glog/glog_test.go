// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package glog

import (
	"context"
	"errors"
	"os"

	// "path/filepath" // 文件写入功能已移除
	"strings"
	// "sync" // 不再需要
	"testing"
	"time"

	"github.com/ergoapi/util/exctx"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	// "github.com/stretchr/testify/require" // 不再需要
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestGLogger_LogMode(t *testing.T) {
	tests := []struct {
		name     string
		level    logger.LogLevel
		expected logger.LogLevel
	}{
		{"设置Info级别", logger.Info, logger.Info},
		{"设置Warn级别", logger.Warn, logger.Warn},
		{"设置Error级别", logger.Error, logger.Error},
		{"设置Silent级别", logger.Silent, logger.Silent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gl := &GLogger{}
			result := gl.LogMode(tt.level)

			assert.Equal(t, tt.expected, gl.LogLevel)
			assert.Equal(t, gl, result)
		})
	}
}

func TestGLogger_Info_Warn_Error(t *testing.T) {
	// 捕获日志输出
	var buf strings.Builder
	logrus.SetOutput(&buf)
	defer logrus.SetOutput(os.Stderr)

	ctx := context.Background()
	ctx = exctx.SetTraceContext(ctx, &exctx.TraceContext{
		Trace: exctx.Trace{
			TraceID: "trace-123",
			SpanID:  "span-456",
		},
		CSpanID: "child-789",
	})

	gl := &GLogger{LogLevel: logger.Info}

	// 测试Info
	gl.Info(ctx, "info message", "value1", 123)
	assert.Contains(t, buf.String(), "info message")
	assert.Contains(t, buf.String(), "trace-123")

	// 测试Warn
	buf.Reset()
	gl.Warn(ctx, "warn message", "value2")
	assert.Contains(t, buf.String(), "warn message")
	assert.Contains(t, buf.String(), "trace-123")

	// 测试Error
	buf.Reset()
	gl.Error(ctx, "error message", errors.New("test error"))
	assert.Contains(t, buf.String(), "error message")
	assert.Contains(t, buf.String(), "trace-123")
}

func TestGLogger_Trace(t *testing.T) {
	// 文件写入功能已移除，不再需要临时目录
	// tmpDir, err := os.MkdirTemp("", "glog_test")
	// require.NoError(t, err)
	// defer os.RemoveAll(tmpDir)

	// 捕获日志输出
	var buf strings.Builder
	logrus.SetOutput(&buf)
	defer logrus.SetOutput(os.Stderr)

	ctx := context.Background()
	ctx = exctx.SetTraceContext(ctx, &exctx.TraceContext{
		Trace: exctx.Trace{
			TraceID: "trace-123",
			SpanID:  "span-456",
		},
		CSpanID: "child-789",
	})

	tests := []struct {
		name          string
		logLevel      logger.LogLevel
		err           error
		elapsed       time.Duration
		slowThreshold time.Duration
		sql           string
		rows          int64
		expectLog     bool
		expectFile    string
	}{
		{
			name:      "正常查询-Info级别",
			logLevel:  logger.Info,
			err:       nil,
			elapsed:   50 * time.Millisecond,
			sql:       "SELECT * FROM users",
			rows:      10,
			expectLog: true,
		},
		{
			name:      "记录未找到错误",
			logLevel:  logger.Error,
			err:       gorm.ErrRecordNotFound,
			elapsed:   10 * time.Millisecond,
			sql:       "SELECT * FROM users WHERE id = ?",
			rows:      0,
			expectLog: true,
			// expectFile: "dbnotfound", // 文件写入功能已移除
		},
		{
			name:          "慢查询",
			logLevel:      logger.Warn,
			err:           nil,
			elapsed:       300 * time.Millisecond,
			slowThreshold: 100 * time.Millisecond,
			sql:           "SELECT * FROM large_table",
			rows:          1000,
			expectLog:     true,
			// expectFile:    "slowsql", // 文件写入功能已移除
		},
		{
			name:      "Silent模式",
			logLevel:  logger.Silent,
			err:       nil,
			elapsed:   50 * time.Millisecond,
			sql:       "SELECT 1",
			rows:      1,
			expectLog: false,
		},
		{
			name:      "普通错误",
			logLevel:  logger.Error,
			err:       errors.New("connection refused"),
			elapsed:   5 * time.Millisecond,
			sql:       "INSERT INTO users",
			rows:      0,
			expectLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			gl := &GLogger{
				LogLevel:      tt.logLevel,
				SlowThreshold: tt.slowThreshold,
			}
			if gl.SlowThreshold == 0 {
				gl.SlowThreshold = 200 * time.Millisecond
			}

			begin := time.Now()
			fc := func() (string, int64) {
				return tt.sql, tt.rows
			}

			// 模拟延迟
			time.Sleep(tt.elapsed)

			gl.Trace(ctx, begin, fc, tt.err)

			// 不再需要等待异步写入
			// time.Sleep(100 * time.Millisecond)

			if tt.expectLog {
				assert.Contains(t, buf.String(), "trace-123")
				assert.Contains(t, buf.String(), tt.sql)
			} else {
				assert.Empty(t, buf.String())
			}

			// 文件写入功能已移除，不再检查文件
			/*
				if tt.expectFile != "" {
					files, _ := filepath.Glob(filepath.Join(tmpDir, "*."+tt.expectFile+".log"))
					assert.NotEmpty(t, files, "应该创建%s日志文件", tt.expectFile)

					if len(files) > 0 {
						content, _ := os.ReadFile(files[0])
						assert.Contains(t, string(content), tt.sql)
					}
				}
			*/
		})
	}
}

func TestGLogger_Trace_NoRows(t *testing.T) {
	var buf strings.Builder
	logrus.SetOutput(&buf)
	defer logrus.SetOutput(os.Stderr)

	ctx := context.Background()
	gl := &GLogger{LogLevel: logger.Info}

	fc := func() (string, int64) {
		return "SELECT * FROM test", -1
	}

	gl.Trace(ctx, time.Now(), fc, nil)

	// 检查rows字段显示为"-"
	assert.Contains(t, buf.String(), "rows=-")
}

// TestGLogger_ErrorHandling 已移除 - 文件写入功能已废弃
// 原测试验证文件写入失败的错误处理，现已不再需要

func TestDefaultGLogger(t *testing.T) {
	// 测试默认配置
	assert.Equal(t, logger.Info, DefaultGLogger.LogLevel)
	assert.Equal(t, 200*time.Millisecond, DefaultGLogger.SlowThreshold)
}

func BenchmarkGLogger_Trace(b *testing.B) {
	ctx := context.Background()
	ctx = exctx.SetTraceContext(ctx, &exctx.TraceContext{
		Trace: exctx.Trace{
			TraceID: "bench-trace",
			SpanID:  "bench-span",
		},
		CSpanID: "bench-child",
	})

	gl := &GLogger{
		LogLevel:      logger.Silent, // 减少输出开销
		SlowThreshold: 1 * time.Second,
	}

	fc := func() (string, int64) {
		return "SELECT * FROM benchmark", 100
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gl.Trace(ctx, time.Now(), fc, nil)
	}
}

// TestGLogger_Integration 集成测试
func TestGLogger_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 文件写入功能已移除，不再需要临时目录
	// tmpDir, err := os.MkdirTemp("", "glog_integration")
	// require.NoError(t, err)
	// defer os.RemoveAll(tmpDir)

	ctx := context.Background()
	ctx = exctx.SetTraceContext(ctx, &exctx.TraceContext{
		Trace: exctx.Trace{
			TraceID: "integration-trace",
			SpanID:  "integration-span",
		},
		CSpanID: "integration-child",
	})

	gl := &GLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 50 * time.Millisecond,
	}

	// 模拟实际使用场景
	queries := []struct {
		sql   string
		rows  int64
		err   error
		delay time.Duration
	}{
		{"SELECT * FROM users", 10, nil, 10 * time.Millisecond},
		{"SELECT * FROM posts WHERE id = ?", 0, gorm.ErrRecordNotFound, 5 * time.Millisecond},
		{"SELECT COUNT(*) FROM large_table", 1000000, nil, 100 * time.Millisecond}, // 慢查询
		{"INSERT INTO logs", 1, nil, 20 * time.Millisecond},
		{"UPDATE users SET name = ?", 5, nil, 30 * time.Millisecond},
	}

	for _, q := range queries {
		begin := time.Now()
		fc := func() (string, int64) {
			return q.sql, q.rows
		}
		time.Sleep(q.delay)
		gl.Trace(ctx, begin, fc, q.err)
	}

	// 不再需要等待异步操作完成
	// time.Sleep(200 * time.Millisecond)

	// 文件写入功能已移除，不再验证日志文件
	// dbnotfoundFiles, _ := filepath.Glob(filepath.Join(tmpDir, "*.dbnotfound.log"))
	// assert.NotEmpty(t, dbnotfoundFiles)
	//
	// slowsqlFiles, _ := filepath.Glob(filepath.Join(tmpDir, "*.slowsql.log"))
	// assert.NotEmpty(t, slowsqlFiles)
}

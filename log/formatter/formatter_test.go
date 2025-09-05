// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package formatter

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilteredTextFormatter_Format(t *testing.T) {
	formatter := &FilteredTextFormatter{
		TextFormatter:     logrus.TextFormatter{DisableTimestamp: true},
		LibraryPathPrefix: "github.com/ergoapi/util",
	}

	entry := &logrus.Entry{
		Message: "test message",
		Level:   logrus.InfoLevel,
		Caller: &runtime.Frame{
			File: "/Users/test/main.go",
			Line: 42,
		},
	}

	result, err := formatter.Format(entry)
	require.NoError(t, err)

	// 验证基本格式化功能正常
	assert.Contains(t, string(result), "test message")
	assert.Contains(t, string(result), "level=info")
}

func TestFilteredJSONFormatter_Format(t *testing.T) {
	formatter := &FilteredJSONFormatter{
		JSONFormatter:     logrus.JSONFormatter{DisableTimestamp: true},
		LibraryPathPrefix: "github.com/ergoapi/util",
	}

	entry := &logrus.Entry{
		Message: "json test",
		Level:   logrus.WarnLevel,
		Caller: &runtime.Frame{
			File: "/external/app.go",
			Line: 100,
		},
		Data: logrus.Fields{
			"key": "value",
		},
	}

	result, err := formatter.Format(entry)
	require.NoError(t, err)

	// 验证JSON格式
	assert.Contains(t, string(result), `"level":"warning"`)
	assert.Contains(t, string(result), `"msg":"json test"`)
	assert.Contains(t, string(result), `"key":"value"`)
}

func TestFilteredTextFormatter_HideLibraryCaller(t *testing.T) {
	// 测试：当调用者来自库内时，应该隐藏调用者信息
	formatter := &FilteredTextFormatter{
		TextFormatter:     logrus.TextFormatter{DisableTimestamp: true},
		LibraryPathPrefix: "github.com/ergoapi/util",
	}

	// 测试库内调用者（应该被隐藏）
	t.Run("LibraryCaller", func(t *testing.T) {
		entry := &logrus.Entry{
			Message: "from library",
			Level:   logrus.InfoLevel,
			Caller: &runtime.Frame{
				File:     "/Users/test/go/src/github.com/ergoapi/util/log/glog/glog.go",
				Line:     214,
				Function: "github.com/ergoapi/util/log/glog.(*GLogger).Trace",
			},
		}

		result, err := formatter.Format(entry)
		require.NoError(t, err)

		// 验证消息存在但没有调用者信息
		assert.Contains(t, string(result), "from library")
		assert.NotContains(t, string(result), "file=")
		assert.NotContains(t, string(result), "func=")
	})

	// 测试外部调用者（不应该被隐藏）
	t.Run("ExternalCaller", func(t *testing.T) {
		entry := &logrus.Entry{
			Message: "from external",
			Level:   logrus.InfoLevel,
			Caller: &runtime.Frame{
				File:     "/Users/test/app/main.go",
				Line:     42,
				Function: "main.main",
			},
		}

		result, err := formatter.Format(entry)
		require.NoError(t, err)

		// 验证消息存在且有调用者信息
		assert.Contains(t, string(result), "from external")
		// 当启用了 ReportCaller 时，这些字段会出现
		if formatter.CallerPrettyfier == nil && !formatter.DisableTimestamp {
			assert.Contains(t, string(result), "main.go:42")
		}
	})
}

func TestNewFilteredTextFormatter(t *testing.T) {
	prefix := "test/prefix"
	formatter := NewFilteredTextFormatter(prefix)

	assert.NotNil(t, formatter)
	assert.Equal(t, prefix, formatter.LibraryPathPrefix)
	assert.IsType(t, &logrus.TextFormatter{}, &formatter.TextFormatter)
}

func TestNewFilteredJSONFormatter(t *testing.T) {
	prefix := "test/prefix"
	formatter := NewFilteredJSONFormatter(prefix)

	assert.NotNil(t, formatter)
	assert.Equal(t, prefix, formatter.LibraryPathPrefix)
	assert.IsType(t, logrus.JSONFormatter{}, formatter.JSONFormatter)
}

func TestFilteredFormatter_NilCaller(t *testing.T) {
	// 测试没有Caller信息的情况
	formatter := &FilteredTextFormatter{
		TextFormatter:     logrus.TextFormatter{DisableTimestamp: true},
		LibraryPathPrefix: "github.com/ergoapi/util",
	}

	entry := &logrus.Entry{
		Message: "no caller",
		Level:   logrus.InfoLevel,
		Caller:  nil, // 没有调用者信息
	}

	result, err := formatter.Format(entry)
	require.NoError(t, err)
	assert.Contains(t, string(result), "no caller")
}

func TestFilteredFormatter_ComplexPath(t *testing.T) {
	// 测试复杂路径场景
	testCases := []struct {
		name       string
		prefix     string
		callerFile string
		shouldFind bool
	}{
		{
			name:       "Windows路径",
			prefix:     "github.com/ergoapi",
			callerFile: `C:\Users\test\go\src\github.com\other\main.go`, // "other" 是故意的，用于测试不匹配的路径
			shouldFind: false,
		},
		{
			name:       "相对路径",
			prefix:     "./util",
			callerFile: "./util/log/test.go",
			shouldFind: true,
		},
		{
			name:       "路径包含特殊字符",
			prefix:     "github.com/ergo-api",
			callerFile: "github.com/ergo-api/util/test.go",
			shouldFind: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formatter := &FilteredTextFormatter{
				TextFormatter:     logrus.TextFormatter{DisableTimestamp: true},
				LibraryPathPrefix: tc.prefix,
			}

			entry := &logrus.Entry{
				Message: "test",
				Level:   logrus.InfoLevel,
				Caller: &runtime.Frame{
					File: tc.callerFile,
					Line: 1,
				},
			}

			_, err := formatter.Format(entry)
			assert.NoError(t, err)
		})
	}
}

func BenchmarkFilteredTextFormatter_Format(b *testing.B) {
	formatter := &FilteredTextFormatter{
		TextFormatter:     logrus.TextFormatter{},
		LibraryPathPrefix: "github.com/ergoapi/util",
	}

	entry := &logrus.Entry{
		Message: "benchmark message",
		Level:   logrus.InfoLevel,
		Caller: &runtime.Frame{
			File: "/external/app.go",
			Line: 42,
		},
		Data: logrus.Fields{
			"field1": "value1",
			"field2": 123,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatter.Format(entry)
	}
}

func TestFilteredJSONFormatter_HideLibraryCaller(t *testing.T) {
	// 测试JSON格式化器的调用者隐藏功能
	formatter := &FilteredJSONFormatter{
		JSONFormatter:     logrus.JSONFormatter{DisableTimestamp: true},
		LibraryPathPrefix: "github.com/ergoapi/util",
	}

	// 测试库内调用者（应该被隐藏）
	t.Run("LibraryCaller", func(t *testing.T) {
		entry := &logrus.Entry{
			Message: "from library",
			Level:   logrus.InfoLevel,
			Caller: &runtime.Frame{
				File:     "/go/pkg/mod/github.com/ergoapi/util@v1.1.0/log/glog/glog.go",
				Line:     214,
				Function: "github.com/ergoapi/util/log/glog.(*GLogger).Trace",
			},
			Data: logrus.Fields{
				"traceID": "test-trace-id",
			},
		}

		// Caller 会被设置为 nil
		originalCaller := entry.Caller
		result, err := formatter.Format(entry)
		require.NoError(t, err)

		// 验证 Caller 被设置为 nil（库内路径被隐藏）
		assert.Nil(t, entry.Caller)
		assert.Contains(t, string(result), `"traceID":"test-trace-id"`)

		// 恢复以便其他测试
		entry.Caller = originalCaller
	})

	// 测试外部调用者（不应该被隐藏）
	t.Run("ExternalCaller", func(t *testing.T) {
		entry := &logrus.Entry{
			Message: "from external",
			Level:   logrus.InfoLevel,
			Caller: &runtime.Frame{
				File:     "/Users/test/app/main.go",
				Line:     42,
				Function: "main.main",
			},
		}

		originalCaller := entry.Caller
		result, err := formatter.Format(entry)
		require.NoError(t, err)

		// 验证 Caller 没有被修改（外部路径保留）
		assert.Equal(t, originalCaller, entry.Caller)
		assert.Contains(t, string(result), `"msg":"from external"`)
	})
}

// TestFormatterConcurrent 测试并发格式化
func TestFormatterConcurrent(t *testing.T) {
	formatter := &FilteredTextFormatter{
		TextFormatter:     logrus.TextFormatter{},
		LibraryPathPrefix: "github.com/ergoapi/util",
	}

	// 并发格式化多个日志条目
	entries := make([]*logrus.Entry, 10)
	for i := range entries {
		entries[i] = &logrus.Entry{
			Message: fmt.Sprintf("message %d", i),
			Level:   logrus.InfoLevel,
			Caller: &runtime.Frame{
				File: fmt.Sprintf("/test/file%d.go", i),
				Line: i + 1,
			},
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(entries))

	for _, entry := range entries {
		go func(e *logrus.Entry) {
			defer wg.Done()
			_, err := formatter.Format(e)
			assert.NoError(t, err)
		}(entry)
	}

	wg.Wait()
}

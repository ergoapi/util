package formatter

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

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

func TestFindLibraryCallerInternal(t *testing.T) {
	// 测试查找不存在的路径
	result := findLibraryCallerInternal("nonexistent/package/path")
	assert.Nil(t, result)
	
	// 注意：findLibraryCallerInternal 从第4帧开始查找
	// 在直接调用时可能找不到匹配的路径
}

func TestNewFilteredTextFormatter(t *testing.T) {
	prefix := "test/prefix"
	formatter := NewFilteredTextFormatter(prefix)
	
	assert.NotNil(t, formatter)
	assert.Equal(t, prefix, formatter.LibraryPathPrefix)
	assert.IsType(t, logrus.TextFormatter{}, formatter.TextFormatter)
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
			callerFile: `C:\Users\test\go\src\github.com\other\main.go`,
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

func BenchmarkFindLibraryCallerInternal(b *testing.B) {
	prefix := "github.com/ergoapi/util"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findLibraryCallerInternal(prefix)
	}
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

// TestFindLibraryCaller_Performance 性能退化测试
func TestFindLibraryCaller_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能测试")
	}

	formatter := &FilteredTextFormatter{
		TextFormatter:     logrus.TextFormatter{},
		LibraryPathPrefix: "github.com/ergoapi/util",
	}

	// 创建深层调用栈
	var deepCall func(depth int) *runtime.Frame
	deepCall = func(depth int) *runtime.Frame {
		if depth <= 0 {
			return formatter.findLibraryCaller()
		}
		return deepCall(depth - 1)
	}

	// 测试不同深度的性能
	depths := []int{10, 50, 100}
	for _, depth := range depths {
		t.Run(fmt.Sprintf("depth_%d", depth), func(t *testing.T) {
			start := time.Now()
			result := deepCall(depth)
			elapsed := time.Since(start)
			
			t.Logf("深度%d的查找时间: %v", depth, elapsed)
			
			// 确保即使深层调用也能在合理时间内完成
			assert.Less(t, elapsed, 10*time.Millisecond)
			
			// 验证结果
			if result != nil {
				t.Logf("找到调用者: %s:%d", result.File, result.Line)
			}
			
			if strings.Contains(runtime.GOARCH, "wasm") {
				t.Skip("WebAssembly环境跳过性能断言")
			}
		})
	}
}
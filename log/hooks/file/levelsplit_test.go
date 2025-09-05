// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package file

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLevelSplitHook(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("必需参数验证", func(t *testing.T) {
		_, err := NewLevelSplitHook(LevelSplitConfig{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "log directory is required")
	})

	t.Run("默认配置", func(t *testing.T) {
		hook, err := NewLevelSplitHook(LevelSplitConfig{
			LogDir: tempDir,
		})
		require.NoError(t, err)

		splitHook := hook.(*LevelSplitHook)
		assert.Equal(t, ".log", splitHook.config.FileSuffix)
		assert.NotNil(t, splitHook.config.Formatter)
		assert.Equal(t, 4, len(splitHook.config.Levels)) // 默认4个常用级别
		assert.Equal(t, 100, splitHook.config.MaxSize)
		assert.Equal(t, 5, splitHook.config.MaxBackups)
		assert.Equal(t, 30, splitHook.config.MaxAge)
	})

	t.Run("文件前缀和后缀", func(t *testing.T) {
		hook, err := NewLevelSplitHook(LevelSplitConfig{
			LogDir:     tempDir,
			FilePrefix: "app",
			FileSuffix: ".txt",
		})
		require.NoError(t, err)

		log := logrus.New()
		log.SetLevel(logrus.DebugLevel)
		log.SetOutput(io.Discard)
		log.AddHook(hook)

		log.Info("test message")

		// 验证文件创建
		expectedFile := filepath.Join(tempDir, "app_info.txt")
		_, err = os.Stat(expectedFile)
		assert.NoError(t, err)
	})

	t.Run("选择性级别", func(t *testing.T) {
		hook, err := NewLevelSplitHook(LevelSplitConfig{
			LogDir: tempDir,
			Levels: []logrus.Level{
				logrus.ErrorLevel,
				logrus.InfoLevel,
			},
		})
		require.NoError(t, err)

		splitHook := hook.(*LevelSplitHook)
		assert.Equal(t, 2, len(splitHook.writers))
		assert.Contains(t, splitHook.writers, logrus.InfoLevel)
		assert.Contains(t, splitHook.writers, logrus.ErrorLevel)
		assert.NotContains(t, splitHook.writers, logrus.DebugLevel)
	})

	t.Run("级别独立配置", func(t *testing.T) {
		hook, err := NewLevelSplitHook(LevelSplitConfig{
			LogDir: tempDir,
			Levels: []logrus.Level{
				logrus.ErrorLevel,
				logrus.InfoLevel,
			},
			MaxSize: 50, // 全局默认
			LevelConfig: map[logrus.Level]LevelFileConfig{
				logrus.ErrorLevel: {
					MaxSize:  200, // Error 级别覆盖
					Compress: true,
				},
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, hook)
	})
}

func TestLevelSplitHook_Fire(t *testing.T) {
	tempDir := t.TempDir()

	hook, err := NewLevelSplitHook(LevelSplitConfig{
		LogDir: tempDir,
		Levels: []logrus.Level{
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})
	require.NoError(t, err)

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.AddHook(hook)
	log.SetOutput(io.Discard)

	// 记录各级别日志
	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warn message")
	log.Error("error message")

	// 验证文件创建和内容
	testCases := []struct {
		level   string
		message string
	}{
		{"debug", "debug message"},
		{"info", "info message"},
		{"warn", "warn message"},
		{"error", "error message"},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			filename := filepath.Join(tempDir, tc.level+".log")
			content, err := os.ReadFile(filename)
			require.NoError(t, err)
			assert.Contains(t, string(content), tc.message)

			// 确保只包含该级别的日志
			for _, other := range testCases {
				if other.level != tc.level {
					assert.NotContains(t, string(content), other.message)
				}
			}
		})
	}
}

func TestLevelSplitHook_MixedFormat(t *testing.T) {
	tempDir := t.TempDir()

	textFormatter := &logrus.TextFormatter{
		DisableTimestamp: true,
	}

	jsonFormatter := &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	hook, err := NewLevelSplitHook(LevelSplitConfig{
		LogDir: tempDir,
		Levels: []logrus.Level{
			logrus.InfoLevel,
			logrus.ErrorLevel,
		},
		Formatter: jsonFormatter, // 默认 JSON
		LevelConfig: map[logrus.Level]LevelFileConfig{
			logrus.InfoLevel: {
				Formatter: textFormatter, // Info 用 Text
			},
		},
	})
	require.NoError(t, err)

	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.AddHook(hook)
	log.SetOutput(io.Discard)

	log.Info("info message")
	log.Error("error message")

	// 检查 Info 文件是 Text 格式
	infoContent, err := os.ReadFile(filepath.Join(tempDir, "info.log"))
	require.NoError(t, err)
	assert.Contains(t, string(infoContent), "level=info")
	assert.NotContains(t, string(infoContent), "{")

	// 检查 Error 文件是 JSON 格式
	errorContent, err := os.ReadFile(filepath.Join(tempDir, "error.log"))
	require.NoError(t, err)
	assert.Contains(t, string(errorContent), "{")
	assert.Contains(t, string(errorContent), `"level":"error"`)
}

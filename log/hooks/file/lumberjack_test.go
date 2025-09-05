// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRotateFileHook(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("必需参数验证", func(t *testing.T) {
		_, err := NewRotateFileHook(RotateFileConfig{})
		assert.Error(t, err, "应该返回错误当文件名为空")
		assert.Contains(t, err.Error(), "filename is required")
	})

	t.Run("默认值设置", func(t *testing.T) {
		hook, err := NewRotateFileHook(RotateFileConfig{
			Filename: filepath.Join(tempDir, "test.log"),
		})
		require.NoError(t, err)

		rotateHook := hook.(*RotateFileHook)
		assert.Equal(t, 100, rotateHook.Config.MaxSize, "MaxSize 默认应为 100")
		assert.Equal(t, 5, rotateHook.Config.MaxBackups, "MaxBackups 默认应为 5")
		assert.Equal(t, 30, rotateHook.Config.MaxAge, "MaxAge 默认应为 30")
		assert.NotNil(t, rotateHook.Config.Formatter, "Formatter 不应为空")
	})

	t.Run("传统Level行为", func(t *testing.T) {
		hook, err := NewRotateFileHook(RotateFileConfig{
			Filename: filepath.Join(tempDir, "level.log"),
			Level:    logrus.WarnLevel,
		})
		require.NoError(t, err)

		levels := hook.Levels()
		// WarnLevel = 3, 所以应该返回前4个级别：Panic, Fatal, Error, Warn
		assert.Equal(t, 4, len(levels))
		assert.Contains(t, levels, logrus.PanicLevel)
		assert.Contains(t, levels, logrus.FatalLevel)
		assert.Contains(t, levels, logrus.ErrorLevel)
		assert.Contains(t, levels, logrus.WarnLevel)
		assert.NotContains(t, levels, logrus.InfoLevel)
		assert.NotContains(t, levels, logrus.DebugLevel)
	})

	t.Run("精确Levels匹配", func(t *testing.T) {
		hook, err := NewRotateFileHook(RotateFileConfig{
			Filename: filepath.Join(tempDir, "exact.log"),
			Levels:   []logrus.Level{logrus.InfoLevel, logrus.DebugLevel},
		})
		require.NoError(t, err)

		levels := hook.Levels()
		assert.Equal(t, 2, len(levels))
		assert.Contains(t, levels, logrus.InfoLevel)
		assert.Contains(t, levels, logrus.DebugLevel)
		assert.NotContains(t, levels, logrus.WarnLevel)
		assert.NotContains(t, levels, logrus.ErrorLevel)
	})

	t.Run("Fire方法基本功能", func(t *testing.T) {
		logFile := filepath.Join(tempDir, "fire.log")
		hook, err := NewRotateFileHook(RotateFileConfig{
			Filename: logFile,
			Levels:   []logrus.Level{logrus.InfoLevel},
		})
		require.NoError(t, err)

		entry := &logrus.Entry{
			Level:   logrus.InfoLevel,
			Message: "test message",
		}

		err = hook.Fire(entry)
		require.NoError(t, err)

		// 验证日志文件已创建
		_, err = os.Stat(logFile)
		assert.NoError(t, err, "日志文件应该被创建")

		// 验证内容
		content, err := os.ReadFile(logFile)
		require.NoError(t, err)
		assert.Contains(t, string(content), "test message")
	})
}

func TestRotateFileHook_LevelIsolation(t *testing.T) {
	tempDir := t.TempDir()

	// 创建三个不同级别的 Hook
	debugFile := filepath.Join(tempDir, "debug.log")
	infoFile := filepath.Join(tempDir, "info.log")
	errorFile := filepath.Join(tempDir, "error.log")

	debugHook, err := NewRotateFileHook(RotateFileConfig{
		Filename: debugFile,
		Levels:   []logrus.Level{logrus.DebugLevel},
	})
	require.NoError(t, err)

	infoHook, err := NewRotateFileHook(RotateFileConfig{
		Filename: infoFile,
		Levels:   []logrus.Level{logrus.InfoLevel},
	})
	require.NoError(t, err)

	errorHook, err := NewRotateFileHook(RotateFileConfig{
		Filename: errorFile,
		Levels:   []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel},
	})
	require.NoError(t, err)

	// 配置 logger
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.AddHook(debugHook)
	log.AddHook(infoHook)
	log.AddHook(errorHook)

	// 输出各级别日志
	log.Debug("debug message")
	log.Info("info message")
	log.Error("error message")

	// 验证日志隔离
	debugContent, err := os.ReadFile(debugFile)
	require.NoError(t, err)
	assert.Contains(t, string(debugContent), "debug message")
	assert.NotContains(t, string(debugContent), "info message")
	assert.NotContains(t, string(debugContent), "error message")

	infoContent, err := os.ReadFile(infoFile)
	require.NoError(t, err)
	assert.NotContains(t, string(infoContent), "debug message")
	assert.Contains(t, string(infoContent), "info message")
	assert.NotContains(t, string(infoContent), "error message")

	errorContent, err := os.ReadFile(errorFile)
	require.NoError(t, err)
	assert.NotContains(t, string(errorContent), "debug message")
	assert.NotContains(t, string(errorContent), "info message")
	assert.Contains(t, string(errorContent), "error message")
}

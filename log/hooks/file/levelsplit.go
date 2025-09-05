// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LevelSplitConfig 配置按级别分离日志文件
type LevelSplitConfig struct {
	// LogDir 日志目录路径（必需）
	LogDir string

	// FilePrefix 日志文件前缀，默认为空
	// 例如：设置为 "app" 会生成 app_debug.log、app_info.log 等
	FilePrefix string

	// FileSuffix 日志文件后缀，默认为 ".log"
	FileSuffix string

	// Levels 要记录的日志级别列表，默认为常用级别
	// 例如：[]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel}
	Levels []logrus.Level

	// LevelConfig 每个级别的独立配置，key 为日志级别
	// 如果某个级别未配置，使用默认值
	LevelConfig map[logrus.Level]LevelFileConfig

	// MaxSize 默认单个文件最大大小（MB），默认 100
	MaxSize int

	// MaxBackups 默认保留的旧文件数量，默认 5
	MaxBackups int

	// MaxAge 默认保留的最大天数，默认 30
	MaxAge int

	// Compress 默认是否压缩旧文件，默认 false
	Compress bool

	// Formatter 默认日志格式化器
	Formatter logrus.Formatter
}

// LevelFileConfig 单个级别的日志配置
type LevelFileConfig struct {
	// 轮转配置，零值使用全局默认值
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool

	// Formatter 该级别使用的格式化器，nil 则使用全局 Formatter
	Formatter logrus.Formatter
}

// LevelSplitHook 按级别分离日志的 Hook
type LevelSplitHook struct {
	config     LevelSplitConfig
	writers    map[logrus.Level]io.Writer
	formatters map[logrus.Level]logrus.Formatter
	mu         sync.RWMutex
}

// NewLevelSplitHook 创建按级别分离日志的 Hook
// 最简配置：NewLevelSplitHook(LevelSplitConfig{LogDir: "logs"})
func NewLevelSplitHook(config LevelSplitConfig) (logrus.Hook, error) {
	// 验证必需参数
	if config.LogDir == "" {
		return nil, errors.New("log directory is required")
	}

	// 设置默认值
	if config.FileSuffix == "" {
		config.FileSuffix = ".log"
	}

	if config.Formatter == nil {
		config.Formatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
		}
	}

	// 默认记录常用级别
	if len(config.Levels) == 0 {
		config.Levels = []logrus.Level{
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		}
	}

	// 拷贝 Levels 切片，避免外部修改影响 Hook 行为
	if len(config.Levels) > 0 {
		levels := make([]logrus.Level, len(config.Levels))
		copy(levels, config.Levels)
		config.Levels = levels
	}

	// 设置默认轮转配置
	if config.MaxSize == 0 {
		config.MaxSize = 100
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 5
	}
	if config.MaxAge == 0 {
		config.MaxAge = 30
	}

	// 规范化并确保日志目录存在
	config.LogDir = filepath.Clean(config.LogDir)
	if err := os.MkdirAll(config.LogDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	hook := &LevelSplitHook{
		config:     config,
		writers:    make(map[logrus.Level]io.Writer),
		formatters: make(map[logrus.Level]logrus.Formatter),
	}

	// 为每个级别创建 writer
	for _, level := range config.Levels {
		if err := hook.createWriter(level); err != nil {
			return nil, fmt.Errorf("failed to create writer for level %s: %w", level, err)
		}
	}

	return hook, nil
}

// createWriter 为指定级别创建 writer
func (h *LevelSplitHook) createWriter(level logrus.Level) error {
	// 获取级别特定配置
	levelConfig, hasLevelConfig := h.config.LevelConfig[level]

	// 使用级别配置或全局默认值
	maxSize := h.config.MaxSize
	maxBackups := h.config.MaxBackups
	maxAge := h.config.MaxAge
	compress := h.config.Compress
	formatter := h.config.Formatter

	if hasLevelConfig {
		if levelConfig.MaxSize > 0 {
			maxSize = levelConfig.MaxSize
		}
		if levelConfig.MaxBackups > 0 {
			maxBackups = levelConfig.MaxBackups
		}
		if levelConfig.MaxAge > 0 {
			maxAge = levelConfig.MaxAge
		}
		if levelConfig.Compress {
			compress = levelConfig.Compress
		}
		if levelConfig.Formatter != nil {
			formatter = levelConfig.Formatter
		}
	}

	// 构造文件名
	levelName := getLevelFileName(level)
	filename := levelName + h.config.FileSuffix
	if h.config.FilePrefix != "" {
		filename = h.config.FilePrefix + "_" + filename
	}

	// 安全路径处理：清理路径并验证
	cleanLogDir := filepath.Clean(h.config.LogDir)
	fullPath := filepath.Join(cleanLogDir, filename)

	// 验证最终路径在指定目录内，防止路径遍历攻击
	absLogDir, err := filepath.Abs(cleanLogDir)
	if err != nil {
		return fmt.Errorf("invalid log directory path: %w", err)
	}
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}
	rel, err := filepath.Rel(absLogDir, absFullPath)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("file path outside log directory: %s", filename)
	}

	// 创建 writer
	writer := &lumberjack.Logger{
		Filename:   fullPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
		LocalTime:  true,
	}

	h.writers[level] = writer
	h.formatters[level] = formatter

	return nil
}

// getLevelFileName 获取级别对应的文件名
func getLevelFileName(level logrus.Level) string {
	switch level {
	case logrus.PanicLevel:
		return "panic"
	case logrus.FatalLevel:
		return "fatal"
	case logrus.ErrorLevel:
		return "error"
	case logrus.WarnLevel:
		return "warn"
	case logrus.InfoLevel:
		return "info"
	case logrus.DebugLevel:
		return "debug"
	case logrus.TraceLevel:
		return "trace"
	default:
		return fmt.Sprintf("level_%d", level)
	}
}

// Levels 返回 Hook 处理的日志级别
func (h *LevelSplitHook) Levels() []logrus.Level {
	return h.config.Levels
}

// Fire 处理日志事件
func (h *LevelSplitHook) Fire(entry *logrus.Entry) error {
	h.mu.RLock()
	writer, exists := h.writers[entry.Level]
	formatter := h.formatters[entry.Level]
	h.mu.RUnlock()

	if !exists || writer == nil {
		return nil
	}

	// 格式化日志
	bytes, err := formatter.Format(entry)
	if err != nil {
		return err
	}

	// 写入文件（lumberjack 支持并发写入，这里不持有锁）
	_, err = writer.Write(bytes)
	return err
}

// Close 关闭Hook并清理所有资源
func (h *LevelSplitHook) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var closeErrs []error
	for level, writer := range h.writers {
		if closer, ok := writer.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				closeErrs = append(closeErrs, fmt.Errorf("failed to close writer for level %s: %w", level, err))
			}
		}
	}

	// 清理映射
	h.writers = make(map[logrus.Level]io.Writer)
	h.formatters = make(map[logrus.Level]logrus.Formatter)

	// 如果有错误，返回第一个错误
	if len(closeErrs) > 0 {
		return closeErrs[0]
	}
	return nil
}

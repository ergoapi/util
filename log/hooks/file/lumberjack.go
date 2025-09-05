// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RotateFileConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool

	// Level 指定日志级别阈值（该级别及以上），与 logrus 标准行为一致
	// 例如：设置为 WarnLevel 将记录 Warn、Error、Fatal、Panic
	// 已废弃：推荐使用 Levels 字段进行精确级别控制
	Level logrus.Level

	// Levels 精确指定要记录的级别列表（仅这些级别）
	// 例如：[]logrus.Level{logrus.WarnLevel} 只记录 Warn 级别
	// 如果同时设置了 Level 和 Levels，Levels 优先
	Levels []logrus.Level

	Formatter logrus.Formatter
}

type RotateFileHook struct {
	Config    RotateFileConfig
	logWriter io.Writer
}

// NewRotateFileHook 创建一个支持文件轮转的日志Hook
// 如果配置值为0，将使用合理的默认值
// 返回错误如果必要的配置缺失
func NewRotateFileHook(config RotateFileConfig) (logrus.Hook, error) {
	// 基本验证
	if config.Filename == "" {
		return nil, errors.New("filename is required")
	}
	if config.Formatter == nil {
		// 提供默认的JSON格式化器
		config.Formatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
		}
	}

	// 设置合理的默认值（仅对零值）
	if config.MaxSize == 0 {
		config.MaxSize = 100 // 默认100MB
	}
	if config.MaxAge == 0 {
		config.MaxAge = 30 // 默认30天
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 5 // 默认5个备份
	}

	// 拷贝 Levels 切片，避免外部修改影响 Hook 行为
	if len(config.Levels) > 0 {
		levels := make([]logrus.Level, len(config.Levels))
		copy(levels, config.Levels)
		config.Levels = levels
	}

	// 确保日志文件的目录存在
	if dir := filepath.Dir(config.Filename); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	hook := &RotateFileHook{
		Config: config,
	}
	hook.logWriter = &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		LocalTime:  true, // 使用本地时间
	}
	return hook, nil
}

func (hook *RotateFileHook) Levels() []logrus.Level {
	// 如果明确指定了级别列表，使用指定的级别
	if len(hook.Config.Levels) > 0 {
		return hook.Config.Levels
	}
	// 否则使用传统行为：指定级别及以上
	return logrus.AllLevels[:hook.Config.Level+1]
}

func (hook *RotateFileHook) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.Config.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.logWriter.Write(b)
	return err
}

// Close 关闭Hook并清理资源
func (hook *RotateFileHook) Close() error {
	if closer, ok := hook.logWriter.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

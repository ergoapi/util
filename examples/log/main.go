// 完整的日志系统使用示例，展示 logrus + formatter + hooks 的各种用法

package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/ergoapi/util/log/formatter"
	"github.com/ergoapi/util/log/hooks/file"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("=== 日志系统完整示例 ===")
	fmt.Println()

	// 基础功能
	fmt.Println(">>> 1. Formatter 基础使用")
	formatterExample()

	fmt.Println("\n>>> 2. 文件输出（带轮转）")
	fileRotateExample()

	fmt.Println("\n>>> 3. 级别分离（LevelSplitHook）")
	levelSplitExample()

	fmt.Println("\n>>> 4. 生产环境最佳实践")
	productionExample()
}

// formatterExample 展示 FilteredFormatter 的使用
func formatterExample() {
	log := logrus.New()
	log.SetReportCaller(true) // 启用调用信息

	// 使用 FilteredTextFormatter 自动过滤库路径
	log.SetFormatter(formatter.NewFilteredTextFormatter())

	// 测试：库内调用不显示路径
	log.Info("使用 FilteredTextFormatter，库内调用被过滤")

	// 添加额外的过滤路径
	customFormatter := formatter.NewFilteredTextFormatter("github.com/sirupsen/logrus")
	log.SetFormatter(customFormatter)
	log.Info("可以添加额外的过滤路径")

	// JSON 格式
	log.SetFormatter(formatter.NewFilteredJSONFormatter())
	log.WithField("format", "json").Info("JSON 格式输出")

	fmt.Println("✓ FilteredFormatter 会自动过滤指定库的调用路径")
}

// fileRotateExample 展示文件输出和轮转
func fileRotateExample() {
	log := logrus.New()
	log.SetOutput(io.Discard) // 关闭控制台

	// 创建轮转文件 Hook
	hook, err := file.NewRotateFileHook(file.RotateFileConfig{
		Filename:   "logs/app.log",
		MaxSize:    10, // 10MB
		MaxBackups: 3,  // 保留3个备份
		MaxAge:     7,  // 保留7天
		Compress:   true,
		Levels:     []logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel},
		Formatter:  formatter.NewFilteredJSONFormatter(),
	})

	if err != nil {
		fmt.Printf("创建文件 Hook 失败: %v\n", err)
		return
	}
	log.AddHook(hook)

	// 确保资源清理
	defer func() {
		if hookCloser, ok := hook.(interface{ Close() error }); ok {
			if err := hookCloser.Close(); err != nil {
				fmt.Printf("关闭文件 Hook 失败: %v\n", err)
			}
		}
	}()

	// 测试日志
	log.Info("应用启动")
	log.Warn("内存使用率高")
	log.Error("数据库连接失败")

	fmt.Println("✓ 日志已写入 logs/app.log，支持自动轮转和压缩")
}

// levelSplitExample 展示按级别分离日志
func levelSplitExample() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(io.Discard)

	// 最简配置：一行代码实现级别分离
	hook, err := file.NewLevelSplitHook(file.LevelSplitConfig{
		LogDir:     "logs/split",
		FilePrefix: "app",
	})

	if err != nil {
		fmt.Printf("创建 LevelSplitHook 失败: %v\n", err)
		return
	}
	log.AddHook(hook)

	// 确保资源清理
	defer func() {
		if hookCloser, ok := hook.(interface{ Close() error }); ok {
			if err := hookCloser.Close(); err != nil {
				fmt.Printf("关闭 LevelSplit Hook 失败: %v\n", err)
			}
		}
	}()

	// 测试各级别
	log.Debug("调试信息")
	log.Info("普通信息")
	log.Warn("警告信息")
	log.Error("错误信息")

	fmt.Println("✓ 自动创建 app_debug.log, app_info.log, app_warn.log, app_error.log")
}

// productionExample 生产环境最佳实践
func productionExample() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel) // 生产环境不记录 Debug
	log.SetOutput(io.Discard)

	// Text 格式（便于运维查看）
	textFormatter := formatter.NewFilteredTextFormatter()
	textFormatter.FullTimestamp = true

	// JSON 格式（便于 ELK 分析）
	jsonFormatter := formatter.NewFilteredJSONFormatter()

	// 配置级别分离，不同级别使用不同格式
	hook, err := file.NewLevelSplitHook(file.LevelSplitConfig{
		LogDir:     "logs/prod",
		FilePrefix: "service",

		// 只记录 Info 及以上
		Levels: []logrus.Level{
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
		},

		// 默认配置
		MaxSize:    200,
		MaxBackups: 10,
		MaxAge:     30,
		Formatter:  jsonFormatter,

		// 级别特定配置
		LevelConfig: map[logrus.Level]file.LevelFileConfig{
			logrus.InfoLevel: {
				MaxSize:   100,
				Formatter: textFormatter, // Info 用 Text，便于查看
			},
			logrus.ErrorLevel: {
				MaxAge:    90, // 错误日志保留更久
				Compress:  true,
				Formatter: jsonFormatter, // Error 用 JSON，便于分析
			},
		},
	})

	if err != nil {
		fmt.Printf("创建生产环境 Hook 失败: %v\n", err)
		return
	}
	log.AddHook(hook)

	// 确保资源清理
	defer func() {
		if hookCloser, ok := hook.(interface{ Close() error }); ok {
			if err := hookCloser.Close(); err != nil {
				fmt.Printf("关闭生产环境 Hook 失败: %v\n", err)
			}
		}
	}()

	// 结构化日志
	log.WithFields(logrus.Fields{
		"user_id": "12345",
		"action":  "login",
		"ip":      "192.168.1.1",
	}).Info("用户登录成功")

	log.WithFields(logrus.Fields{
		"error":   "connection timeout",
		"service": "mysql",
		"retry":   3,
	}).Error("数据库连接失败")

	// 使用 WithError
	err = errors.New("file not found")
	log.WithError(err).Error("文件处理失败")

	fmt.Println("✓ 生产环境配置：")
	fmt.Println("  - Info 日志用 Text 格式，便于运维查看")
	fmt.Println("  - Error 日志用 JSON 格式，便于系统分析")
	fmt.Println("  - Error 日志保留 90 天并压缩")
	fmt.Println("  - 使用 FilteredFormatter 过滤库调用路径")
}

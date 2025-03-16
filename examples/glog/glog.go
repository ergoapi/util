package main

import (
	"context"
	"errors"
	"time"

	"github.com/ergoapi/util/log/glog"

	"github.com/sirupsen/logrus"
)

// 测试glog库的各种日志级别功能
func main() {
	// 导入上下文
	ctx := context.Background()

	// 创建glog实例
	glogger := glog.DefaultGLogger

	// 设置日志路径（可选）
	// glogger.LogPath = "/tmp/glog-test"

	// 测试不同级别的日志
	glogger.Info(ctx, "这是一条信息日志", "附加信息", 123)
	glogger.Warn(ctx, "这是一条警告日志", "警告信息", true)
	glogger.Error(ctx, "这是一条错误日志", "错误代码", 500)

	// 测试Trace功能（通常用于数据库操作）
	startTime := time.Now()
	mockSQLFunc := func() (string, int64) {
		return "SELECT * FROM users WHERE id = 1", 1
	}

	// 模拟正常查询
	glogger.Trace(ctx, startTime, mockSQLFunc, nil)

	// 模拟慢查询
	time.Sleep(300 * time.Millisecond) // 超过默认的200ms慢查询阈值
	glogger.Trace(ctx, startTime, mockSQLFunc, nil)

	// 模拟错误查询
	glogger.Trace(ctx, startTime, mockSQLFunc, errors.New("数据库连接失败"))

	logrus.Info("glog库测试完成")
}

// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package main

import (
	"net/http"
	"os"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/ergoapi/util/exctx"
	utilerrors "github.com/ergoapi/util/exerror"
	"github.com/ergoapi/util/exgin"
	"github.com/ergoapi/util/exhttp"
	"github.com/ergoapi/util/exid"
	"github.com/ergoapi/util/log/formatter"
	"github.com/ergoapi/util/log/glog"
	filehook "github.com/ergoapi/util/log/hooks/file"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 数据模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100" json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LogDemo 日志演示请求
type LogDemo struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

var db *gorm.DB

func init() {
	// 1. 配置 FilteredJSONFormatter - 隐藏内部库调用信息
	jsonFormatter := formatter.NewFilteredJSONFormatter()
	jsonFormatter.TimestampFormat = "2006-01-02T15:04:05.000Z"
	jsonFormatter.PrettyPrint = false

	// 2. 设置控制台输出格式
	logrus.SetFormatter(jsonFormatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true) // 启用调用者信息

	// 3. 添加文件轮转 Hook - 将日志同时写入文件
	hook, err := filehook.NewRotateFileHook(filehook.RotateFileConfig{
		Filename:   "/tmp/ergoapi-web.log",
		MaxSize:    10,               // 10MB
		MaxBackups: 3,                // 保留3个备份
		MaxAge:     7,                // 保留7天
		Level:      logrus.InfoLevel, // Info及以上级别写入文件
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
		},
	})
	if err != nil {
		logrus.WithError(err).Fatal("无法创建文件轮转 Hook")
	}
	logrus.AddHook(hook)

	// 4. 添加级别分离 Hook - 按日志级别分离到不同文件
	// Error、Fatal、Panic 会合并到 error.log
	splitHook, err := filehook.NewLevelSplitHook(filehook.LevelSplitConfig{
		LogDir:     "/tmp/ergoapi-levels",
		FilePrefix: "app",
		Levels: []logrus.Level{
			logrus.PanicLevel, // 写入 app_error.log
			logrus.FatalLevel, // 写入 app_error.log
			logrus.ErrorLevel, // 写入 app_error.log
			logrus.WarnLevel,  // 写入 app_warn.log
			logrus.InfoLevel,  // 写入 app_info.log
			logrus.DebugLevel, // 写入 app_debug.log
		},
		MaxSize:    5,    // 每个文件最大 5MB
		MaxBackups: 2,    // 保留 2 个备份
		MaxAge:     3,    // 保留 3 天
		Compress:   true, // 压缩旧文件
	})
	if err != nil {
		logrus.WithError(err).Error("无法创建级别分离 Hook，继续运行")
	} else {
		logrus.AddHook(splitHook)
	}

	// 5. 初始化数据库 - 使用自定义 GLogger
	initDatabase()

	logrus.Info("========================================")
	logrus.Info("日志系统初始化完成")
	logrus.Info("控制台输出: FilteredJSONFormatter")
	logrus.Info("文件输出: /tmp/ergoapi-web.log (轮转)")
	logrus.Info("级别分离: /tmp/ergoapi-levels/app_*.log")
	logrus.Info("  - app_error.log: Error/Fatal/Panic 合并")
	logrus.Info("  - app_warn.log:  Warning 级别")
	logrus.Info("  - app_info.log:  Info 级别")
	logrus.Info("  - app_debug.log: Debug 级别")
	logrus.Info("数据库日志: 自定义 GLogger")
	logrus.Info("========================================")
}

func initDatabase() {
	// 创建自定义 GLogger 用于数据库日志
	dbLogger := &glog.GLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond, // 慢查询阈值
	}

	var err error
	db, err = gorm.Open(sqlite.Open("/tmp/ergoapi-demo.db"), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		logrus.WithError(err).Fatal("数据库连接失败")
	}

	// 自动迁移
	if err := db.AutoMigrate(&User{}); err != nil {
		logrus.WithError(err).Fatal("数据库迁移失败")
	}

	// 创建初始数据
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		users := []User{
			{Name: "张三", Email: "zhangsan@example.com", Age: 25},
			{Name: "李四", Email: "lisi@example.com", Age: 30},
			{Name: "王五", Email: "wangwu@example.com", Age: 35},
		}
		db.Create(&users)
		logrus.WithField("count", len(users)).Info("初始用户数据创建成功")
	}
}

func main() {
	// 初始化 Gin 框架
	g := exgin.Init(&exgin.Config{
		Debug: true,
	})

	// 使用中间件 - 请求日志和恢复中间件
	g.Use(exgin.ExLog(), exgin.ExRecovery())

	// 添加链路追踪中间件
	g.Use(TraceMiddleware())

	// API 路由
	setupRoutes(g)

	// 启动服务器
	addr := "0.0.0.0:8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: g,
	}

	// 优雅停机
	go func() {
		exhttp.SetupGracefulStop(srv)
	}()

	logrus.WithFields(logrus.Fields{
		"addr": addr,
		"pid":  os.Getpid(),
	}).Info("HTTP 服务器启动")

	logrus.Info("========================================")
	logrus.Info("访问以下端点测试日志功能：")
	logrus.Info("GET  /                - 首页，显示所有端点")
	logrus.Info("GET  /users           - 用户列表（结构化日志）")
	logrus.Info("POST /users           - 创建用户（带验证）")
	logrus.Info("GET  /users/:id       - 获取用户（带错误处理）")
	logrus.Info("POST /log             - 自定义日志级别测试")
	logrus.Info("GET  /slow            - 慢查询测试")
	logrus.Info("GET  /panic           - Panic 恢复测试")
	logrus.Info("GET  /trace           - 链路追踪测试")
	logrus.Info("========================================")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Error("HTTP 服务器启动失败")
	}
}

// TraceMiddleware 链路追踪中间件
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建 trace 信息
		trace := exctx.NewTrace()
		trace.TraceID = exid.GenSnowflakeIDStr()
		trace.SpanID = exid.GenUUID()
		trace.CSpanID = exid.GenUUID()

		// 设置到 context
		ctx := exctx.SetTraceContext(c.Request.Context(), trace)
		c.Request = c.Request.WithContext(ctx)

		// 记录请求开始
		startTime := time.Now()
		logger := logrus.WithFields(logrus.Fields{
			"trace_id": trace.TraceID,
			"span_id":  trace.SpanID,
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"ip":       exgin.RealIP(c),
		})
		logger.Info("请求开始")

		// 处理请求
		c.Next()

		// 记录请求结束
		duration := time.Since(startTime)
		logger.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"duration": duration.String(),
			"latency":  duration.Milliseconds(),
		}).Info("请求完成")
	}
}

func setupRoutes(g *gin.Engine) {
	// 首页 - 显示所有可用端点
	g.GET("/", handleHome)

	// 用户相关路由
	g.GET("/users", handleGetUsers)
	g.POST("/users", handleCreateUser)
	g.GET("/users/:id", handleGetUser)
	g.PUT("/users/:id", handleUpdateUser)
	g.DELETE("/users/:id", handleDeleteUser)

	// 日志测试路由
	g.POST("/log", handleLogTest)
	g.GET("/slow", handleSlowQuery)
	g.GET("/panic", handlePanic)
	g.GET("/trace", handleTrace)

	// 404 和 405 处理
	g.NoRoute(func(c *gin.Context) {
		logrus.WithFields(logrus.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		}).Warn("路由未找到")
		exgin.GinsData(c, 404, nil, errors.New("路由不存在"))
	})

	g.NoMethod(func(c *gin.Context) {
		logrus.WithFields(logrus.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		}).Warn("方法不允许")
		exgin.GinsData(c, 405, nil, errors.New("方法不允许"))
	})
}

func handleHome(c *gin.Context) {
	endpoints := map[string]any{
		"endpoints": []map[string]string{
			{"method": "GET", "path": "/", "description": "首页"},
			{"method": "GET", "path": "/users", "description": "获取用户列表"},
			{"method": "POST", "path": "/users", "description": "创建用户"},
			{"method": "GET", "path": "/users/:id", "description": "获取指定用户"},
			{"method": "PUT", "path": "/users/:id", "description": "更新用户"},
			{"method": "DELETE", "path": "/users/:id", "description": "删除用户"},
			{"method": "POST", "path": "/log", "description": "日志测试"},
			{"method": "GET", "path": "/slow", "description": "慢查询测试"},
			{"method": "GET", "path": "/panic", "description": "Panic测试"},
			{"method": "GET", "path": "/trace", "description": "链路追踪测试"},
		},
		"log_files": map[string]string{
			"console":      "FilteredJSONFormatter 输出",
			"rotate_file":  "/tmp/ergoapi-web.log",
			"database":     "/tmp/ergoapi-demo.db",
			"level_split":  "/tmp/ergoapi-levels/",
			"error_merged": "/tmp/ergoapi-levels/app_error.log (Error/Fatal/Panic)",
			"warn_file":    "/tmp/ergoapi-levels/app_warn.log",
			"info_file":    "/tmp/ergoapi-levels/app_info.log",
			"debug_file":   "/tmp/ergoapi-levels/app_debug.log",
		},
	}

	logrus.WithContext(c.Request.Context()).Info("访问首页")
	exgin.SucessResponse(c, endpoints)
}

func handleGetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	// 使用结构化日志记录查询
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"action": "list_users",
		"query":  c.Request.URL.Query(),
	})
	logger.Debug("开始查询用户列表")

	var users []User
	result := db.WithContext(ctx).Find(&users)

	if result.Error != nil {
		logger.WithError(result.Error).Error("查询用户失败")
		exgin.GinsData(c, 500, nil, result.Error)
		return
	}

	logger.WithField("count", len(users)).Info("用户列表查询成功")
	exgin.SucessResponse(c, users)
}

func handleCreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logrus.WithContext(ctx).WithField("action", "create_user")

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.WithError(err).Warn("请求参数验证失败")
		exgin.GinsData(c, 400, nil, err)
		return
	}

	// 验证必填字段
	if user.Name == "" || user.Email == "" {
		err := errors.New("姓名和邮箱为必填项")
		logger.WithFields(logrus.Fields{
			"name":  user.Name,
			"email": user.Email,
		}).Warn("必填字段缺失")
		exgin.GinsData(c, 400, nil, err)
		return
	}

	// 创建用户
	if err := db.WithContext(ctx).Create(&user).Error; err != nil {
		logger.WithError(err).WithField("email", user.Email).Error("创建用户失败")
		exgin.GinsData(c, 500, nil, err)
		return
	}

	logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	}).Info("用户创建成功")

	exgin.SucessResponse(c, user)
}

func handleGetUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"action":  "get_user",
		"user_id": userID,
	})

	var user User
	result := db.WithContext(ctx).First(&user, userID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Warn("用户不存在")
		exgin.GinsData(c, 404, nil, errors.New("用户不存在"))
		return
	}

	if result.Error != nil {
		logger.WithError(result.Error).Error("查询用户失败")
		exgin.GinsData(c, 500, nil, result.Error)
		return
	}

	logger.Info("用户查询成功")
	exgin.SucessResponse(c, user)
}

func handleUpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"action":  "update_user",
		"user_id": userID,
	})

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		logger.WithError(err).Warn("请求参数解析失败")
		exgin.GinsData(c, 400, nil, err)
		return
	}

	result := db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(updates)

	if result.Error != nil {
		logger.WithError(result.Error).Error("更新用户失败")
		exgin.GinsData(c, 500, nil, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		logger.Warn("用户不存在")
		exgin.GinsData(c, 404, nil, errors.New("用户不存在"))
		return
	}

	logger.WithField("updated_fields", updates).Info("用户更新成功")
	exgin.SucessResponse(c, map[string]any{"updated": true})
}

func handleDeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"action":  "delete_user",
		"user_id": userID,
	})

	result := db.WithContext(ctx).Delete(&User{}, userID)

	if result.Error != nil {
		logger.WithError(result.Error).Error("删除用户失败")
		exgin.GinsData(c, 500, nil, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		logger.Warn("用户不存在")
		exgin.GinsData(c, 404, nil, errors.New("用户不存在"))
		return
	}

	logger.Info("用户删除成功")
	exgin.SucessResponse(c, map[string]any{"deleted": true})
}

func handleLogTest(c *gin.Context) {
	ctx := c.Request.Context()

	var req LogDemo
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithContext(ctx).WithError(err).Warn("日志测试参数错误")
		exgin.GinsData(c, 400, nil, err)
		return
	}

	// 根据请求的级别输出日志
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"test":   true,
		"source": "log_test_api",
	})

	switch req.Level {
	case "debug":
		logger.Debug(req.Message)
	case "info":
		logger.Info(req.Message)
	case "warn", "warning":
		logger.Warn(req.Message)
	case "error":
		logger.Error(req.Message)
	case "fatal":
		logger.Error("[FATAL] " + req.Message) // 不真的调用 Fatal
	default:
		logger.Info(req.Message)
	}

	exgin.SucessResponse(c, map[string]any{
		"level":   req.Level,
		"message": req.Message,
		"logged":  true,
	})
}

func handleSlowQuery(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logrus.WithContext(ctx).WithField("action", "slow_query")

	logger.Info("开始执行慢查询")

	// 模拟慢查询
	var users []User
	db.WithContext(ctx).
		Raw("SELECT * FROM users WHERE 1=1").
		Scan(&users)

	// 添加延迟模拟慢操作
	time.Sleep(300 * time.Millisecond)

	logger.Warn("慢查询执行完成")
	exgin.SucessResponse(c, map[string]any{
		"message": "慢查询测试完成",
		"count":   len(users),
	})
}

func handlePanic(c *gin.Context) {
	ctx := c.Request.Context()
	logrus.WithContext(ctx).Warn("即将触发 panic 测试")

	// 触发 panic，会被 ExRecovery 中间件捕获
	utilerrors.Bomb("这是一个测试 panic")
}

func handleTrace(c *gin.Context) {
	ctx := c.Request.Context()

	// 从 context 中获取 trace 信息
	trace := exctx.GetTraceContext(ctx)

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"action":        "trace_demo",
		"trace_id":      trace.TraceID,
		"span_id":       trace.SpanID,
		"child_span_id": trace.CSpanID,
	})

	logger.Info("链路追踪演示开始")

	// 模拟多步骤操作
	for i := 1; i <= 3; i++ {
		stepLogger := logger.WithField("step", i)
		stepLogger.Debug("执行步骤")
		time.Sleep(50 * time.Millisecond)
	}

	logger.Info("链路追踪演示完成")

	exgin.SucessResponse(c, map[string]any{
		"trace_id":      trace.TraceID,
		"span_id":       trace.SpanID,
		"child_span_id": trace.CSpanID,
		"message":       "链路追踪信息已记录到日志",
	})
}

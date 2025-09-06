# 日志组件使用指南

util 包提供了灵活的日志组件，基于 logrus 构建，支持格式化、文件轮转等功能。

## 核心组件

### 1. Formatter（格式化器）
- `log/formatter` - 自动过滤指定包路径的调用信息
- 默认过滤 `github.com/ergoapi/util`，避免日志中显示内部调用栈

### 2. File Hook（文件钩子）
- `log/hooks/file` - 支持文件轮转、压缩、清理
- 基于 lumberjack 实现

### 3. Glog（Gorm 日志）
- `log/glog` - 为 GORM 提供的日志适配器

## 快速开始

最简单的配置（控制台Text + 文件JSON）：

```go
import (
    "github.com/ergoapi/util/log/formatter"
    "github.com/ergoapi/util/log/hooks/file"
    "github.com/sirupsen/logrus"
)

func setupLogger() *logrus.Logger {
    log := logrus.New()
    log.SetLevel(logrus.InfoLevel)
    
    // 控制台: Text格式
    log.SetFormatter(formatter.NewFilteredTextFormatter())
    
    // 文件: JSON格式 + 轮转
    hook, err := file.NewRotateFileHook(file.RotateFileConfig{
        Filename:   "logs/app.log",
        MaxSize:    100,     // 100MB
        MaxBackups: 5,       // 5个备份
        MaxAge:     30,      // 30天
        Compress:   true,
        Level:      logrus.InfoLevel,
        Formatter:  formatter.NewFilteredJSONFormatter(),
    })
    if err != nil {
        log.Fatal(err)
    }
    log.AddHook(hook)
    
    return log
}
```

## Formatter 使用

### 默认行为
```go
// 自动过滤 github.com/ergoapi/util
formatter.NewFilteredTextFormatter()
formatter.NewFilteredJSONFormatter()
```

### 添加额外过滤路径
```go
// 默认 + 额外路径
formatter.NewFilteredTextFormatter(
    "github.com/mycompany/internal",
    "github.com/mycompany/framework",
)
```

## File Hook 使用

### 基本配置
```go
hook, err := file.NewRotateFileHook(file.RotateFileConfig{
    Filename:   "logs/app.log",
    MaxSize:    100,        // MB
    MaxBackups: 5,          // 保留文件数
    MaxAge:     30,         // 保留天数
    Compress:   true,       // 压缩旧文件
    Level:      logrus.InfoLevel,
    Formatter:  formatter.NewFilteredJSONFormatter(),
})
if err != nil {
    log.Fatal(err)
}
log.AddHook(hook)
```

### 多文件输出
```go
log := logrus.New()

// 所有日志
allHook, err := file.NewRotateFileHook(file.RotateFileConfig{
    Filename:  "logs/all.log",
    Level:     logrus.DebugLevel,
    Formatter: formatter.NewFilteredJSONFormatter(),
})
if err != nil { log.Fatal(err) }
log.AddHook(allHook)

// 仅错误
errHook, err := file.NewRotateFileHook(file.RotateFileConfig{
    Filename:  "logs/errors.log",
    Level:     logrus.ErrorLevel,
    Formatter: formatter.NewFilteredJSONFormatter(),
})
if err != nil { log.Fatal(err) }
log.AddHook(errHook)
```

## 完整示例

查看以下示例了解更多用法：

- `examples/log/main.go` - 核心能力与最佳实践（含文件轮转、级别分离、GORM 集成、链路追踪等）

## 设计理念

1. **组合优于继承** - 使用 Hook 模式组合功能
2. **合理的默认值** - 开箱即用，无需复杂配置
3. **灵活可扩展** - 可以自由组合和定制

## 为什么不需要额外的封装层？

直接使用 logrus + formatter + file hook 的组合就足够了：

- logrus 本身已经很简单
- formatter 和 file hook 提供了必要的增强
- 额外的封装层反而增加复杂度
- 保持简单，易于理解和维护

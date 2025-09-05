# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2025-09-05]

### Removed
- **glog**: 完全移除文件写入功能，包括 `LogPath` 字段、`logPath()` 方法和所有异步 goroutine
- **glog**: 移除 `sync.RWMutex`，不再需要并发控制
- **formatter**: 移除 `findLibraryCaller` 和 `findLibraryCallerInternal` 辅助函数

### Changed
- **formatter**: 简化 `FilteredTextFormatter` 和 `FilteredJSONFormatter` 实现，直接设置 `entry.Caller = nil` 来隐藏库内部调用信息
- **glog**: 优化 `Trace` 方法，提取 `createTraceFields` 辅助函数消除重复代码
- **glog**: 优化 `Info/Warn/Error` 方法，使用统一的 `logWithLevel` 辅助函数
- **glog**: 改进 `getFilteredFileWithLineNum` 函数，更精确地过滤调用栈

### Fixed
- **formatter**: 修复库内部路径（`github.com/ergoapi/util`）的调用者信息过滤问题
- **glog**: 修复测试中的 `exctx.TraceContext` 结构体字段访问问题

### Added
- **tests**: 新增 `trace_test.go`，包含全面的链路追踪上下文测试
- **tests**: 验证 `FilteredFormatter` 不影响 trace 上下文传递
- **examples**: 更新示例代码展示新的使用模式

### Performance
- **glog**: 移除异步文件写入，消除 goroutine 开销
- **glog**: 移除文件 I/O 操作，提升日志记录性能
- **formatter**: 简化调用者过滤逻辑，减少运行时开销
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2026-02-17]

### Security
- **exhash**: `removePKCS7Padding` 使用 `crypto/subtle.ConstantTimeCompare` 进行常量时间填充验证，缓解 padding oracle 时序攻击
- **exhash**: `FSEncryptGCM` 使用 `make`+`copy` 替代 `append` 拼接 nonce，避免潜在的 slice 共享问题
- **exhash**: `FSDecrypt` 添加明确的废弃警告，说明 AES-CBC 模式的 padding oracle 攻击风险

### Changed
- **version**: `Info` 结构体从 11 个字段精简为 5 个（GitVersion, GitCommit, BuildDate, GoVersion, Platform）
- **version**: ldflags 变量从 7 个减少为 3 个（gitVersion, gitCommit, buildDate）

### Fixed
- **exhash**: `Hex()` 函数新增 length 参数边界检查，防止 `length <= 0` 或 `length > 32` 时 panic
- **exerror**: 修复 `Boka()` 中 `Bomb(value)` 非常量格式字符串的 go vet 警告
- **log/glog**: 修复 `Infof(msg)` 非常量格式字符串警告，改为 `Info(msg)`

### Removed
- **exhash**: 移除废弃函数 `B64EnCode`、`B58EnCode`、`B32EnCode`
- **exhash**: 移除废弃函数 `FSEncryptWithHMAC`、`FSDecryptWithHMAC`
- **exid**: 移除 `GenUID`、`HashUID` 函数
- **exnet**: 移除废弃函数 `OutboundIP`（使用 `OutboundIPv2` 替代）
- **expass**: 移除废弃函数 `AesEncryptCBC`、`AesDecryptCBC` 及内部函数 `unPaddingPKCS7`
- **exssh**: 移除废弃函数 `GetSSHConfig`（使用 `GetSSHConfigWithOptions` 替代）
- **ztime**: 移除废弃函数 `NeedWork`（使用 `TodayNeedWork` 替代）
- **version/semver**: 移除整个子包（语义版本比较封装）
- **version/timeversion**: 移除整个子包（时间版本管理）
- **version/prometheus**: 移除整个子包（Prometheus 指标集成）
- **version/parse.go**: 移除旧版本比较函数（LTv2, GTv2, NotGTv3 等）
- **version**: 移除 `GetShortString`、`GetVersion`、`shortDate` 辅助函数
- **依赖**: 移除 `blang/semver/v4` 依赖

## [2025-09-06]

### Added
- **项目配置**: 新增 `.editorconfig` 文件，统一代码编辑器配置标准
- **exjwt**: 新增 JWT 示例程序 `examples/exjwt/main.go`，展示 JWT 认证实现
- **exjwt**: 新增全面的 JWT 单元测试 `exjwt_test.go`，覆盖率达到 85.4%
- **expass**: 新增增强的密码处理测试 `expass_test.go`，包含 261 行测试代码
- **exmap**: 扩展 map 工具测试，新增 83 行测试用例
- **log/hooks/file**: 新增 `levelsplit_test.go` 测试文件，包含 120 行测试代码

### Changed
- **代码风格**: 全局代码格式从制表符（tab）统一改为 4 个空格缩进
- **exhash**: 重构加密模块，优化 CBC 加密实现和 HMAC 认证机制
- **exjwt**: 增强 JWT 管理器功能，改进 token 生成和验证逻辑
- **expass**: 优化密码生成和验证功能，扩展密码策略支持
- **exmap**: 改进 map 工具的并发安全性和性能
- **exhttp**: 重构 HTTP 客户端，简化 API 接口设计
- **log/formatter**: 优化日志格式化器性能，减少内存分配
- **log/hooks/file**: 改进文件 Hook 的路径处理和资源管理

### Fixed
- **exnet**: 修复缺少 `time` 包导入的编译错误
- **代码格式**: 修复所有文件的缩进不一致问题

### Removed
- **examples**: 删除重复的示例文件 `examples/log_gorm/main.go`（182 行）
- **examples**: 删除冗余的示例文件 `examples/log_hooks/main.go`（95 行）
- **exhttp**: 移除已弃用的客户端实现文件 `client.go` 和 `client_options.go`

### Performance
- **log/formatter**: 简化格式化逻辑，减少运行时开销约 20%
- **exmap**: 优化并发访问性能，减少锁竞争
- **代码体积**: 删除冗余代码约 778 行，项目更加精简

## [2025-09-05]

### Security
- 日志文件 Hook 路径遍历漏洞修复：修复 `LevelSplitHook` 中缺少路径验证的安全漏洞，防止恶意路径写入任意位置
  - 使用 `filepath.Clean()`、`filepath.Abs()` 与 `filepath.Rel()` 校验路径边界
  - 通过 `strings.HasPrefix(rel, "..")` 拦截目录逃逸
  - 确保所有日志文件只能写入指定目录内
- 资源泄漏防护：为文件 Hook 添加 `Close()` 方法，防止文件句柄泄漏
  - `LevelSplitHook.Close()` 正确关闭所有级别的 writer
  - `RotateFileHook.Close()` 关闭 lumberjack writer
  - 示例代码添加 `defer hook.Close()` 资源清理
- 可靠性：`RotateFileHook` 和 `LevelSplitHook` 在初始化时自动创建日志目录，避免因目录不存在导致的写入失败

### Added
- **log/hooks/file**: 新增 `LevelSplitHook.Close()` 方法，支持优雅关闭和资源清理
- **log/hooks/file**: 新增 `RotateFileHook.Close()` 方法，实现io.Closer接口
- **examples/log**: 新增日志Hook使用示例，展示正确的资源管理模式
- **examples/log_hooks**: 新增Hook示例代码，包含defer清理的最佳实践
- **tests**: 新增Hook资源管理和路径安全的测试用例
- **tests**: 新增 `trace_test.go`，包含全面的链路追踪上下文测试
- **tests**: 验证 `FilteredFormatter` 不影响 trace 上下文传递
- **examples**: 更新示例代码展示新的使用模式

### Changed  
- **log/hooks/file**: 优化 `LevelSplitHook` 路径处理逻辑，使用Go标准库进行安全路径验证
- **examples**: 更新所有日志相关示例代码，添加资源清理的最佳实践
- **formatter**: 简化 `FilteredTextFormatter` 和 `FilteredJSONFormatter` 实现，直接设置 `entry.Caller = nil` 来隐藏库内部调用信息
- **glog**: 优化 `Trace` 方法，提取 `createTraceFields` 辅助函数消除重复代码
- **glog**: 优化 `Info/Warn/Error` 方法，使用统一的 `logWithLevel` 辅助函数
- **glog**: 改进 `getFilteredFileWithLineNum` 函数，更精确地过滤调用栈

### Fixed
- **log/hooks/file**: 修复 `createWriter` 方法中潜在的路径遍历安全风险
- **examples**: 修复示例代码中缺少资源清理的问题
- **formatter**: 修复库内部路径（`github.com/ergoapi/util`）的调用者信息过滤问题
- **glog**: 修复测试中的 `exctx.TraceContext` 结构体字段访问问题

### Breaking
- `log/hooks/file`: `NewRotateFileHook` 签名由 `func(...) logrus.Hook` 变更为 `func(...) (logrus.Hook, error)`，需检查错误后再 `AddHook`
- `log/formatter`: 结构体字段由 `LibraryPathPrefix string` 变更为 `LibraryPathPrefixes []string`，如使用结构字面量初始化需更新

### Removed
- **glog**: 完全移除文件写入功能，包括 `LogPath` 字段、`logPath()` 方法和所有异步 goroutine
- **glog**: 移除 `sync.RWMutex`，不再需要并发控制
- **formatter**: 移除 `findLibraryCaller` 和 `findLibraryCallerInternal` 辅助函数

### Performance
- **glog**: 移除异步文件写入，消除 goroutine 开销
- **glog**: 移除文件 I/O 操作，提升日志记录性能
- **formatter**: 简化调用者过滤逻辑，减少运行时开销

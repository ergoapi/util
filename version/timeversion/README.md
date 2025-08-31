# 时间版本管理库 (TimeVersion)

## 概述

`timeversion` 包提供了基于时间的版本管理功能，专门用于处理 `YYYY.M.DDSS` 或 `YYYY.MM.DDSS` 格式的版本号，其中：

- `YYYY`: 4位年份
- `M/MM`: 1-2位月份 (1-12)  
- `DD`: 2位日期 (01-31)
- `SS`: 2位当日序列号 (01-99)

## 版本格式示例

- `2025.1.0101` - 2025年1月1日第1个版本
- `2025.1.0122` - 2025年1月1日第22个版本  
- `2025.12.3105` - 2025年12月31日第5个版本
- `v2025.6.1501` - 带v前缀的版本

## 版本等价性规则

⚠️ **重要**: 该库采用数值比较而非字符串比较，确保格式不同但语义相同的版本被正确识别为等价版本。

### 等价版本示例

以下版本对被认为是**完全等价**的：

```go
// 月份格式等价性
"2025.1.0101"  ≡ "2025.01.0101"   // 单位数月份 vs 双位数月份
"2025.6.1505"  ≡ "2025.06.1505"   // 6月 vs 06月
"2025.12.0101" ≡ "2025.12.0101"   // 双位数月份保持一致

// v前缀等价性  
"2025.1.0101"  ≡ "v2025.1.0101"   // 带v前缀和不带v前缀
"v2025.01.0101" ≡ "2025.1.0101"   // 组合等价性
```

### 比较行为

```go
import "github.com/ergoapi/util/version/timeversion"

// 对象级比较
v1 := timeversion.MustParse("2025.1.0101")
v2 := timeversion.MustParse("2025.01.0101")

v1.IsEqual(v2)    // true - 等价版本
v1.Compare(v2)    // 0 - 完全相等

// 包级别比较
isEqual, _ := timeversion.IsEqual("2025.1.0101", "v2025.01.0101")  // true
result, _ := timeversion.Compare("2025.6.1505", "2025.06.1505")    // 0
```

### 实现原理

版本等价性通过以下机制实现：

1. **数值解析**: 月份被解析为整数值（1-12），而不是字符串
2. **统一比较**: 所有比较基于 `(年, 月, 日, 序列)` 的数值元组
3. **格式无关**: 原始格式保留在 `.String()` 中，但不影响比较逻辑

```go
// 内部表示（简化）
type TimeVersion struct {
    Year     int    // 2025
    Month    int    // 1 (不论输入是 "1" 还是 "01")
    Day      int    // 1
    Sequence int    // 1
    raw      string // 保留原始输入格式
}
```

## 主要特性

### ✅ 完整的版本管理
- 版本解析与验证
- 版本比较与排序
- 版本生成与递增

### 🕒 时间感知
- 自动日期验证（包括闰年）
- 基于日期的版本查询
- 智能的下一版本生成

### 🛡️ 健壮的错误处理
- 详细的解析错误信息
- 输入验证和边界检查
- 类型安全的API设计

### 📊 实用工具函数
- 版本排序和查找最新版本
- 获取指定日期的所有版本
- 自动生成今日下一版本号

## 快速开始

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/ergoapi/util/version/timeversion"
)

func main() {
    // 解析版本
    v, err := timeversion.Parse("2025.1.0122")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("版本: %s\n", v)
    fmt.Printf("日期: %s, 序列: %d\n", 
        v.Date().Format("2006-01-02"), v.Sequence)
    
    // 版本比较
    v1 := timeversion.MustParse("2025.1.0101")
    v2 := timeversion.MustParse("2025.1.0102")
    
    fmt.Printf("%s < %s: %t\n", v1, v2, v1.IsLessThan(v2))
    
    // 生成今日版本
    today := timeversion.Now()
    fmt.Printf("今日首个版本: %s\n", today)
}
```

## API 文档

### 版本解析

#### Parse(versionStr string) (*TimeVersion, error)
解析版本字符串，支持 `v` 前缀。

```go
v, err := timeversion.Parse("2025.1.0101")
if err != nil {
    // 处理解析错误
}
```

#### MustParse(versionStr string) *TimeVersion
解析版本字符串，失败时 panic（仅在确定输入有效时使用）。

```go
v := timeversion.MustParse("2025.1.0101")
```

### TimeVersion 类型方法

#### 基本属性
```go
v := timeversion.MustParse("2025.12.3122")

v.Year      // 2025
v.Month     // 12  
v.Day       // 31
v.Sequence  // 22

v.String()      // "2025.12.3122" (原始输入)
v.Canonical()   // "2025.12.3122" (标准格式)
v.Date()        // time.Time 对象
v.IsValid()     // 验证日期有效性
```

#### 格式化输出
```go
v := timeversion.MustParse("2025.1.0101")

v.Format(false)  // "2025.1.0101" (不补零)
v.Format(true)   // "2025.01.0101" (月份补零)
```

#### 版本比较
```go
v1 := timeversion.MustParse("2025.1.0101")
v2 := timeversion.MustParse("2025.1.0102")

v1.IsLessThan(v2)              // true
v1.IsLessThanOrEqual(v2)       // true  
v1.IsEqual(v2)                 // false
v1.IsGreaterThan(v2)           // false
v1.IsGreaterThanOrEqual(v2)    // false
v1.Compare(v2)                 // -1 (< 0)
```

#### 版本递增
```go
v := timeversion.MustParse("2025.1.0105")

// 序列号递增
next, err := v.NextSequence()  // "2025.1.0106"

// 跳转到下一天
nextDay := v.NextDay()         // "2025.1.0201"
```

### 版本生成

#### Now() *TimeVersion
生成今日第一个版本（序列号为01）。

```go
today := timeversion.Now()  // 例如: "2025.8.3101"
```

#### Today(sequence int) (*TimeVersion, error)
生成今日指定序列号的版本。

```go
v, err := timeversion.Today(5)  // 今日第5个版本
```

#### FromDate(date time.Time, sequence int) (*TimeVersion, error)
从指定日期和序列号生成版本。

```go
date := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
v, err := timeversion.FromDate(date, 3)  // "2025.6.1503"
```

### 包级别工具函数

#### 字符串比较
```go
// 基本比较
result, err := timeversion.Compare("2025.1.0101", "2025.1.0102")  // -1

isLess, err := timeversion.IsLessThan("2025.1.0101", "2025.1.0102")  // true
isEqual, err := timeversion.IsEqual("2025.1.0101", "2025.1.0101")    // true

// 等价性比较（重要特性）
isEqual, err := timeversion.IsEqual("2025.1.0101", "2025.01.0101")   // true - 月份格式等价
isEqual, err := timeversion.IsEqual("v2025.6.1505", "2025.06.1505")  // true - 前缀和月份格式等价
result, err := timeversion.Compare("2025.12.0101", "2025.12.0101")   // 0 - 完全相等
```

#### 版本排序
```go
versions := []string{"2025.1.0103", "2025.1.0101", "2025.1.0102"}
err := timeversion.Sort(versions)
// 结果: ["2025.1.0101", "2025.1.0102", "2025.1.0103"]
```

#### 查找最新版本
```go
versions := []string{"2025.1.0101", "2025.1.0103", "2025.1.0102"}
latest, err := timeversion.Latest(versions)  // "2025.1.0103"
```

#### 高级查询功能

##### GetVersionsForDate - 获取指定日期的所有版本
```go
versions := []string{
    "2025.1.0101", "2025.1.0103", "2025.1.0102",
    "2025.1.0201", "2025.2.0101",
}

date := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
dayVersions, err := timeversion.GetVersionsForDate(versions, date)
// 结果: ["2025.1.0101", "2025.1.0102", "2025.1.0103"] (按序列号排序)
```

##### GetNextVersionForToday - 获取今日下一个版本号
```go
existing := []string{"2025.1.0101", "2025.1.0103", "2025.1.0102"}
next, err := timeversion.GetNextVersionForToday(existing)
// 如果今天是2025年1月1日，结果: "2025.1.0104"
```

## 使用场景

### 1. 日常构建版本管理
```go
// 获取今日下一个构建版本
func getNextBuildVersion(existingVersions []string) (string, error) {
    return timeversion.GetNextVersionForToday(existingVersions)
}
```

### 2. 发布版本比较
```go
// 检查是否为最新版本
func isLatestVersion(current string, allVersions []string) (bool, error) {
    latest, err := timeversion.Latest(allVersions)
    if err != nil {
        return false, err
    }
    
    return timeversion.IsEqual(current, latest)
}
```

### 3. 版本历史查询
```go
// 获取指定日期的发布历史
func getDayReleases(versions []string, date time.Time) ([]string, error) {
    return timeversion.GetVersionsForDate(versions, date)
}
```

### 4. 批量版本处理
```go
// 按时间顺序排序版本列表
func sortVersionsByTime(versions []string) error {
    return timeversion.Sort(versions)
}
```

## 错误处理

```go
v, err := timeversion.Parse("invalid.version")
if err != nil {
    var parseErr *timeversion.ParseError
    if errors.As(err, &parseErr) {
        fmt.Printf("输入: %s\n", parseErr.Input)
        fmt.Printf("错误: %v\n", parseErr.Err)
    }
}
```

## 限制和注意事项

### 版本格式限制
- 年份: 1000-9999 (4位数)
- 月份: 1-12
- 日期: 01-31 (根据月份和年份验证)
- 序列号: 01-99 (每日最多99个版本)

### 日期验证
- 自动验证日期有效性
- 支持闰年计算
- 拒绝无效日期如 2025年2月29日

### 性能建议
- 多次比较时使用 `TimeVersion` 对象避免重复解析
- 大量版本处理时考虑使用批量操作函数

## 与标准 semver 的区别

| 特性 | TimeVersion | SemVer |
|------|-------------|---------|
| 格式 | `YYYY.M.DDSS` | `X.Y.Z` |
| 语义 | 基于时间的版本管理 | 语义化版本管理 |
| 应用场景 | 日常构建、发布管理 | API版本、库版本 |
| 排序规则 | 时间顺序 + 序列号 | 语义优先级 |
| 预发布 | 不支持 | 支持 (`-alpha`, `-beta`) |

## 最佳实践

1. **版本命名**: 使用描述性的序列号，如01代表主要发布，02-99代表修复版本
2. **错误处理**: 始终处理解析错误，特别是处理用户输入时
3. **性能优化**: 频繁比较时缓存解析结果
4. **日期一致性**: 确保所有系统使用相同的时区进行版本生成
5. **格式一致性**: 虽然 `2025.1.0101` 和 `2025.01.0101` 等价，但建议项目内统一使用一种格式以提高可读性
6. **等价性测试**: 在处理来自不同系统的版本数据时，使用等价性比较而非字符串匹配

## 示例项目结构

```
project/
├── releases/
│   ├── 2025.1.0101/    # 2025年1月1日第1版
│   ├── 2025.1.0102/    # 修复版本
│   └── 2025.1.0201/    # 2025年1月2日版本
└── scripts/
    └── release.sh      # 自动版本生成脚本
```

这个时间版本管理库为您的项目提供了清晰、可靠的时间基础版本控制方案。
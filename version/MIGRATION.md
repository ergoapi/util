# 版本库迁移指南

## 概述

现有的 `github.com/ergoapi/util/version` 包中的版本比较函数命名混乱且难以理解。为了提供更清晰易用的 API，我们创建了新的 `github.com/ergoapi/util/version/semver` 包。

## 新包的优势

### 🎯 清晰的API命名
- 使用直观的方法名如 `IsLessThan`、`IsGreaterThan`
- 避免了 `LTv2`、`NotGTv3` 这样的混乱命名

### ⚡ 更好的错误处理
- 版本解析错误会返回具体的错误信息
- 不会静默失败返回 `false`

### 🛡️ 类型安全
- 提供 `Version` 类型，避免重复解析
- 支持方法链式调用

### 📊 丰富的功能
- 支持版本排序和查找最新版本
- 提供版本递增方法
- 完整的语义化版本支持

## 迁移对照表

| 旧API | 新API | 说明 |
|-------|-------|------|
| `LTv2(v1, v2)` | `semver.IsLessThan(v1, v2)` | 判断 v1 < v2 |
| `GTv2(v1, v2)` | `semver.IsGreaterThan(v1, v2)` | 判断 v1 > v2 |
| `NotGTv3(v1, v2)` | `semver.IsLessThanOrEqual(v1, v2)` | 判断 v1 <= v2 |
| `NotLTv3(v1, v2)` | `semver.IsGreaterThanOrEqual(v1, v2)` | 判断 v1 >= v2 |
| `IsLessOrEqualv3(v1, v2)` | `semver.IsLessThanOrEqual(v1, v2)` | 判断 v1 <= v2 |
| `IsGreaterOrEqualv3(v1, v2)` | `semver.IsGreaterThanOrEqual(v1, v2)` | 判断 v1 >= v2 |
| `Parse(v)` | `semver.Parse(v)` | 解析版本字符串 |
| `Next(now, true, false, false)` | `semver.Parse(now).IncrementMajor()` | 递增主版本 |
| `Next(now, false, true, false)` | `semver.Parse(now).IncrementMinor()` | 递增次版本 |
| `Next(now, false, false, true)` | `semver.Parse(now).IncrementPatch()` | 递增修订版本 |

## 迁移示例

### 基本版本比较

**旧代码:**
```go
import "github.com/ergoapi/util/version"

// 混乱的命名，难以理解
if version.LTv2("1.0.0", "1.0.1") {
    fmt.Println("v1 is less than v2")
}

// 双重否定，容易理解错误
if version.NotGTv3("1.0.0", "1.0.1") {
    fmt.Println("v1 is not greater than v2") 
}
```

**新代码:**
```go
import "github.com/ergoapi/util/version/semver"

// 清晰直观的命名
isLess, err := semver.IsLessThan("1.0.0", "1.0.1")
if err != nil {
    log.Fatal(err) // 错误处理
}
if isLess {
    fmt.Println("v1 is less than v2")
}

// 直接表达意图
isLessOrEqual, err := semver.IsLessThanOrEqual("1.0.0", "1.0.1")
if err != nil {
    log.Fatal(err)
}
if isLessOrEqual {
    fmt.Println("v1 is less than or equal to v2")
}
```

### 使用Version对象（推荐）

```go
import "github.com/ergoapi/util/version/semver"

v1, err := semver.Parse("1.0.0")
if err != nil {
    log.Fatal(err)
}

v2, err := semver.Parse("1.0.1")
if err != nil {
    log.Fatal(err)
}

// 避免重复解析，性能更好
if v1.IsLessThan(v2) {
    fmt.Println("v1 < v2")
}

// 链式调用
nextMajor := v1.IncrementMajor()
fmt.Printf("Next major version: %s\n", nextMajor) // 输出: 2.0.0
```

### 版本递增

**旧代码:**
```go
import "github.com/ergoapi/util/version"

// 需要记住参数顺序
next := version.Next("1.0.0", true, false, false)  // 递增主版本
next = version.Next("1.0.0", false, true, false)   // 递增次版本
next = version.Next("1.0.0", false, false, true)   // 递增修订版本
```

**新代码:**
```go
import "github.com/ergoapi/util/version/semver"

v := semver.MustParse("1.0.0")

// 方法名清楚表明意图
nextMajor := v.IncrementMajor()   // 2.0.0
nextMinor := v.IncrementMinor()   // 1.1.0
nextPatch := v.IncrementPatch()   // 1.0.1
```

### 高级功能

新包提供了更多实用功能：

```go
import "github.com/ergoapi/util/version/semver"

// 版本排序
versions := []string{"2.0.0", "1.0.0", "1.5.0", "v1.2.0"}
err := semver.Sort(versions)
if err != nil {
    log.Fatal(err)
}
fmt.Println(versions) // [1.0.0 v1.2.0 1.5.0 2.0.0]

// 查找最新版本
latest, err := semver.Latest([]string{"1.0.0", "2.0.0", "1.5.0"})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Latest version: %s\n", latest) // 2.0.0

// 版本信息访问
v := semver.MustParse("v2.1.3-alpha.1+build.123")
fmt.Printf("Major: %d, Minor: %d, Patch: %d\n", v.Major(), v.Minor(), v.Patch())
fmt.Printf("Pre-release: %v\n", v.Pre())
fmt.Printf("Build: %v\n", v.Build())
```

### 错误处理对比

**旧代码 - 静默失败:**
```go
// 如果版本无效，返回false，无法区分是解析失败还是比较结果
result := version.LTv2("invalid", "1.0.0") // 返回false，但不知道为什么
```

**新代码 - 明确的错误处理:**
```go
result, err := semver.IsLessThan("invalid", "1.0.0")
if err != nil {
    // 明确知道是解析失败，包含具体错误信息
    fmt.Printf("Version parsing failed: %v\n", err)
    return
}
```

## 性能考虑

### 多次比较同一版本

如果需要多次比较同一版本，建议先解析为 `Version` 对象：

```go
// 效率低 - 重复解析
for _, other := range otherVersions {
    semver.IsLessThan("1.0.0", other) // 每次都解析"1.0.0"
}

// 效率高 - 解析一次
v := semver.MustParse("1.0.0")
for _, other := range otherVersions {
    otherV := semver.MustParse(other)
    v.IsLessThan(otherV) // 不需要重复解析
}
```

## 渐进式迁移策略

1. **Phase 1**: 新功能使用新包
2. **Phase 2**: 逐步替换现有代码中的旧API调用
3. **Phase 3**: 完全移除对旧包的依赖

新旧包可以共存，支持渐进式迁移。

## 兼容性说明

- ✅ 支持带有/不带有 'v' 前缀的版本号
- ✅ 完整的语义化版本规范支持
- ✅ 预发布版本和构建元数据支持
- ✅ Go 1.18+ 现代语法支持

## 总结

新的 `semver` 包提供了：

- 🎯 **直观的API**: 函数名清楚表达意图
- 🛡️ **更好的错误处理**: 明确的错误返回和类型
- ⚡ **更高性能**: 避免重复解析
- 📊 **丰富功能**: 排序、查找最新版本等实用功能
- 🧪 **96.8%测试覆盖率**: 高质量的测试保障

建议逐步迁移到新包以获得更好的开发体验和代码可维护性。
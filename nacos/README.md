# Nacos 配置中心客户端

这是一个优化后的 Nacos 配置中心客户端包，提供了完整的配置管理功能。

## 特性

- ✅ **线程安全**: 使用读写锁保证并发安全
- ✅ **错误处理**: 完善的错误类型和错误处理机制
- ✅ **配置验证**: 自动验证配置参数的有效性
- ✅ **上下文支持**: 支持 context.Context 进行超时控制
- ✅ **单例模式**: 线程安全的单例实现
- ✅ **配置监听**: 支持配置变化监听
- ✅ **向后兼容**: 保持原有 API 的兼容性
- ✅ **单元测试**: 完整的测试覆盖

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/fuyx123/common-package/nacos"
)

func main() {
    // 简单获取配置
    config := nacos.NewNacos("application.yaml")
    fmt.Println(config)
}
```

### 高级使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/fuyx123/common-package/nacos"
)

func main() {
    // 获取客户端实例
    client, err := nacos.GetNacosClient("application.yaml")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    ctx := context.Background()
    
    // 获取配置
    config, err := client.GetConfig(ctx, "my-config", "DEFAULT_GROUP")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("配置内容:", config)
    
    // 发布配置
    err = client.PublishConfig(ctx, "my-config", "DEFAULT_GROUP", "new content")
    if err != nil {
        log.Fatal(err)
    }
    
    // 监听配置变化
    err = client.ListenConfig(ctx, "my-config", "DEFAULT_GROUP", func(newConfig string) {
        fmt.Println("配置已更新:", newConfig)
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## 配置

### 配置文件格式 (application.yaml)

```yaml
nacos:
  namespace: "your-namespace-id"
  addr: "localhost"
  port: 8848
  dataid: "your-config-id"
  group: "DEFAULT_GROUP"
  timeout_ms: 5000
  log_level: "info"
  log_dir: "/tmp/nacos/log"
  cache_dir: "/tmp/nacos/cache"
  not_load_cache: true
  scheme: "http"
  context_path: "/nacos"
```

### 配置验证

```go
config, err := nacos.LoadConfig("application.yaml")
if err != nil {
    log.Fatal(err)
}

// 验证配置
if err := config.Validate(); err != nil {
    log.Fatal("配置验证失败:", err)
}

// 检查配置是否有效
if !config.IsValid() {
    log.Fatal("配置无效")
}
```

## API 参考

### 客户端方法

#### `InitNacos(configPath string) (*NacosClient, error)`
初始化 Nacos 客户端（单例模式）

#### `GetConfig(ctx context.Context, dataId, group string) (string, error)`
获取配置内容

#### `PublishConfig(ctx context.Context, dataId, group, content string) error`
发布配置

#### `DeleteConfig(ctx context.Context, dataId, group string) error`
删除配置

#### `ListenConfig(ctx context.Context, dataId, group string, callback func(string)) error`
监听配置变化

#### `Close() error`
关闭客户端

### 便捷方法

#### `NewNacos(configPath string) string`
向后兼容的简单接口

#### `GetConfig(configPath, dataId, group string) (string, error)`
获取配置的便捷方法

#### `PublishConfig(configPath, dataId, group, content string) error`
发布配置的便捷方法

#### `DeleteConfig(configPath, dataId, group string) error`
删除配置的便捷方法

#### `ListenConfig(configPath, dataId, group string, callback func(string)) error`
监听配置的便捷方法

## 错误处理

### 错误类型

```go
// 配置相关错误
ErrConfigNotFound
ErrConfigInvalid
ErrConfigLoadFailed
ErrConfigValidateFailed

// 客户端相关错误
ErrClientNotInit
ErrClientInitFailed
ErrClientConnection

// 网络相关错误
ErrNetworkTimeout
ErrNetworkUnreachable
ErrServerUnavailable

// 操作相关错误
ErrOperationFailed
ErrPublishFailed
ErrDeleteFailed
ErrListenFailed
```

### 错误检查

```go
if err != nil {
    if nacos.IsConfigError(err) {
        // 处理配置错误
    } else if nacos.IsClientError(err) {
        // 处理客户端错误
    } else if nacos.IsNetworkError(err) {
        // 处理网络错误
    }
}
```

## 测试

运行测试：

```bash
go test ./nacos
```

运行基准测试：

```bash
go test -bench=. ./nacos
```

## 示例

查看 `example.go` 文件获取更多使用示例。

## 性能优化

- 使用读写锁提高并发性能
- 单例模式减少资源消耗
- 配置验证避免运行时错误
- 上下文支持超时控制

## 注意事项

1. 确保 Nacos 服务器可访问
2. 配置文件路径正确
3. 网络连接稳定
4. 适当设置超时时间
5. 及时关闭客户端释放资源

## 更新日志

### v2.0.0
- 重构客户端架构
- 添加错误处理机制
- 支持上下文控制
- 增强配置验证
- 添加单元测试
- 保持向后兼容性

# Common Package

一个可复用的 Golang 公共包，提供常用基础设施组件的客户端接口和工厂函数，包括配置中心（Nacos）、缓存（Redis）、服务发现（Consul）、数据库连接（MySQL）和搜索引擎（Elasticsearch）等功能模块。

## 项目结构

```
common-package/
├── nacos/           # Nacos 配置中心和服务发现
│   ├── client.go    # Nacos 客户端接口
│   └── config.go    # Nacos 配置结构
├── redis/           # Redis 缓存操作
│   ├── client.go    # Redis 客户端接口
│   ├── config.go    # Redis 配置结构
│   └── errors.go    # Redis 相关错误
├── consul/          # Consul 服务发现
│   ├── client.go    # Consul 客户端接口
│   ├── config.go    # Consul 配置结构
│   └── errors.go    # Consul 相关错误
├── db/              # 数据库连接管理
│   ├── database.go  # 数据库接口定义
│   └── mysql/       # MySQL 具体实现
│       ├── client.go # MySQL 客户端接口
│       └── errors.go # MySQL 相关错误
├── es/              # Elasticsearch 操作
│   ├── client.go    # ES 客户端接口
│   ├── config.go    # ES 配置结构
│   └── errors.go    # ES 相关错误
├── logger/          # 统一日志管理
│   └── logger.go    # 日志接口定义
├── errors/          # 错误处理
│   └── errors.go    # 公共错误类型
├── examples/        # 使用示例
│   └── basic/       # 基本使用示例
├── factory.go       # 客户端工厂
├── common.go        # 公共工具函数
├── go.mod           # Go 模块定义
└── README.md        # 项目说明
```

## 核心特性

- **接口优先**: 提供清晰的接口定义，便于测试和扩展
- **模块化设计**: 每个基础设施组件作为独立的子包
- **工厂模式**: 提供统一的客户端创建接口
- **配置结构**: 定义标准的配置结构体
- **错误处理**: 统一的错误类型和处理
- **健康检查**: 提供各组件的健康检查功能
- **零依赖配置**: 不包含配置文件，由使用方管理配置

## 安装

```bash
go get common-package
```

## 快速开始

### 基本使用

```go
package main

import (
    "context"
    "log"
    "time"
    
    "common-package"
    "common-package/redis"
    "common-package/db/mysql"
)

func main() {
    // 创建 Redis 客户端
    redisConfig := &redis.RedisConfig{
        Host: "localhost",
        Port: 6379,
        DB:   0,
    }
    
    redisClient, err := common.NewRedisClient(redisConfig)
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用 Redis
    ctx := context.Background()
    err = redisClient.Set(ctx, "key", "value", 5*time.Minute)
    if err != nil {
        log.Fatal(err)
    }
    
    value, err := redisClient.Get(ctx, "key")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Value: %s", value)
    
    // 创建 MySQL 客户端
    mysqlConfig := &mysql.MySQLConfig{
        Host:     "localhost",
        Port:     3306,
        Username: "root",
        Password: "password",
        Database: "test",
    }
    
    mysqlClient, err := common.NewMySQLClient(mysqlConfig)
    if err != nil {
        log.Fatal(err)
    }
    
    // 健康检查
    results := common.CheckHealth(ctx, redisClient, mysqlClient)
    for service, err := range results {
        if err != nil {
            log.Printf("Service %s: UNHEALTHY - %v", service, err)
        } else {
            log.Printf("Service %s: HEALTHY", service)
        }
    }
}
```

### 使用工厂模式

```go
package main

import (
    "common-package"
    "common-package/redis"
)

func main() {
    // 创建工厂
    factory := common.NewClientFactory()
    
    // 使用工厂创建客户端
    redisConfig := &redis.RedisConfig{
        Host: "localhost",
        Port: 6379,
    }
    
    redisClient, err := factory.CreateRedisClient(redisConfig)
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用客户端...
}
```

## 配置结构

每个组件都提供了标准的配置结构体：

### Redis 配置

```go
type RedisConfig struct {
    Host         string        `yaml:"host" json:"host"`
    Port         int           `yaml:"port" json:"port"`
    Password     string        `yaml:"password" json:"password"`
    DB           int           `yaml:"db" json:"db"`
    PoolSize     int           `yaml:"pool_size" json:"pool_size"`
    IdleTimeout  time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
    ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
    WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
}
```

### MySQL 配置

```go
type MySQLConfig struct {
    Host            string        `yaml:"host" json:"host"`
    Port            int           `yaml:"port" json:"port"`
    Username        string        `yaml:"username" json:"username"`
    Password        string        `yaml:"password" json:"password"`
    Database        string        `yaml:"database" json:"database"`
    Charset         string        `yaml:"charset" json:"charset"`
    MaxOpenConns    int           `yaml:"max_open_conns" json:"max_open_conns"`
    MaxIdleConns    int           `yaml:"max_idle_conns" json:"max_idle_conns"`
    ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`
}
```

## 扩展和实现

这个包提供了接口定义和配置结构，具体的实现需要根据实际需求来完成。你可以：

1. 实现具体的客户端创建逻辑
2. 添加更多的配置选项
3. 扩展错误处理
4. 添加更多的基础设施组件

## 许可证

MIT License

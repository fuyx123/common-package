package nacos

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ExampleBasicUsage 基本使用示例
func ExampleBasicUsage() {
	// 简单获取配置
	config := NewNacos("application.yaml")
	fmt.Printf("获取到的配置: %s\n", config)
}

// ExampleAdvancedUsage 高级使用示例
func ExampleAdvancedUsage() {
	// 使用上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取客户端实例
	client, err := GetNacosClient("application.yaml")
	if err != nil {
		log.Fatalf("初始化客户端失败: %v", err)
	}
	defer client.Close()

	// 获取配置
	config, err := client.GetConfig(ctx, "my-config", "DEFAULT_GROUP")
	if err != nil {
		log.Printf("获取配置失败: %v", err)
		return
	}
	fmt.Printf("配置内容: %s\n", config)

	// 发布配置
	err = client.PublishConfig(ctx, "my-config", "DEFAULT_GROUP", "new config content")
	if err != nil {
		log.Printf("发布配置失败: %v", err)
		return
	}
	fmt.Println("配置发布成功")

	// 监听配置变化
	err = client.ListenConfig(ctx, "my-config", "DEFAULT_GROUP", func(newConfig string) {
		fmt.Printf("配置已更新: %s\n", newConfig)
	})
	if err != nil {
		log.Printf("监听配置失败: %v", err)
		return
	}
	fmt.Println("开始监听配置变化...")

	// 保持程序运行以监听配置变化
	time.Sleep(30 * time.Second)
}

// ExampleErrorHandling 错误处理示例
func ExampleErrorHandling() {
	client, err := GetNacosClient("application.yaml")
	if err != nil {
		if IsClientError(err) {
			log.Printf("客户端错误: %v", err)
		} else if IsConfigError(err) {
			log.Printf("配置错误: %v", err)
		} else if IsNetworkError(err) {
			log.Printf("网络错误: %v", err)
		} else {
			log.Printf("未知错误: %v", err)
		}
		return
	}

	ctx := context.Background()
	config, err := client.GetConfig(ctx, "non-existent-config", "DEFAULT_GROUP")
	if err != nil {
		if IsNetworkError(err) {
			log.Printf("网络连接问题，请检查Nacos服务器: %v", err)
		} else {
			log.Printf("获取配置失败: %v", err)
		}
		return
	}

	fmt.Printf("配置内容: %s\n", config)
}

// ExampleConfigValidation 配置验证示例
func ExampleConfigValidation() {
	// 加载配置
	config, err := LoadConfig("application.yaml")
	if err != nil {
		log.Printf("加载配置失败: %v", err)
		return
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		log.Printf("配置验证失败: %v", err)
		return
	}

	// 检查配置是否有效
	if !config.IsValid() {
		log.Println("配置无效")
		return
	}

	// 获取服务器URL
	serverURL := config.GetServerURL()
	fmt.Printf("Nacos服务器URL: %s\n", serverURL)

	fmt.Println("配置验证通过")
}

// ExampleConcurrentUsage 并发使用示例
func ExampleConcurrentUsage() {
	client, err := GetNacosClient("application.yaml")
	if err != nil {
		log.Fatalf("初始化客户端失败: %v", err)
	}
	defer client.Close()

	// 并发获取多个配置
	configs := []string{"config1", "config2", "config3"}
	results := make(chan string, len(configs))

	for _, configName := range configs {
		go func(name string) {
			ctx := context.Background()
			config, err := client.GetConfig(ctx, name, "DEFAULT_GROUP")
			if err != nil {
				log.Printf("获取配置 %s 失败: %v", name, err)
				results <- ""
				return
			}
			results <- fmt.Sprintf("%s: %s", name, config)
		}(configName)
	}

	// 收集结果
	for i := 0; i < len(configs); i++ {
		result := <-results
		if result != "" {
			fmt.Println(result)
		}
	}
}

// ExampleWithCustomConfig 使用自定义配置示例
func ExampleWithCustomConfig() {
	// 创建自定义配置
	customConfig := &Config{
		Nacos: NacosConfig{
			Namespace: "custom-namespace",
			Addr:      "localhost",
			Port:      8848,
			Dataid:    "custom-config",
			Group:     "CUSTOM_GROUP",
			TimeoutMs: 10000,
			LogLevel:  "debug",
		},
	}

	// 验证配置
	if err := customConfig.Validate(); err != nil {
		log.Printf("自定义配置验证失败: %v", err)
		return
	}

	fmt.Printf("自定义配置验证通过: %+v\n", customConfig)
}

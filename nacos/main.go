package nacos

import (
	"context"
	"fmt"
	"log"
)

// NewNacos 创建Nacos客户端并获取配置（向后兼容的简单接口）
func NewNacos(configPath string) string {
	client, err := InitNacos(configPath)
	if err != nil {
		log.Printf("初始化Nacos客户端失败: %v", err)
		return ""
	}

	ctx := context.Background()
	data, err := client.GetConfig(ctx, conf.Nacos.Dataid, conf.Nacos.Group)
	if err != nil {
		log.Printf("获取配置失败: %v", err)
		return ""
	}

	return data
}

// NewNacosWithContext 创建Nacos客户端并获取配置（支持上下文）
func NewNacosWithContext(ctx context.Context, configPath string) (string, error) {
	client, err := InitNacos(configPath)
	if err != nil {
		return "", fmt.Errorf("初始化Nacos客户端失败: %w", err)
	}

	data, err := client.GetConfig(ctx, "", "")
	if err != nil {
		return "", fmt.Errorf("获取配置失败: %w", err)
	}

	return data, nil
}

// GetNacosClient 获取Nacos客户端实例
func GetNacosClient(configPath string) (*NacosClient, error) {
	return InitNacos(configPath)
}

// GetConfig 获取配置的便捷方法
func GetConfig(configPath, dataId, group string) (string, error) {
	client, err := InitNacos(configPath)
	if err != nil {
		return "", fmt.Errorf("初始化Nacos客户端失败: %w", err)
	}

	ctx := context.Background()
	return client.GetConfig(ctx, dataId, group)
}

// PublishConfig 发布配置的便捷方法
func PublishConfig(configPath, dataId, group, content string) error {
	client, err := InitNacos(configPath)
	if err != nil {
		return fmt.Errorf("初始化Nacos客户端失败: %w", err)
	}

	ctx := context.Background()
	return client.PublishConfig(ctx, dataId, group, content)
}

// DeleteConfig 删除配置的便捷方法
func DeleteConfig(configPath, dataId, group string) error {
	client, err := InitNacos(configPath)
	if err != nil {
		return fmt.Errorf("初始化Nacos客户端失败: %w", err)
	}

	ctx := context.Background()
	return client.DeleteConfig(ctx, dataId, group)
}

// ListenConfig 监听配置变化的便捷方法
func ListenConfig(configPath, dataId, group string, callback func(string)) error {
	client, err := InitNacos(configPath)
	if err != nil {
		return fmt.Errorf("初始化Nacos客户端失败: %w", err)
	}

	ctx := context.Background()
	return client.ListenConfig(ctx, dataId, group, callback)
}

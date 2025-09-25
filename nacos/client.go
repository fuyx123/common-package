package nacos

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// NacosClient 封装了Nacos配置中心客户端
type NacosClient struct {
	client config_client.IConfigClient
	config *Config
	mu     sync.RWMutex
}

var (
	instance *NacosClient
	once     sync.Once
)

// InitNacos 初始化Nacos客户端（单例模式）
func InitNacos(configPath string) (*NacosClient, error) {
	var initErr error

	once.Do(func() {
		// 加载配置文件
		config, err := LoadConfig(configPath)
		conf = config
		if err != nil {
			initErr = fmt.Errorf("加载配置文件失败: %w", err)
			return
		}

		// 验证配置
		if err := config.Validate(); err != nil {
			initErr = fmt.Errorf("配置验证失败: %w", err)
			return
		}

		// 创建客户端配置
		clientConfig := constant.ClientConfig{
			NamespaceId:         config.Nacos.Namespace,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			LogLevel:            "info", // 改为info级别，减少日志输出
		}

		// 创建服务器配置
		serverConfigs := []constant.ServerConfig{
			{
				IpAddr:      config.Nacos.Addr,
				ContextPath: "/nacos",
				Port:        config.Nacos.Port,
				Scheme:      "http",
			},
		}

		// 创建Nacos客户端
		configClient, err := clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
		if err != nil {
			initErr = fmt.Errorf("创建Nacos客户端失败: %w", err)
			return
		}

		instance = &NacosClient{
			client: configClient,
			config: &config,
		}

		log.Printf("Nacos客户端初始化成功，服务器: %s:%d", config.Nacos.Addr, config.Nacos.Port)
	})

	return instance, initErr
}

// GetConfig 获取配置
func (c *NacosClient) GetConfig(ctx context.Context, dataId, group string) (string, error) {
	if c == nil || c.client == nil {
		return "", fmt.Errorf("Nacos客户端未初始化")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// 使用默认值如果参数为空
	if dataId == "" {
		dataId = c.config.Nacos.Dataid
	}
	if group == "" {
		group = c.config.Nacos.Group
	}

	config, err := c.client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return "", fmt.Errorf("获取配置失败 [DataId: %s, Group: %s]: %w", dataId, group, err)
	}

	return config, nil
}

// PublishConfig 发布配置
func (c *NacosClient) PublishConfig(ctx context.Context, dataId, group, content string) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("Nacos客户端未初始化")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// 使用默认值如果参数为空
	if dataId == "" {
		dataId = c.config.Nacos.Dataid
	}
	if group == "" {
		group = c.config.Nacos.Group
	}

	success, err := c.client.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	})
	if err != nil {
		return fmt.Errorf("发布配置失败 [DataId: %s, Group: %s]: %w", dataId, group, err)
	}

	if !success {
		return fmt.Errorf("发布配置失败，返回false")
	}

	return nil
}

// DeleteConfig 删除配置
func (c *NacosClient) DeleteConfig(ctx context.Context, dataId, group string) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("Nacos客户端未初始化")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// 使用默认值如果参数为空
	if dataId == "" {
		dataId = c.config.Nacos.Dataid
	}
	if group == "" {
		group = c.config.Nacos.Group
	}

	success, err := c.client.DeleteConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return fmt.Errorf("删除配置失败 [DataId: %s, Group: %s]: %w", dataId, group, err)
	}

	if !success {
		return fmt.Errorf("删除配置失败，返回false")
	}

	return nil
}

// ListenConfig 监听配置变化
func (c *NacosClient) ListenConfig(ctx context.Context, dataId, group string, callback func(string)) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("Nacos客户端未初始化")
	}

	// 使用默认值如果参数为空
	if dataId == "" {
		dataId = c.config.Nacos.Dataid
	}
	if group == "" {
		group = c.config.Nacos.Group
	}

	err := c.client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			if callback != nil {
				callback(data)
			}
		},
	})
	if err != nil {
		return fmt.Errorf("监听配置失败 [DataId: %s, Group: %s]: %w", dataId, group, err)
	}

	return nil
}

// Close 关闭客户端
func (c *NacosClient) Close() error {
	if c == nil || c.client == nil {
		return nil
	}

	// Nacos SDK 没有提供显式的关闭方法
	// 这里可以添加清理逻辑
	log.Println("Nacos客户端已关闭")
	return nil
}

// GetClient 获取原始客户端（用于高级用法）
func (c *NacosClient) GetClient() config_client.IConfigClient {
	if c == nil {
		return nil
	}
	return c.client
}

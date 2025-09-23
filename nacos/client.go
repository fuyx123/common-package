package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"sync"
)

type NacosClient struct {
	client config_client.IConfigClient
}

var (
	instance *NacosClient
	once     sync.Once
)

// nacos-配置中心 / 服务发现和注册
// 设计模式   单例模式     工厂模式
func InitNacos(configPath string) config_client.IConfigClient {
	// 调用LoadConfig加载配置文件
	nacos, err := LoadConfig(configPath)
	if err != nil {
		panic("nacos配置解析失败")
	}
	once.Do(func() {
		clientConfig := constant.ClientConfig{
			NamespaceId:         Conf.Nacos.Namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			LogLevel:            "debug",
		}
		// 创建nacos的连接（先通过nacos官方提供的结构体进行赋值）
		serverConfigs := []constant.ServerConfig{
			{
				IpAddr:      Conf.Nacos.Addr,
				ContextPath: "/nacos",
				Port:        nacos.Nacos.Port,
				Scheme:      "http",
			},
		}
		// 通过New来实例化nacos
		configClient, err := clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
		if err != nil {
			fmt.Println(err)
		}
		instance.client = configClient

	})
	return instance.client // interface的接口
}

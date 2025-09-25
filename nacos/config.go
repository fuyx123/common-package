package nacos

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/viper"
)

var conf Config

// Config Nacos配置结构
type Config struct {
	Nacos NacosConfig `mapstructure:"nacos"`
}

// NacosConfig Nacos具体配置
type NacosConfig struct {
	Namespace string `mapstructure:"namespace"`
	Addr      string `mapstructure:"addr"`
	Port      uint64 `mapstructure:"port"`
	Dataid    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
	// 新增配置项
	TimeoutMs    int64  `mapstructure:"timeout_ms"`
	LogLevel     string `mapstructure:"log_level"`
	LogDir       string `mapstructure:"log_dir"`
	CacheDir     string `mapstructure:"cache_dir"`
	NotLoadCache bool   `mapstructure:"not_load_cache"`
	Scheme       string `mapstructure:"scheme"`
	ContextPath  string `mapstructure:"context_path"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Nacos: NacosConfig{
			TimeoutMs:    5000,
			LogLevel:     "info",
			LogDir:       "/tmp/nacos/log",
			CacheDir:     "/tmp/nacos/cache",
			NotLoadCache: true,
			Scheme:       "http",
			ContextPath:  "/nacos",
			Group:        "DEFAULT_GROUP",
		},
	}
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (Config, error) {
	viper.SetConfigFile(configPath)

	// 设置默认值
	viper.SetDefault("nacos.timeout_ms", 5000)
	viper.SetDefault("nacos.log_level", "info")
	viper.SetDefault("nacos.log_dir", "/tmp/nacos/log")
	viper.SetDefault("nacos.cache_dir", "/tmp/nacos/cache")
	viper.SetDefault("nacos.not_load_cache", true)
	viper.SetDefault("nacos.scheme", "http")
	viper.SetDefault("nacos.context_path", "/nacos")
	viper.SetDefault("nacos.group", "DEFAULT_GROUP")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Nacos.Addr == "" {
		return fmt.Errorf("nacos地址不能为空")
	}

	// 验证IP地址格式
	if net.ParseIP(c.Nacos.Addr) == nil {
		// 如果不是IP，尝试解析域名
		if _, err := net.LookupHost(c.Nacos.Addr); err != nil {
			return fmt.Errorf("无效的nacos地址: %s", c.Nacos.Addr)
		}
	}

	if c.Nacos.Port == 0 {
		return fmt.Errorf("nacos端口不能为0")
	}

	if c.Nacos.Port > 65535 {
		return fmt.Errorf("nacos端口不能超过65535")
	}

	if c.Nacos.Dataid == "" {
		return fmt.Errorf("dataid不能为空")
	}

	if c.Nacos.Group == "" {
		return fmt.Errorf("group不能为空")
	}

	// 验证日志级别
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if c.Nacos.LogLevel != "" {
		found := false
		for _, level := range validLogLevels {
			if strings.ToLower(c.Nacos.LogLevel) == level {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("无效的日志级别: %s，支持: %v", c.Nacos.LogLevel, validLogLevels)
		}
	}

	// 验证协议
	if c.Nacos.Scheme != "" {
		validSchemes := []string{"http", "https"}
		found := false
		for _, scheme := range validSchemes {
			if strings.ToLower(c.Nacos.Scheme) == scheme {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("无效的协议: %s，支持: %v", c.Nacos.Scheme, validSchemes)
		}
	}

	return nil
}

// GetServerURL 获取服务器URL
func (c *Config) GetServerURL() string {
	scheme := c.Nacos.Scheme
	if scheme == "" {
		scheme = "http"
	}

	contextPath := c.Nacos.ContextPath
	if contextPath == "" {
		contextPath = "/nacos"
	}

	return fmt.Sprintf("%s://%s:%d%s", scheme, c.Nacos.Addr, c.Nacos.Port, contextPath)
}

// IsValid 检查配置是否有效
func (c *Config) IsValid() bool {
	return c.Validate() == nil
}

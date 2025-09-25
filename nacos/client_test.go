package nacos

import (
	"context"
	"testing"
	"time"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Nacos: NacosConfig{
					Addr:   "localhost",
					Port:   8848,
					Dataid: "test-config",
					Group:  "DEFAULT_GROUP",
				},
			},
			wantErr: false,
		},
		{
			name: "empty address",
			config: Config{
				Nacos: NacosConfig{
					Addr:   "",
					Port:   8848,
					Dataid: "test-config",
					Group:  "DEFAULT_GROUP",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			config: Config{
				Nacos: NacosConfig{
					Addr:   "localhost",
					Port:   0,
					Dataid: "test-config",
					Group:  "DEFAULT_GROUP",
				},
			},
			wantErr: true,
		},
		{
			name: "empty dataid",
			config: Config{
				Nacos: NacosConfig{
					Addr:   "localhost",
					Port:   8848,
					Dataid: "",
					Group:  "DEFAULT_GROUP",
				},
			},
			wantErr: true,
		},
		{
			name: "empty group",
			config: Config{
				Nacos: NacosConfig{
					Addr:   "localhost",
					Port:   8848,
					Dataid: "test-config",
					Group:  "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: Config{
				Nacos: NacosConfig{
					Addr:     "localhost",
					Port:     8848,
					Dataid:   "test-config",
					Group:    "DEFAULT_GROUP",
					LogLevel: "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid scheme",
			config: Config{
				Nacos: NacosConfig{
					Addr:   "localhost",
					Port:   8848,
					Dataid: "test-config",
					Group:  "DEFAULT_GROUP",
					Scheme: "ftp",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// 验证默认值
	if config.Nacos.TimeoutMs != 5000 {
		t.Errorf("Expected TimeoutMs = 5000, got %d", config.Nacos.TimeoutMs)
	}

	if config.Nacos.LogLevel != "info" {
		t.Errorf("Expected LogLevel = 'info', got %s", config.Nacos.LogLevel)
	}

	if config.Nacos.Scheme != "http" {
		t.Errorf("Expected Scheme = 'http', got %s", config.Nacos.Scheme)
	}

	if config.Nacos.Group != "DEFAULT_GROUP" {
		t.Errorf("Expected Group = 'DEFAULT_GROUP', got %s", config.Nacos.Group)
	}
}

func TestGetServerURL(t *testing.T) {
	config := &Config{
		Nacos: NacosConfig{
			Addr:        "localhost",
			Port:        8848,
			Scheme:      "http",
			ContextPath: "/nacos",
		},
	}

	expected := "http://localhost:8848/nacos"
	actual := config.GetServerURL()

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestIsValid(t *testing.T) {
	validConfig := &Config{
		Nacos: NacosConfig{
			Addr:   "localhost",
			Port:   8848,
			Dataid: "test-config",
			Group:  "DEFAULT_GROUP",
		},
	}

	if !validConfig.IsValid() {
		t.Error("Expected valid config to be valid")
	}

	invalidConfig := &Config{
		Nacos: NacosConfig{
			Addr:   "",
			Port:   8848,
			Dataid: "test-config",
			Group:  "DEFAULT_GROUP",
		},
	}

	if invalidConfig.IsValid() {
		t.Error("Expected invalid config to be invalid")
	}
}

func TestNacosError(t *testing.T) {
	// 测试错误创建
	err := NewNacosError("TEST_ERROR", "测试错误", nil)
	if err.Error() != "[TEST_ERROR] 测试错误" {
		t.Errorf("Expected '[TEST_ERROR] 测试错误', got '%s'", err.Error())
	}

	// 测试错误包装
	wrappedErr := WrapError(err, "包装错误")
	if wrappedErr.Error() != "[TEST_ERROR] 包装错误: [TEST_ERROR] 测试错误" {
		t.Errorf("Unexpected wrapped error: %s", wrappedErr.Error())
	}

	// 测试错误类型检查
	if !IsConfigError(ErrConfigNotFound) {
		t.Error("Expected ErrConfigNotFound to be a config error")
	}

	if !IsClientError(ErrClientNotInit) {
		t.Error("Expected ErrClientNotInit to be a client error")
	}

	if IsConfigError(ErrClientNotInit) {
		t.Error("Expected ErrClientNotInit not to be a config error")
	}
}

func TestContextTimeout(t *testing.T) {
	// 测试上下文超时
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// 等待超时
	time.Sleep(10 * time.Millisecond)

	select {
	case <-ctx.Done():
		// 上下文已取消，这是预期的
	default:
		t.Error("Expected context to be cancelled")
	}
}

// 基准测试
func BenchmarkConfigValidation(b *testing.B) {
	config := &Config{
		Nacos: NacosConfig{
			Addr:   "localhost",
			Port:   8848,
			Dataid: "test-config",
			Group:  "DEFAULT_GROUP",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Validate()
	}
}

func BenchmarkGetServerURL(b *testing.B) {
	config := &Config{
		Nacos: NacosConfig{
			Addr:        "localhost",
			Port:        8848,
			Scheme:      "http",
			ContextPath: "/nacos",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.GetServerURL()
	}
}

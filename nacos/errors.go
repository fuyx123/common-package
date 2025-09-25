package nacos

import (
	"fmt"
	"net"
)

// NacosError 自定义Nacos错误类型
type NacosError struct {
	Code    string
	Message string
	Err     error
}

func (e *NacosError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *NacosError) Unwrap() error {
	return e.Err
}

// 预定义的错误类型
var (
	// 配置相关错误
	ErrConfigNotFound       = &NacosError{Code: "CONFIG_NOT_FOUND", Message: "配置未找到"}
	ErrConfigInvalid        = &NacosError{Code: "CONFIG_INVALID", Message: "配置无效"}
	ErrConfigLoadFailed     = &NacosError{Code: "CONFIG_LOAD_FAILED", Message: "配置加载失败"}
	ErrConfigValidateFailed = &NacosError{Code: "CONFIG_VALIDATE_FAILED", Message: "配置验证失败"}

	// 客户端相关错误
	ErrClientNotInit    = &NacosError{Code: "CLIENT_NOT_INIT", Message: "客户端未初始化"}
	ErrClientInitFailed = &NacosError{Code: "CLIENT_INIT_FAILED", Message: "客户端初始化失败"}
	ErrClientConnection = &NacosError{Code: "CLIENT_CONNECTION", Message: "客户端连接失败"}

	// 网络相关错误
	ErrNetworkTimeout     = &NacosError{Code: "NETWORK_TIMEOUT", Message: "网络超时"}
	ErrNetworkUnreachable = &NacosError{Code: "NETWORK_UNREACHABLE", Message: "网络不可达"}
	ErrServerUnavailable  = &NacosError{Code: "SERVER_UNAVAILABLE", Message: "服务器不可用"}

	// 操作相关错误
	ErrOperationFailed = &NacosError{Code: "OPERATION_FAILED", Message: "操作失败"}
	ErrPublishFailed   = &NacosError{Code: "PUBLISH_FAILED", Message: "发布配置失败"}
	ErrDeleteFailed    = &NacosError{Code: "DELETE_FAILED", Message: "删除配置失败"}
	ErrListenFailed    = &NacosError{Code: "LISTEN_FAILED", Message: "监听配置失败"}
)

// NewNacosError 创建新的Nacos错误
func NewNacosError(code, message string, err error) *NacosError {
	return &NacosError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IsNetworkError 检查是否为网络错误
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否为网络超时错误
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout() || netErr.Temporary()
	}

	// 检查自定义错误
	if nacosErr, ok := err.(*NacosError); ok {
		switch nacosErr.Code {
		case "NETWORK_TIMEOUT", "NETWORK_UNREACHABLE", "SERVER_UNAVAILABLE", "CLIENT_CONNECTION":
			return true
		}
	}

	return false
}

// IsConfigError 检查是否为配置错误
func IsConfigError(err error) bool {
	if err == nil {
		return false
	}

	if nacosErr, ok := err.(*NacosError); ok {
		switch nacosErr.Code {
		case "CONFIG_NOT_FOUND", "CONFIG_INVALID", "CONFIG_LOAD_FAILED", "CONFIG_VALIDATE_FAILED":
			return true
		}
	}

	return false
}

// IsClientError 检查是否为客户端错误
func IsClientError(err error) bool {
	if err == nil {
		return false
	}

	if nacosErr, ok := err.(*NacosError); ok {
		switch nacosErr.Code {
		case "CLIENT_NOT_INIT", "CLIENT_INIT_FAILED", "CLIENT_CONNECTION":
			return true
		}
	}

	return false
}

// WrapError 包装错误
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

	if nacosErr, ok := err.(*NacosError); ok {
		return &NacosError{
			Code:    nacosErr.Code,
			Message: message,
			Err:     nacosErr,
		}
	}

	return &NacosError{
		Code:    "UNKNOWN",
		Message: message,
		Err:     err,
	}
}

package nacos

import "github.com/spf13/viper"

type Config struct {
	Nacos struct {
		Namespace string `mapstructure:"namespace"`
		Addr      string `mapstructure:"addr"`
		Port      uint64 `mapstructure:"port"`
		Dataid    string `mapstructure:"dataid"`
		Group     string `mapstructure:"group"`
	} `mapstructure:"nacos"`
}

var Conf Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (Config, error) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	//解析
	err := viper.Unmarshal(&Conf)
	return Conf, err
}

package nacos

import "github.com/nacos-group/nacos-sdk-go/v2/vo"

func NewNacos(configPath string) string {
	client := InitNacos(configPath)
	data, _ := client.GetConfig(vo.ConfigParam{
		DataId: Conf.Nacos.Dataid,
		Group:  Conf.Nacos.Group,
	})
	return data
}

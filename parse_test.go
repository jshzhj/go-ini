package go_ini

import (
	"testing"
)

//节点
type Config struct {
	ServerConf ServerConfig `ini:"server"`
	Mysql      MysqlConfig
}

//选项
type ServerConfig struct {
	Ipaa string `default:"192.168.1.1" ini:"ip"`
	Port uint   `default:"8080"`
}

//选项
type MysqlConfig struct {
	Username string
	Password string
}

func TestConfigToStruct(t *testing.T) {
	iniPath := "config.ini" //配置文件路径
	conf := &Config{}

	err := Parse(iniPath, conf)
	if err != nil {
		t.Error(err)
	}

	t.Logf("ini到结构体解析成功! %#v", conf)
}

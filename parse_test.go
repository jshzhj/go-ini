package go_ini

import (
	"testing"
)

//节点
type Config struct {
	ServerConf ServerConfig `ini:"server"`
	MysqlConf  MysqlConfig  `ini:"mysql"`
}

//选项
type ServerConfig struct {
	Ip   string `ini:"ip"`
	Port uint   `ini:"port"`
}

//选项
type MysqlConfig struct {
	Username string `ini:"username"`
	Password string `ini:"password"`
}

func TestConfigToStruct(t *testing.T) {
	iniPath := "config.ini"

	var conf Config

	err := Parse(iniPath, &conf)
	if err != nil {
		t.Error(err)
	}

	t.Logf("ini到结构体解析成功! %#v", conf)
}


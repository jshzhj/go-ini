package go_ini

import (
	"testing"
)

//节点
// `ini:"server"`为配置文件节点到结构体字段的映射,可写可不写,不写默认为结构体字段名小写
type Config struct {
	ServerConf ServerConfig `ini:"server"` //有映射(tag),则结构体字段名随意写
	Mysql      MysqlConfig  //无映射(tag),则结构体字段名要与配置文件节点名对应
}

//选项
type ServerConfig struct {
	Ip   string `default:"192.168.1.1" ini:"ip"` //default:默认值(可有可无) ini:映射(同上,可有可无)
	Port uint   `default:"8080"`
}

//选项
type MysqlConfig struct {
	Username string
	Password string
}

func TestConfigToStruct(t *testing.T) {
	iniPath := "config.ini" //配置文件路径
	conf := &Config{}       //实力化结构体

	err := Parse(iniPath, conf)
	if err != nil {
		t.Error(err)
	}

	t.Logf("ini到结构体解析成功! %#v", conf)
}

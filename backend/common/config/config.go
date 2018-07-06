package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type RedisConf struct {
	Host     string
	Port     string
	Password string
	Db       int
}

type DBConf struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	Target   string
	MaxIdle  int
	MaxOpen  int
	Showsql  bool
	Location string
}

type System struct {
	Name   string
	Prefix string
	Cors   []string
}

type Conf struct {
	App struct {
		Mode string
		Log  string
	}
	Quant struct {
		System   System
		Redis    RedisConf
		Database DBConf
		Rpc      struct {
			Host string
			Port string
		}
	}
	Time struct {
		System   System
		Redis    RedisConf
		Database DBConf
	}
	User struct {
		System   System
		Redis    RedisConf
		Database DBConf
	}
}

var onceConfig *Conf = &Conf{}

func Get() *Conf {
	if "" == onceConfig.App.Mode {
		configPath := os.Getenv("ConfigPath")
		if configPath == "" {
			configPath = "../config/.config.yaml"
		}
		yamlFile, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}
		config := &Conf{}
		err = yaml.Unmarshal(yamlFile, config)
		if err != nil {
			panic(err)
		}
		onceConfig = config
	}
	return onceConfig
}

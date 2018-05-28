package common

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type RedisConf struct {
	Name     string
	Host     string
	Port     string
	Password string
	Db       int
}

type DBConf struct {
	Name     string
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

type Conf struct {
	App struct {
		Mode string
	}
	Quant struct {
		Redis    RedisConf
		Database DBConf
	}
}

var onceConfig *Conf = &Conf{}

func GetConfig() *Conf {
	if (Conf{}) == *onceConfig {
		configPath := os.Getenv("ConfigPath")
		if configPath == "" {
			configPath = "./config/.config.yaml"
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

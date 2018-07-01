package main

import (
	"evolution/backend/common/config"
	"evolution/backend/quant/models"
	"fmt"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func main() {
	yamlFile, err := ioutil.ReadFile("../../config/.config.yaml")
	if err != nil {
		panic(err)
	}
	conf := &config.Conf{}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		panic(err)
	}
	exec(conf.Quant.Database, conf.App.Mode)
}

func exec(dbConfig config.DBConf, mode string) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Target)

	engine, err := xorm.NewEngine(dbConfig.Type, source)
	if err != nil {
		panic(err)
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	engine.TZLocation = location
	engine.StoreEngine("InnoDB")
	engine.Charset("utf8")
	if mode == "debug" {
		err = engine.DropTables(new(models.Signal), new(models.Pool), new(models.MapPoolItem), new(models.Classify), new(models.Item))
		if err != nil {
			panic(err)
		}
	}
	err = engine.Sync2(new(models.Signal), new(models.Pool), new(models.Classify), new(models.Item), new(models.MapPoolItem), new(models.MapClassifyItem))
	if err != nil {
		panic(err)
	}
}

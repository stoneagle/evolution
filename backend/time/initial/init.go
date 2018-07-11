package main

import (
	"evolution/backend/common/config"
	"evolution/backend/time/models"
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	yaml "gopkg.in/yaml.v2"
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
	exec(conf.Time.Database, conf.App.Mode)
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
	// if mode == "debug" {
	// 	err = engine.DropTables(new(models.Area), new(models.Country), new(models.Field), new(models.Entity), new(models.Phase), new(models.Resource), new(models.Quest), new(models.QuestTeam), new(models.QuestTimeTable), new(models.QuestTarget), new(models.QuestEntity), new(models.Project))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	err = engine.Sync2(new(models.Area), new(models.Country), new(models.Field), new(models.Entity), new(models.Phase), new(models.Resource), new(models.Quest), new(models.QuestTeam), new(models.QuestTimeTable), new(models.QuestTarget), new(models.QuestEntity), new(models.Project))
	if err != nil {
		panic(err)
	}
}

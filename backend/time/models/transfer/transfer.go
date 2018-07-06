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

type SrcConf struct {
	Database config.DBConf
}

func main() {
	// read dest config
	yamlFile, err := ioutil.ReadFile("../../../config/.config.yaml")
	if err != nil {
		panic(err)
	}
	conf := &config.Conf{}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		panic(err)
	}
	dbConfig := conf.Time.Database
	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Target)
	destEng, err := xorm.NewEngine(dbConfig.Type, source)
	if err != nil {
		panic(err)
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	destEng.TZLocation = location
	destEng.StoreEngine("InnoDB")
	destEng.Charset("utf8")

	// read src config
	srcYamlFile, err := ioutil.ReadFile("./.config.yaml")
	if err != nil {
		panic(err)
	}
	srcConf := &SrcConf{}
	err = yaml.Unmarshal(srcYamlFile, srcConf)
	if err != nil {
		panic(err)
	}
	srcSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", srcConf.Database.User, srcConf.Database.Password, srcConf.Database.Host, srcConf.Database.Port, srcConf.Database.Target)
	srcEng, err := xorm.NewEngine(srcConf.Database.Type, srcSource)
	if err != nil {
		panic(err)
	}
	location, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	srcEng.TZLocation = location
	srcEng.StoreEngine("InnoDB")
	srcEng.Charset("utf8")
	// srcEng.ShowSQL(true)

	// new(php.Area).Transfer(srcEng, destEng)
	// new(php.Country).Transfer(srcEng, destEng)
	// initField(destEng)
	initPhase(destEng)
	// new(php.EntityAsset).Transfer(srcEng, destEng)
	// new(php.EntityCircle).Transfer(srcEng, destEng)
	// new(php.EntityLife).Transfer(srcEng, destEng)
	// new(php.EntityQuest).Transfer(srcEng, destEng)
	// new(php.EntityWork).Transfer(srcEng, destEng)
	// new(php.EntitySkill).Transfer(srcEng, destEng)
}

func initField(des *xorm.Engine) {
	initMap := map[int]string{
		1: "技能",
		2: "资产",
		3: "作品",
		4: "圈子",
		5: "经历",
		6: "日常",
	}
	news := make([]models.Field, 0)
	for k, v := range initMap {
		tmp := models.Field{}
		tmp.Name = v
		tmp.Id = k
		news = append(news, tmp)
	}
	affected, err := des.Insert(&news)
	if err != nil {
		fmt.Printf("field transfer error:%v\r\n", err.Error())
	} else {
		fmt.Printf("field transfer success:%v\r\n", affected)
	}
}

func initPhase(des *xorm.Engine) {
	initMap := map[int][]string{
		1: []string{"探索", "学习", "原理", "突破"},
		2: []string{"探索", "实践", "行业", "发展"},
		3: []string{"探索", "欣赏", "思想", "创作"},
		4: []string{"探索", "融入", "影响", "领导"},
		5: []string{"探索", "尝试", "挑战", "记录"},
		6: []string{"杂务", "规划"},
	}
	news := make([]models.Phase, 0)
	for k, slice := range initMap {
		for l, v := range slice {
			tmp := models.Phase{}
			tmp.Name = v
			tmp.Desc = ""
			tmp.Level = l + 1
			tmp.FieldId = k
			news = append(news, tmp)
		}
	}
	affected, err := des.Insert(&news)
	if err != nil {
		fmt.Printf("phase transfer error:%v\r\n", err.Error())
	} else {
		fmt.Printf("phase transfer success:%v\r\n", affected)
	}
}

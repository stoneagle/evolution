package main

import (
	"evolution/backend/common/config"
	"evolution/backend/time/models"
	"evolution/backend/time/models/php"
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	yaml "gopkg.in/yaml.v2"
)

var (
	FieldMap map[int]string = map[int]string{
		1: "技能",
		2: "资产",
		3: "作品",
		4: "圈子",
		5: "经历",
		6: "日常",
	}
	FieldColorMap map[int]string = map[int]string{
		1: "#57b94c",
		2: "#5187c6",
		3: "#edba3c",
		4: "#FF9900",
		5: "#ee4e75",
		6: "#B0B0B0",
	}
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

	userId := 1
	// new(php.Area).Transfer(srcEng, destEng)
	// initAreaType(destEng)
	// new(php.Country).Transfer(srcEng, destEng)
	// initField(destEng)
	// initPhase(destEng)
	// new(php.EntityAsset).Transfer(srcEng, destEng)
	// new(php.EntityCircle).Transfer(srcEng, destEng)
	// new(php.EntityQuest).Transfer(srcEng, destEng)
	// new(php.EntityLife).Transfer(srcEng, destEng)
	// new(php.EntityWork).Transfer(srcEng, destEng)
	// new(php.EntitySkill).Transfer(srcEng, destEng)
	// new(php.TargetEntityLink).Transfer(srcEng, destEng, userId)
	// new(php.Target).Transfer(srcEng, destEng, userId)
	new(php.Project).Transfer(srcEng, destEng, userId)
	// initUserResourceTime(destEng, userId)
}

func initUserResourceTime(des *xorm.Engine, userId int) {
	actionsJoin := make([]models.ActionJoin, 0)
	sql := des.Unscoped().Table("action").Join("INNER", "task", "task.id = action.task_id")
	err := sql.Find(&actionsJoin)
	if err != nil {
		fmt.Printf("action get error:%v\r\n", err.Error())
		return
	}

	actions := make([]models.Action, 0)
	for _, one := range actionsJoin {
		one.Action.Task = one.Task
		actions = append(actions, one.Action)
	}
	userResourceTime := map[int]int{}
	for _, one := range actions {
		resourceTime, ok := userResourceTime[one.Task.ResourceId]
		if !ok {
			userResourceTime[one.Task.ResourceId] = 0
		}
		resourceTime += one.Time
		userResourceTime[one.Task.ResourceId] = resourceTime
	}
	updateNum := 0
	for resourceId, sumTime := range userResourceTime {
		userResource := models.UserResource{}
		userResource.Time = sumTime
		_, err = des.Where("resource_id = ?", resourceId).And("user_id = ?", userId).Update(&userResource)
		if err != nil {
			fmt.Printf("user resource time update error:%v\r\n", err.Error())
			return
		}
		updateNum++
	}
	fmt.Printf("user resource time init success:%v\r\n", updateNum)
	return
}

func initField(des *xorm.Engine) {
	news := make([]models.Field, 0)
	for k, v := range FieldMap {
		tmp := models.Field{}
		tmp.Name = v
		tmp.Id = k
		color, ok := FieldColorMap[k]
		if !ok {
			fmt.Printf("field color not exist")
			return
		}
		tmp.Color = color
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
	thresholdMap := map[int]int{
		1: 80,
		2: 400,
		3: 2000,
		4: 10000,
	}
	news := make([]models.Phase, 0)
	for k, slice := range initMap {
		for l, v := range slice {
			tmp := models.Phase{}
			tmp.Name = v
			tmp.Desc = ""
			tmp.Level = l + 1
			tmp.Threshold = thresholdMap[tmp.Level]
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

func initAreaType(des *xorm.Engine) {
	var err error
	news := make([]models.Area, 0)
	des.Find(&news)
	for _, one := range news {
		tmp := new(models.Area)
		if one.Parent == 0 {
			tmp.Type = models.AreaTypeRoot
			_, err = des.Id(one.Id).Update(tmp)
		} else {
			asParent := make([]models.Area, 0)
			des.Where("parent = ?", one.Id).Find(&asParent)
			if len(asParent) > 0 {
				tmp.Type = models.AreaTypeNode
				_, err = des.Id(one.Id).Update(tmp)
			} else {
				tmp.Type = models.AreaTypeLeaf
				_, err = des.Id(one.Id).Update(tmp)
			}
		}
		if err != nil {
			fmt.Printf("%v area type update error:%v\r\n", one.Name, err.Error())
			return
		}
	}
	fmt.Printf("area type init success\r\n")
}

package database

import (
	"evolution/backend/common/config"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

var mapXormEngine map[string]*xorm.Engine = make(map[string]*xorm.Engine)

func SetXorm(key string, e *xorm.Engine) {
	mapXormEngine[key] = e
}

func GetXorm(key string) *xorm.Engine {
	e, ok := mapXormEngine[key]
	if !ok {
		panic("engine not exist")
	}
	return e
}

func SetProjectXorm(dbConfig config.DBConf, name string) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&interpolateParams=true&parseTime=true&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Target)

	engine, err := xorm.NewEngine(dbConfig.Type, source)
	if err != nil {
		panic(err)
	}
	engine.SetMaxIdleConns(dbConfig.MaxIdle)
	engine.SetMaxOpenConns(dbConfig.MaxOpen)

	location, err := time.LoadLocation(dbConfig.Location)
	if err != nil {
		panic(err)
	}
	engine.DatabaseTZ = location
	engine.TZLocation = location
	if dbConfig.Showsql {
		engine.ShowSQL(true)
	}

	SetXorm(name, engine)
}

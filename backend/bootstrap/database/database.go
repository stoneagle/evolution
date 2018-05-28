package database

import (
	"fmt"
	"quant/backend/bootstrap"
	"quant/backend/common"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func Configure(b *bootstrap.Bootstrapper) {
	quantDB := b.Config.Quant.Database
	setProjectEngine(quantDB)

	redisConf := b.Config.Quant.Redis
	common.SetRedis(redisConf.Host, redisConf.Port, redisConf.Password, redisConf.Db)
}

func setProjectEngine(dbConfig common.DBConf) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Target)

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
	engine.TZLocation = location
	if dbConfig.Showsql {
		engine.ShowSQL(true)
	}

	common.SetEngine(dbConfig.Name, engine)
}

package database

import (
	"evolution/backend/common/config"
	"evolution/backend/common/database"
	"evolution/backend/quant/bootstrap"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func Configure(b *bootstrap.Bootstrapper) {
	quantDB := b.Config.Quant.Database
	setProjectXorm(quantDB)

	redisConf := b.Config.Quant.Redis
	database.SetRedis(redisConf.Host, redisConf.Port, redisConf.Password, redisConf.Db)
}

func setProjectXorm(dbConfig config.DBConf) {
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

	database.SetXorm(dbConfig.Name, engine)
}

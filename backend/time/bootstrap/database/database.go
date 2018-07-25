package database

import (
	"evolution/backend/common/database"
	"evolution/backend/time/bootstrap"

	_ "github.com/go-sql-driver/mysql"
)

func Configure(b *bootstrap.Bootstrapper) {
	database.SetProjectXorm(b.Config.Time.Database, b.Config.Time.System.Name)
	database.SetRedis(b.Config.Time.Redis)
}

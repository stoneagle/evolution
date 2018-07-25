package database

import (
	"evolution/backend/common/database"
	"evolution/backend/system/bootstrap"

	_ "github.com/go-sql-driver/mysql"
)

func Configure(b *bootstrap.Bootstrapper) {
	database.SetProjectXorm(b.Config.System.Database, b.Config.System.System.Name)
	database.SetRedis(b.Config.System.Redis)
}

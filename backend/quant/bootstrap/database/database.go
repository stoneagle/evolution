package database

import (
	"evolution/backend/common/database"
	"evolution/backend/quant/bootstrap"

	_ "github.com/go-sql-driver/mysql"
)

func Configure(b *bootstrap.Bootstrapper) {
	database.SetProjectXorm(b.Config.Quant.Database, b.Config.Quant.System.Name)
	database.SetRedis(b.Config.Quant.Redis)
}

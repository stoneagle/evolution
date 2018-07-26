package main

import (
	"evolution/backend/time/bootstrap"
	"evolution/backend/time/bootstrap/database"
	"evolution/backend/time/bootstrap/route"
	"time"
)

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("time", "wuzhongyang@wzy.com")
	app.Bootstrap()

	app.Configure(database.Configure, route.Configure)
	_, err := time.LoadLocation(app.Config.Time.System.Location)
	if err != nil {
		panic(err)
	}
	return app
}

func main() {
	app := newApp()
	app.Listen(":" + app.Config.Time.System.Port)
}

package main

import (
	"github.com/risasim/projectM5/project/src/server/app"
)

func main() {
	var app app.App
	app.CreateConnection()
	app.Migrate()
	app.CreateRoutes()
	app.Run()
}

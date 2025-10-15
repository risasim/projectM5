package main

import "github.com/risasim/projectM5/project/src/server/state"

func main() {
	var app App
	gameManager := state.NewGameManager()
	app.CreateConnection()
	app.Migrate()
	app.CreateRoutes()
	app.Run(gameManager)
}

package main

import "github.com/risasim/projectM5/project/src/server/state"

func main() {
	app := &App{}
	app.CreateConnection()
	app.Migrate()
	app.SetupLogin()
	app.CreateRoutes()
	gameManager := state.NewGameManager()
	app.Run(gameManager)
}

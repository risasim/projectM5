package main

import "github.com/risasim/projectM5/project/src/server/state"

func main() {
	app := &App{}
	app.InitDatabase()
	app.SetupLogin()
	app.CreateRoutes()
	gameManager := state.NewGameManager()
	app.Run(gameManager)
}

package main

import (
	app2 "github.com/risasim/projectM5/project/src/server/app"
	"github.com/risasim/projectM5/project/src/server/state"
)

func main() {
	app := &app2.App{}
	app.InitDatabase()
	app.SetupLogin()
	app.CreateRoutes()
	gameManager := state.NewGameManager()
	app.Run(gameManager)
}

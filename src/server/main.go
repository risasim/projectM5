package main

import (
	app "github.com/risasim/projectM5/project/src/server/app"
	"github.com/risasim/projectM5/project/src/server/state"
)

func main() {
	app := &app.App{}
	app.InitDatabase()
	app.SetupLogin()
	app.CreateRoutes()

	gameManager := state.NewGameManager()
	// run the broacasters in its own go routines
	go gameManager.BroadcastLeaderBoardHandler()
	go gameManager.BroadcastPisHandler()

	app.Run(gameManager)
}

package state

import "sync"

type GameManager struct {
	gameMode GameMode
	mutex    sync.Mutex
	//wsLeaderBoards
	//wsPis
}

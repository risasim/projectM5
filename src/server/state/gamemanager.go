package state

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
)

type GameManager struct {
	GameStatus GameStatus
	Mutex      sync.Mutex
	// WsLeaderBoards are the leaderboards web socket connections
	WsLeaderBoards map[*websocket.Conn]bool
	// BroadCastLeaderBoard is a channel to broadcast message to all of the leaderboards
	BroadCastLeaderBoard chan []byte
	// WsPis are teh websockets that are connected to pis
	WsPis map[*websocket.Conn]bool
	// BroadcastPis is a channel that will broadcast messages to all of the leaderboards -> all of them
	BroadcastPis chan []byte
}

// NewGameManager initializes a new GameManager
func NewGameManager() *GameManager {
	return &GameManager{
		GameStatus:           idle,
		WsLeaderBoards:       make(map[*websocket.Conn]bool),
		BroadCastLeaderBoard: make(chan []byte),
		WsPis:                make(map[*websocket.Conn]bool),
		BroadcastPis:         make(chan []byte),
	}
}

func (gm *GameManager) WsLeaderBoardHandler(c *gin.Context, upgrader *websocket.Upgrader) {

}

func (gm *GameManager) WsPisHandler(c *gin.Context, upgrader *websocket.Upgrader) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	gm.Mutex.Lock()
	gm.WsLeaderBoards[conn] = true
	gm.Mutex.Unlock()

}

// GameStatus is an enumaration of possible game statuses
type GameStatus int

const (
	idle GameStatus = iota
	active
)

var statusName = map[GameStatus]string{
	idle:   "idle",
	active: "active",
}

func (gs GameStatus) String() string {
	return statusName[gs]
}

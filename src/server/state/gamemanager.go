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
	upgrader     websocket.Upgrader
}

// NewGameManager initializes a new GameManager
func NewGameManager() *GameManager {
	return &GameManager{
		GameStatus:           idle,
		WsLeaderBoards:       make(map[*websocket.Conn]bool),
		BroadCastLeaderBoard: make(chan []byte),
		WsPis:                make(map[*websocket.Conn]bool),
		BroadcastPis:         make(chan []byte),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (gm *GameManager) WsLeaderBoardHandler(c *gin.Context) {
	conn, err := gm.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	gm.Mutex.Lock()
	gm.WsPis[conn] = true
	gm.Mutex.Unlock()

}

func (gm *GameManager) WsPisHandler(c *gin.Context) {
	conn, err := gm.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	gm.Mutex.Lock()
	gm.WsLeaderBoards[conn] = true
	gm.Mutex.Unlock()

}

func (gm *GameManager) handleConnection(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s\n", message)
		//React to message
		// switch based on the message

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

// BroadcastPisHandler does broadcast to all pis
func (gm *GameManager) BroadcastPisHandler() {
	for {
		// Grab the next message from the broadcast channel
		message := <-gm.BroadcastPis

		// Send the message to all connected clients
		gm.Mutex.Lock()
		for pi := range gm.WsPis {
			err := pi.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				pi.Close()
				delete(gm.WsPis, pi)
			}
		}
		gm.Mutex.Unlock()
	}
}

// BroadcastLeaderBoardHandler does broadcast the message about the game state to all of involved
// connections
func (gm *GameManager) BroadcastLeaderBoardHandler() {
	for {
		// Grab the next message from the broadcast channel
		message := <-gm.BroadCastLeaderBoard

		// Send the message to all connected clients
		gm.Mutex.Lock()
		for client := range gm.WsLeaderBoards {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(gm.WsLeaderBoards, client)
			}
		}
		gm.Mutex.Unlock()
	}
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

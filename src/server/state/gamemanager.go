package state

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/risasim/projectM5/project/src/server/communication"
)

type GameManager struct {
	GameStatus     GameStatus
	Mutex          sync.Mutex
	CurrentSession *Session
	// WsLeaderBoards are the leaderboards web socket connections
	WsLeaderBoards map[*websocket.Conn]bool
	// BroadCastLeaderBoard is a channel to broadcast message to all of the leaderboards
	BroadCastLeaderBoard chan []byte
	// WsPis are teh websockets that are connected to pis
	WsPis map[*websocket.Conn]bool
	// BroadcastPis is a channel that will broadcast messages to all of the leaderboards -> all of them
	BroadcastPis chan []byte
	// Game is actual game data
	Game     GameMode
	upgrader websocket.Upgrader
}

// NewGameManager initializes a new GameManager
func NewGameManager() *GameManager {
	return &GameManager{
		GameStatus:           Idle,
		Mutex:                sync.Mutex{},
		CurrentSession:       NewSession(),
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

// CreateNewGame starts a new game session
func (gm *GameManager) CreateNewGame(gameType communication.GameType) error {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	if gm.GameStatus == Created {
		return fmt.Errorf("a game is already active")
	}
	if gm.GameStatus == Started {
		return fmt.Errorf("a game has already started")
	}

	// Initialise a new game session
	gm.CurrentSession = &Session{
		player:   []Player{},
		hitData:  []communication.HitData{},
		GameType: gameType,
	}

	gm.GameStatus = Created
	fmt.Println("Game created")
	return nil
}

func (gm *GameManager) StartGame() error {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()
	gm.GameStatus = Started
	gm.Game.startGame(gm.CurrentSession)
	startMessage := communication.StartedMessage{At: time.Now(), Active: true}
	jsonData, err := json.Marshal(startMessage)
	if err != nil {
		return fmt.Errorf("json messud up innit")
	}
	gm.BroadcastPis <- jsonData
	return nil
}

func (gm *GameManager) EndGame() error {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	gm.GameStatus = Idle
	gm.CurrentSession = nil
	endMessage := communication.EndedMessage{At: time.Now()}
	jsonData, err := json.Marshal(endMessage)
	if err != nil {
		return fmt.Errorf("json fuckup")
	}
	gm.BroadcastPis <- jsonData
	fmt.Println("Game ended")
	return nil
}

// AddPlayer to add a player to the current game session
func (gm *GameManager) AddPlayer(player Player) error {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	if gm.GameStatus == Idle {
		return fmt.Errorf("there is no game to join")
	}
	if gm.GameStatus == Started {
		return fmt.Errorf("a game has already started")
	}

	// Checking if a player is already in the game session
	for _, p := range gm.CurrentSession.player {
		if p.PiSN == player.PiSN {
			return fmt.Errorf("a player with this ID is already in the game")
		}
	}

	gm.CurrentSession.player = append(gm.CurrentSession.player, player)
	return nil
}

// RemovePlayer to remove a player from the current game session
func (gm *GameManager) RemovePlayer(player Player) error {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()
	for i, p := range gm.CurrentSession.player {
		if p.PiSN == player.PiSN {
			gm.CurrentSession.player = append(gm.CurrentSession.player[:i], gm.CurrentSession.player[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("error trying to remove player with this ID")
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
	gm.WsPis[conn] = true
	gm.Mutex.Unlock()

	go gm.handlePiConnection(conn)

}

// handlePiConnection does listen to being hit and in case that
func (gm *GameManager) handlePiConnection(conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		fmt.Printf("Received: %s\n", message)

		var hitData communication.HitData
		if err := json.Unmarshal(message, &hitData); err != nil {
			fmt.Println("Error unmarshalling hit data:", err)
			continue
		}

		res := gm.Game.registerHit(hitData)

		responseJSON, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Error marshalling response:", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

// updateLeaderBoard does send the generated data about the game into the broadcast of the leaderboards
func (gm *GameManager) updateLeaderBoard() {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()
	update := gm.Game.generateData()
	responseJSON, err := json.Marshal(update)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return
	}
	gm.BroadCastLeaderBoard <- responseJSON
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
	Idle GameStatus = iota
	Created
	Started
)

var statusName = map[GameStatus]string{
	Idle:    "idle",
	Created: "created",
	Started: "started",
}

func (gs GameStatus) String() string {
	return statusName[gs]
}

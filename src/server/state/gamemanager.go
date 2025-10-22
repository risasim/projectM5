package state

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/risasim/projectM5/project/src/server/communication"
	"github.com/risasim/projectM5/project/src/server/db"
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
	WsPis map[string]*websocket.Conn
	// BroadcastPis is a channel that will broadcast messages to all of the leaderboards -> all of them
	BroadcastPis   chan []byte
	userRepository db.UserRepositoryInterface
}

// NewGameManager initializes a new GameManager
func NewGameManager(gameType communication.GameType, repo db.UserRepositoryInterface) *GameManager {
	return &GameManager{
		GameStatus:           idle,
		Mutex:                sync.Mutex{},
		CurrentSession:       nil,
		WsLeaderBoards:       make(map[*websocket.Conn]bool),
		BroadCastLeaderBoard: make(chan []byte),
		WsPis:                make(map[string]*websocket.Conn),
		BroadcastPis:         make(chan []byte),
		userRepository:       repo,
	}
}

// StartNewGame starts a new game session
func (gm *GameManager) StartNewGame(gameType communication.GameType) error {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	if gm.GameStatus != idle {
		return fmt.Errorf("a game is already in active")
	}

	// Initialise a new game session
	gm.CurrentSession = &Session{
		player:   []Player{},
		hitData:  []communication.HitData{},
		GameType: gameType,
	}

	gm.GameStatus = Active
	startMessage := communication.StartedMessage{At: time.Now(), Active: true}
	gm.BroadcastToPis(communication.Start, startMessage)
	fmt.Println("Game started")
	return nil
}

func (gm *GameManager) EndGame() {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	gm.GameStatus = idle
	gm.CurrentSession = nil
	endMessage := communication.EndedMessage{At: time.Now()}
	gm.BroadcastToPis(communication.End, endMessage)
	fmt.Println("Game ended")
}

// AddPlayer to add a player to the current game session
func (gm *GameManager) AddPlayer(player Player) error {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	if gm.GameStatus != Active {
		return fmt.Errorf("a game is already in active")
	}

	// Checking if a player is already in the game session
	for _, p := range gm.CurrentSession.player {
		if p.ID == player.ID {
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
		if p.ID == player.ID {
			gm.CurrentSession.player = append(gm.CurrentSession.player[:i], gm.CurrentSession.player[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("error trying to remove player with this ID")
}

// BroadcastToPis to broadcast a message to all Pis
func (gm *GameManager) BroadcastToPis(messageType communication.MsgType, data any) {
	message := map[string]any{
		"messageType": messageType,
		"data":        mustJson(data),
	}

	payload, _ := json.Marshal(message)

	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()
	for _, conn := range gm.WsPis {
		conn.WriteMessage(websocket.TextMessage, payload)
	}
}

// SendNewMusicToPi sends new death sound to pi
func (gm *GameManager) SendNewMusicToPi(username string, b64Sound string, fileName string) error {
	gm.Mutex.Lock()
	conn, ok := gm.WsPis[username]
	gm.Mutex.Unlock()
	if !ok {
		return fmt.Errorf("pi not connected for %s", username)
	}
	message := communication.SetSoundMessage{
		SoundName: fileName,
		Base64:    b64Sound,
	}

	data, _ := json.Marshal(communication.Message{
		MsgType: communication.SetSound,
		Data:    mustJson(message),
	})

	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()
	return conn.WriteMessage(websocket.TextMessage, data)
}

func mustJson(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
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
	gm.WsPis[conn] = true
	gm.Mutex.Unlock()

}

// GameStatus is an enumaration of possible game statuses
type GameStatus int

const (
	idle GameStatus = iota
	Active
)

var statusName = map[GameStatus]string{
	idle:   "idle",
	Active: "active",
}

func (gs GameStatus) String() string {
	return statusName[gs]
}

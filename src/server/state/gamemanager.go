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
		Game:                 NewFreeForAll(NewSession()),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// GenericMessage for checking message type before processing
type GenericMessage struct {
	Type string `json:"type"`
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
		Player:   make([]Player, 0),
		hitData:  make([]communication.HitData, 0),
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
	if gm.CurrentSession.GameType == communication.Infected {
		startMessage = communication.StartedMessage{At: time.Now(), Active: false}
	}
	jsonData, err := json.Marshal(startMessage)
	if err != nil {
		return fmt.Errorf("json messud up innit")
	}
	message := communication.Message{MsgType: communication.Start, Data: jsonData}
	jsonData2, err2 := json.Marshal(message)
	if err2 != nil {
		return fmt.Errorf("json messud up innit")
	}
	gm.BroadcastPis <- jsonData2
	return nil
}

func (gm *GameManager) EndGame() error {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	gm.GameStatus = Idle
	gm.CurrentSession = NewSession()
	endMessage := communication.EndedMessage{At: time.Now()}
	jsonData, err := json.Marshal(endMessage)
	if err != nil {
		return fmt.Errorf("json messud up innit")
	}
	message := communication.Message{MsgType: communication.End, Data: jsonData}
	jsonData2, err2 := json.Marshal(message)
	if err2 != nil {
		return fmt.Errorf("json messud up innit")
	}
	gm.BroadcastPis <- jsonData2
	fmt.Println("Game ended")
	return nil
}

// AddPlayer to add a Player to the current game session
func (gm *GameManager) AddPlayer(player Player) error {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	if gm.GameStatus == Idle {
		return fmt.Errorf("there is no game to join")
	}
	if gm.GameStatus == Started {
		return fmt.Errorf("a game has already started")
	}

	// Checking if a Player is already in the game session
	for _, p := range gm.CurrentSession.Player {
		if p.PiSN == player.PiSN {
			return fmt.Errorf("a Player with this ID is already in the game")
		}
	}

	gm.CurrentSession.Player = append(gm.CurrentSession.Player, player)
	return nil
}

func (gm *GameManager) SessionPlayers() []string {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()
	if gm.CurrentSession.Player == nil || len(gm.CurrentSession.Player) == 0 {
		return []string{}
	}
	usernames := make([]string, len(gm.CurrentSession.Player))
	for i, p := range gm.CurrentSession.Player {
		usernames[i] = p.Username
	}
	return usernames
}

// RemovePlayer to remove a Player from the current game session
func (gm *GameManager) RemovePlayer(player Player) error {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()
	for i, p := range gm.CurrentSession.Player {
		if p.PiSN == player.PiSN {
			gm.CurrentSession.Player = append(gm.CurrentSession.Player[:i], gm.CurrentSession.Player[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("error trying to remove Player with this ID")
}

func (gm *GameManager) WsPisHandler(c *gin.Context) {
	conn, err := gm.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	gm.Mutex.Lock()
	gm.WsPis[conn] = true
	gm.Mutex.Unlock()

	gm.handlePiConnection(conn)
}

func (gm *GameManager) WsLeaderBoardHandler(c *gin.Context) {
	conn, err := gm.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	gm.Mutex.Lock()
	gm.WsLeaderBoards[conn] = true
	gm.Mutex.Unlock()

	gm.handleLeaderBoardConnection(conn)
}

// handleLeaderBoardConnection maintains connection with leaderboard clients
func (gm *GameManager) handleLeaderBoardConnection(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		gm.Mutex.Lock()
		delete(gm.WsLeaderBoards, conn)
		gm.Mutex.Unlock()
		fmt.Println("Leaderboard connection closed and removed from pool")
	}()

	// Set read deadline and pong handler for keepalive
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Channel to signal goroutine to stop
	done := make(chan struct{})
	defer close(done)

	// Send a ping every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Handle pings in a separate goroutine
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					return
				}
				//fmt.Println("Sent ping to leaderboard client")
			case <-done:
				return
			}
		}
	}()

	// Listen for messages from leaderboard (usually just pings/pongs)
	if setupMessage := gm.getLeaderboardMessage(); setupMessage != nil {
		if err := conn.WriteMessage(websocket.TextMessage, setupMessage); err != nil {
			fmt.Println("Error writing message:", err)
		}
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("Leaderboard client disconnected normally")
			} else {
				fmt.Println("Leaderboard connection error:", err)
			}
			break
		}

		// Check if it's a ping/keepalive message
		var genericMsg GenericMessage
		if err := json.Unmarshal(message, &genericMsg); err == nil {
			if genericMsg.Type == "ping" {
				//fmt.Println("Received keepalive ping from leaderboard client")
				conn.SetReadDeadline(time.Now().Add(60 * time.Second))
				continue
			}
		}

		// Log any other messages (leaderboards typically only receive broadcasts)
		fmt.Printf("Leaderboard sent: %s\n", message)
	}
}

// handlePiConnection does listen to being hit and in case that
func (gm *GameManager) handlePiConnection(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		gm.Mutex.Lock()
		delete(gm.WsPis, conn)
		gm.Mutex.Unlock()
		fmt.Println("Pi connection closed and removed from pool")
	}()

	// Set read deadline and pong handler for keepalive
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		//fmt.Println("Received pong from client")
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Channel to signal goroutine to stop
	done := make(chan struct{})
	defer close(done)

	// Send a ping every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Handle pings in a separate goroutine with proper cleanup
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					fmt.Println("Error sending ping:", err)
					return
				}
				//fmt.Println("Sent ping to client")
			case <-done:
				fmt.Println("Stopping Pi ping sender")
				return
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("Client disconnected normally")
			} else {
				fmt.Println("Error reading message:", err)
			}
			break
		}

		fmt.Printf("Received: %s\n", message)

		// First, check if it's a ping/keepalive message
		var genericMsg GenericMessage
		if err := json.Unmarshal(message, &genericMsg); err == nil {
			if genericMsg.Type == "ping" {
				//fmt.Println("Received keepalive ping from client")
				// Reset read deadline on ping
				conn.SetReadDeadline(time.Now().Add(60 * time.Second))
				continue // Skip processing pings as hit data
			}
		}

		// Try to parse as hit data
		var hitData communication.HitData
		if err := json.Unmarshal(message, &hitData); err != nil {
			fmt.Println("Error unmarshalling hit data:", err)
			// Send error response back to client
			errorResponse := map[string]string{
				"error":   "invalid_format",
				"message": "Could not parse hit data",
			}
			responseJSON, _ := json.Marshal(errorResponse)
			conn.WriteMessage(websocket.TextMessage, responseJSON)
			continue
		}

		// Process the hit
		res := gm.Game.registerHit(hitData)
		gm.updateLeaderBoard()

		responseJSON, err1 := json.Marshal(res)
		if err1 != nil {
			fmt.Println("Error marshalling response:", err1)
			continue
		}
		message1 := communication.Message{MsgType: communication.HitResponseMsg, Data: responseJSON}
		jsonData2, err2 := json.Marshal(message1)
		if err2 != nil {
			fmt.Println("Error marshalling response:", err2)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, jsonData2); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}

		fmt.Printf("Processed hit from %s\n", hitData.Victim)
	}
}

// updateLeaderBoard does send the generated data about the game into the broadcast of the leaderboards
func (gm *GameManager) updateLeaderBoard() {
	responseJSON := gm.getLeaderboardMessage()
	if responseJSON != nil {
		gm.BroadCastLeaderBoard <- responseJSON
	}
}

func (gm *GameManager) getLeaderboardMessage() []byte {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()
	update := gm.Game.generateData()
	responseJSON, err := json.Marshal(update)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return nil
	}
	return responseJSON
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

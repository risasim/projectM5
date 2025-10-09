package state

import (
	"gitlab.utwente.nl/computer-systems-project/2025-2026/students-projects/cs25-29/project/src/server/communication"
	"time"
)

// gameMode does prescribe functions that all of the gameModes share
type GameMode interface {
	registerHit() communication.HitResponse
	generateData()
	finished() bool
}

// Hit does store data passed from the pi about a hit
type Hit struct {
	// victim is the id of the Player who is sending the data about being shot
	victim int
	// shooter is the id of the Player who has sent the IR signal and hit the receiver
	shooter int
	// timeStamp of the shot being registered
	timeStamp time.Time
}

// FreeForAll is a game mode where the players are competing against each other, without reviving
type FreeForAll struct {
	// deadPeople stores the people that have already being killed
	deadPeople []Player
	// session that is the GameMode played in
	session Session
}

func (ffl FreeForAll) finished() bool {
	return len(ffl.deadPeople) == len(ffl.session.player)
}

// TeamDeathMatch is a GameMode where players in teams compete to eliminate each-other for 30 minutes
// where people revive after a certain amount of time after dying
type TeamDeathMatch struct {
	// time is session time in minutes
	time int
	// teams is an Array of all teams in the session
	teams []Team
	// session that is the GameMode played in
	session Session
}

func (tdm TeamDeathMatch) finished() bool {
	return false
}

// Team are the collaborating players,they cannot kill each other
type Team struct {
	// name of the team
	name string
	// members are players in the game
	members []Player
	// score of the team (-100 for being killed, 200 for kill)
	score int
}

type Infected struct {
	//infectedPeople that stores the people infected
	infectedPeople []Player
	// session that is the GameMode played in
	session Session
}

// Session does hold the common
type Session struct {
	player  []Player
	status  GameStatus
	hitData []Hit
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

// Holds basic player information
type Player struct {
	// username is the unique name of the user
	username string
	// encodingID is the integer that is encoded by the IR of the assigned PI
	encodingID uint
	// id is the unique integer of the player
	id int
}

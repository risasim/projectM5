package state

import (
	"github.com/risasim/projectM5/project/src/server/communication"
)

// GameMode does prescribe functions that all of the GameModes share
type GameMode interface {
	// registerHit() is a function to respond to getting a hit alert from the infrared receiver
	registerHit() communication.HitResponse
	// generateData is a function to generate the leaderboard data
	generateData()
	// finished is the function to determine if any GameMode is finished
	finished() bool
}

// FreeForAll is a game mode where the players are competing against each other, without reviving
type FreeForAll struct {
	// deadPeople stores the people that have already being killed
	deadPeople []Player
	// session that is the GameMode played in
	session Session
}

func (ffl FreeForAll) registerHit() {

}

// finished returns true if the array length of dead people matches the array length of the player array in the session
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

// finished is the condition to determine if the TeamDeathMatch GameMode is finished
func (tdm TeamDeathMatch) finished() bool {
	// Checking the status of the session to determine if the game is finished
	if tdm.session.status == idle {
		return true
	}
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

// Infected is a GameMode where one person starts the game being infected and their task in to
// infect everyone else by shooting them and once another person gets infected they join the
// infector group to infect others! Infectors cannot die!
type Infected struct {
	//infectedPeople that stores the people infected
	infectedPeople []Player
	// session that is the GameMode played in
	session Session
}

// finished returns true if the array length of infected people matches the array length of the player array in the session
func (inf Infected) finished() bool {
	return len(inf.infectedPeople) == len(inf.session.player)
}

// Session does hold the common
type Session struct {
	player  []Player
	status  GameStatus
	hitData []communication.HitData
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

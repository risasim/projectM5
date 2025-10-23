package state

import (
	"sort"

	"github.com/risasim/projectM5/project/src/server/communication"
)

// GameMode does prescribe functions that all of the GameModes share
type GameMode interface {
	// registerHit() is a function to respond to getting a hit alert from the infrared receiver
	registerHit(dt communication.HitData) communication.HitResponse
	// generateData is a function to generate the leaderboard data
	generateData() communication.LeaderboardMessage
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

func (ffl FreeForAll) registerHit(dt communication.HitData) communication.HitResponse {
	//TODO implement me
	panic("implement me")
}

// generateDate Ensures everyone is in the dead poeple array and then reverses it for the leaderboard
func (ffl FreeForAll) generateData() []Player {
	if len(ffl.deadPeople) != len(ffl.session.player) {
		for _, player := range ffl.session.player {
			found := false
			for _, deadPlayer := range ffl.deadPeople {
				if deadPlayer.ID == player.ID {
					found = true
					break
				}
			}
			if !found {
				ffl.deadPeople = append(ffl.deadPeople, player)
			}
		}
	}
	for i, j := 0, len(ffl.deadPeople)-1; i < j; i, j = i+1, j-1 {
		ffl.deadPeople[i], ffl.deadPeople[j] = ffl.deadPeople[j], ffl.deadPeople[i]
	}
	return ffl.deadPeople
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

func (tdm TeamDeathMatch) registerHit(dt communication.HitData) communication.HitResponse {
	//TODO implement me
	panic("implement me")
}

// generateData sorts teams by score and returns team names in order
func (tdm TeamDeathMatch) generateData() []string {
	sort.Slice(tdm.teams, func(i, j int) bool {
		return tdm.teams[i].score > tdm.teams[j].score
	})
	var teamOrderedName []string
	for _, team := range tdm.teams {
		teamOrderedName = append(teamOrderedName, team.name)
	}
	return teamOrderedName
}

// finished is the condition to determine if the TeamDeathMatch GameMode is finished
func (tdm TeamDeathMatch) finished() bool {
	// Checking the status of the session to determine if the game is finished
	//if tdm.session.status == idle {
	//	return true
	//}
	//return false
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

func (inf Infected) registerHit(dt communication.HitData) communication.HitResponse {
	//TODO implement me
	panic("implement me")
}

// generateData returns reversed list of infected people as people are added as they get infected
func (inf Infected) generateData() []Player {
	for i, j := 0, len(inf.infectedPeople)-1; i < j; i, j = i+1, j-1 {
		inf.infectedPeople[i], inf.infectedPeople[j] = inf.infectedPeople[j], inf.infectedPeople[i]
	}
	return inf.infectedPeople
}

// finished returns true if the array length of infected people matches the array length of the player array in the session
func (inf Infected) finished() bool {
	return len(inf.infectedPeople) == len(inf.session.player)
}

// Session does hold the common
type Session struct {
	player   []Player
	hitData  []communication.HitData
	GameType communication.GameType
}

// Holds basic player information
type Player struct {
	// Username is the unique name of the user
	Username string
	// EncodingID is the integer that is encoded by the IR of the assigned PI
	EncodingID uint
	// ID is the unique integer of the player
	ID int
	// DeathSound is the name of the file that needs to add
	DeathSound string
}

package state

import (
	"encoding/json"
	"fmt"
	"github.com/risasim/projectM5/project/src/server/communication"
)

// GameMode does prescribe functions that all of the GameModes share
type GameMode interface {
	// startGame does all that is needed to run the game
	startGame()
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

// registerHit in freefall does add the user to the death people without reviving
func (ffl *FreeForAll) registerHit(dt communication.HitData) communication.HitResponse {
	for i := range ffl.session.player {
		if ffl.session.player[i].PiSN == dt.Victim {
			ffl.deadPeople = append(ffl.deadPeople, ffl.session.player[i])
			return communication.HitResponse{
				PlaySound: true,
				SoundName: ffl.session.player[i].DeathSound,
				Dead:      true,
				Revive:    false,
				ReviveIn:  0,
			}
		}
	}
	return communication.HitResponse{}
}

// generateData Ensures everyone is in the dead poeple array and then reverses it for the leaderboard
func (ffl *FreeForAll) generateData() communication.LeaderboardMessage {
	dead := make([]communication.LeaderboardPlayer, len(ffl.deadPeople))
	for i, player := range ffl.deadPeople {
		dead[i] = communication.LeaderboardPlayer{
			Username: player.Username,
		}
	}

	details := communication.FreefallLeaderboard{
		DeadPlayers: dead,
	}

	players := make([]communication.LeaderboardPlayer, len(ffl.session.player))
	for i, player := range ffl.session.player {
		players[i] = communication.LeaderboardPlayer{
			Username: player.Username,
		}
	}

	jsonRaw, err := json.Marshal(details)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
	}

	res := communication.LeaderboardMessage{
		GameType: communication.Freefall,
		Data:     jsonRaw,
		Players:  players,
	}
	return res
}

// finished returns true if the array length of dead people matches the array length of the player array in the session
func (ffl *FreeForAll) finished() bool {
	return len(ffl.deadPeople) == len(ffl.session.player)
}

// TeamDeathMatch is a GameMode where players in teams compete to eliminate each-other for 30 minutes
// where people revive after a certain amount of time after dying
type TeamDeathMatch struct {
	// time is session time in minutes
	time int
	// divisions is map from the players to the
	divisions map[string]*Team
	// teams is an Array of all teams in the session
	teams []Team
	// session that is the GameMode played in
	session Session
}

func (tdm *TeamDeathMatch) registerHit(dt communication.HitData) communication.HitResponse {
	tdm.divisions[dt.Victim].score -= 100
	for i := range tdm.session.player {
		if tdm.session.player[i].PiSN == dt.Victim {
			return communication.HitResponse{
				PlaySound: true,
				SoundName: tdm.session.player[i].DeathSound,
				Dead:      true,
				Revive:    true,
				ReviveIn:  30,
			}
		}
	}
	return communication.HitResponse{}
}

// generateData sorts teams by score and returns team names in order
func (tdm *TeamDeathMatch) generateData() communication.LeaderboardMessage {
	teams := make([]communication.DeathMatchTeam, len(tdm.teams))

	for i, team := range tdm.teams {
		members := make([]communication.LeaderboardPlayer, len(team.members))
		for _, member := range team.members {
			members = append(members, communication.LeaderboardPlayer{
				Username: member.Username,
			})
		}
		teams[i] = communication.DeathMatchTeam{
			Name:    team.name,
			Members: members,
			Score:   team.score,
		}
	}

	data := communication.TeamDeathMatchLeaderboard{
		Teams: teams,
	}

	players := make([]communication.LeaderboardPlayer, len(tdm.session.player))
	for i, player := range tdm.session.player {
		players[i] = communication.LeaderboardPlayer{
			Username: player.Username,
		}
	}

	jsonRaw, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
	}

	res := communication.LeaderboardMessage{
		GameType: communication.TeamDeathmatch,
		Data:     jsonRaw,
		Players:  players,
	}

	return res
}

// finished is the condition to determine if the TeamDeathMatch GameMode is finished
func (tdm *TeamDeathMatch) finished() bool {
	for i := range tdm.teams {
		if tdm.teams[i].score <= 0 {
			return true
		}
	}
	return false
}

// startGame does initilise the game by splitting the users into two teams
func (tdm *TeamDeathMatch) startGame() {
	team1 := Team{
		score:   1500,
		name:    "kittens",
		members: make([]Player, 0),
	}

	team2 := Team{
		score:   1500,
		name:    "mittens",
		members: make([]Player, 0),
	}
	for i := range tdm.session.player {
		if i%2 == 0 {
			team1.members = append(team1.members, tdm.session.player[i])
		} else {
			team2.members = append(team2.members, tdm.session.player[i])
		}
	}
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

func (inf *Infected) registerHit(dt communication.HitData) communication.HitResponse {
	//TODO implement me
	panic("implement me")
}

// generateData returns reversed list of infected people as people are added as they get infected
func (inf *Infected) generateData() []Player {
	for i, j := 0, len(inf.infectedPeople)-1; i < j; i, j = i+1, j-1 {
		inf.infectedPeople[i], inf.infectedPeople[j] = inf.infectedPeople[j], inf.infectedPeople[i]
	}
	return inf.infectedPeople
}

// finished returns true if the array length of infected people matches the array length of the player array in the session
func (inf *Infected) finished() bool {
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
	// ID is the unique integer of the player
	PiSN string
	// DeathSound is the name of the file that needs to add
	DeathSound string
}

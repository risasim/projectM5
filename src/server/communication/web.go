package communication

import (
	"encoding/json"
	"fmt"
	"strings"
)

// LeaderboardMessage is the struct for json that will be sent to the web
type LeaderboardMessage struct {
	GameType GameType `json:"game_type"`
	// It will carrt different type of data
	Data    json.RawMessage     `json:"data"`
	Players []LeaderboardPlayer `json:"players"`
}

// FreefallLeaderboard is leaderboard sent to
type FreefallLeaderboard struct {
	// deadPlayers are the players that has been killed
	DeadPlayers []LeaderboardPlayer
}

// TeamDeathMatchLeaderboard is for deathmatch passing to json
type TeamDeathMatchLeaderboard struct {
	Teams []DeathMatchTeam `json:"teams"`
}

// DeathMarchTeam
type DeathMatchTeam struct {
	Name    string `json:"name"`
	Members []LeaderboardPlayer
	Score   int `json:"score"`
}

// InfectedLeaderboard has only the infected users in them
type InfectedLeaderboard struct {
	Infected []LeaderboardPlayer `json:"infected"`
}

// LeaderboardPlayer is player with only the data needed for
type LeaderboardPlayer struct {
	Username string `json:"username"`
}

// GameType is pretty much GameMode but only for decoding
type GameType uint

// Changed something here, need to check later if that is fine, also futher from MsgType to GameType
const (
	Freefall GameType = iota
	TeamDeathmatch
	Infected
)

var (
	GameType_String = map[uint]string{
		0: "Freefall",
		1: "TeamDeathmatch",
		2: "Infected",
	}
	GameType_value = map[string]uint{
		"Freefall":       0,
		"TeamDeathmatch": 1,
		"Infected":       2,
	}
)

func ParseGameType(gmt string) (GameType, error) {
	gmt = strings.TrimSpace(gmt)
	value, ok := GameType_value[gmt]
	if !ok {
		return GameType(0), fmt.Errorf("invalid msg type: %s", gmt)
	}
	return GameType(value), nil
}

func (gmt GameType) String() string {
	return GameType_String[uint(gmt)]
}

// MarshalJSON is encoding to JSON
func (gmt GameType) MarshalJSON() ([]byte, error) {
	return json.Marshal(gmt.String())
}

// UnmarshalJSON is for decoding the msg type enum
func (gmt *GameType) UnmarshalJSON(data []byte) (err error) {
	var gts string
	if err := json.Unmarshal(data, &gts); err != nil {
		return err
	}
	if *gmt, err = ParseGameType(gts); err != nil {
		return err
	}
	return nil
}

package communication

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// StartedMessage is message that should be sent at the start of each game
type StartedMessage struct {
	At time.Time `json:"at"`
	// Active is flag to determine if the pi is active or waits for activation
	Active bool `json:"active"`
}

// EndedMessage is message sent after the end of the game to all players
type EndedMessage struct {
	At time.Time `json:"at"`
}

// HitResponse is what will be sent from server to pi after hit
type HitResponse struct {
	// PlaySound should the pi make a sound
	PlaySound bool `json:"playSound"`
	// SoundName is the name of the file
	SoundName string `json:"soundName"`
	// Dead is whether the pi should still play or not
	Dead bool `json:"dead"`
	// Revive should the pi revive ?
	Revive bool `json:"revive"`
	// ReviveIn if it should revive in what time, after which it should shoot again
	ReviveIn int `json:"reviveIn"`
}

// HitData does store data passed from the pi about a hit
type HitData struct {
	// victim is the id of the Player who is sending the data about being shot
	Victim int `json:"victim"`
	// shooter is the id of the Player who has sent the IR signal and hit the receiver
	Shooter int `json:"shooter"`
	// timeStamp of the shot being registered
	TimeStamp time.Time `json:"timestamp"`
}

// Message is the actual type that will be decoded into
type Message struct {
	MsgType MsgType         `json:"msgtype"`
	Data    json.RawMessage `json:"Data"`
}

// MsgType is an enum that will determine the rest of the message
type MsgType uint

const (
	Auth MsgType = iota
	Start
	HitDataMsg
	HitResponseMsg
	End
)

var (
	MsgType_String = map[uint]string{
		0: "Auth",
		1: "Start",
		2: "HitDataMsg",
		3: "HitResponseMsg",
		4: "End",
	}
	MsgType_value = map[string]uint{
		"Auth":           0,
		"Start":          1,
		"HitDataMsg":     2,
		"HitResponseMsg": 3,
		"End":            4,
	}
)

func ParseMsgType(msg string) (MsgType, error) {
	msg = strings.TrimSpace(strings.ToLower(msg))
	value, ok := MsgType_value[msg]
	if !ok {
		return MsgType(0), fmt.Errorf("invalid msg type: %s", msg)
	}
	return MsgType(value), nil
}

func (mst MsgType) String() string {
	return MsgType_String[uint(mst)]
}

// MarshalJSON is encoding to JSON
func (mst MsgType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mst.String())
}

// UnmarshalJSON is for decoding the msg type enum
func (mst *MsgType) UnmarshalJSON(data []byte) (err error) {
	var msg string
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	if *mst, err = ParseMsgType(msg); err != nil {
		return err
	}
	return nil
}

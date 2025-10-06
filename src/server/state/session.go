package state

type session struct {
	// Version of the current session
	Version uint
	// Players that are connected to the game
	Players []player
	// GameMode defines the mode that the game logic should be derived from
	GameMode string
}

// FreeFall, Infected, Team DeathcMatch, Hot potato?

type player struct {
	// Username is taken username from the db
	Username string
	// NumOfHits is number of times the users Pi has detected
	NumOfHits uint
}

package model

// GetUser represents the sql model of the user in db
type GetUserResponse struct {
	// ID unique user id
	ID uint `json:"id"`
	// Username of the user
	Username string `json:"username"`
	// IsAdmin is a flag if the user is admin
	IsAdmin bool `json:"isAdmin"`
}

// GetUserAuth is used to retrieve and compare data for authentication purposes
type GetUserAuth struct {
	// ID unique user id
	ID uint
	// Username of the user
	Username string
	// IsAdmin is a flag if the user is admin
	IsAdmin bool
	// Password of the user
	Password string
	// DeathSound is the link to the file with the death sound
	DeathSound string
	// ApiKey generated for the user
	ApiKey string
	// PiSN is the serial number of Pi
	PiSN string
}

// PostUser is a struct for the create request adding new user
type PostUser struct {
	// Username of the user
	Username string `json:"username"`
	// IsAdmin is a flag if the user is admin
	IsAdmin bool `json:"isAdmin"`
	// Password of the user
	Password string `json:"password"`
	// PiSN is the serial number of Pi
	PiSN string `json:"piSN"`
}

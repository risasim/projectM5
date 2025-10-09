package model

// GetUser represents the sql model of the user in db
type GetUser struct {
	// ID unique user id
	ID uint
	// Username of the user
	Username string
	// IsAdmin is a flag if the user is admin
	IsAdmin bool
	// Password of the user
	Password string
	// ApiKey generated for the user
	ApiKey string
	// PiSN is the serial number of Pi
	PiSN string
}

// CreateUser is a struct for the create request adding new user
type CreateUser struct {
	// Username of the user
	Username string
	// IsAdmin is a flag if the user is admin
	IsAdmin bool
	// Password of the user
	Password string
	// PiSN is the serial number of Pi
	PiSN string
}

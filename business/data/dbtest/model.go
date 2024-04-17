package dbtest

import "github.com/EnesDemirtas/medisync/business/core/crud/userbus"

// User represents an app user specified for the test.
type User struct {
	userbus.User
	Token string
}

// SeedData represents data that was seeded for the test.
type SeedData struct {
	Users  []User
	Admins []User
}
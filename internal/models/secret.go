package models

type Secret struct {
	Id      string // access token should be string
	OwnerId string
	Secret  string
}

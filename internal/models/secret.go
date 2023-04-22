package models

type Secret struct {
	Id      string `bson:"_id"` // access token should be string
	OwnerId string `bson:"user_id"`
	Token   string `bson:"token"`
}

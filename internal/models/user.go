package models

type User struct {
	Id       string
	Email    string
	Username string
	Password string `json:"-" bson:"Password"`
}

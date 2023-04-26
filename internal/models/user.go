package models

import "errors"

type User struct {
	Id       string `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func IsValidUser(user User) error {
	if len(user.Username) == 0 {
		return errors.New("expected an email or unique identifier")
	}

	if len(user.Password) == 0 {
		return errors.New("expected an password or unique identifier")
	}

	return nil
}

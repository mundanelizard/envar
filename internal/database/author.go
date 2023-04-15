package database

import (
	"fmt"
	"time"
)

type Author struct {
	id        string
	name      string
	email     string
	timestamp time.Time
}

func NewAuthor(name, email string, timestamp time.Time) *Author {
	if len(name) == 0 {
		name = "John Doe"
	}

	if len(email) == 0 {
		email = "john.doe@envi.org"
	}

	return &Author{
		name:      name,
		email:     email,
		timestamp: timestamp,
	}
}

func (aut *Author) Id() string {
	return aut.id
}

func (aut *Author) SetId(id string) {
	aut.id = id
}

func (aut *Author) String() string {
	return fmt.Sprintf("%s <%s> %s", aut.name, aut.email, aut.timestamp.Format(time.RFC3339))
}

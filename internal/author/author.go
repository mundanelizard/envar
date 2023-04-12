package author

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

func New(name, email string, timestamp time.Time) *Author {
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

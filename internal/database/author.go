package database

import (
	"fmt"
	"strings"
	"time"
)

type Author struct {
	name      string
	email     string
	timestamp time.Time
}

func NewAuthor(name string, timestamp time.Time) *Author {
	if len(name) == 0 {
		name = "John Doe"
	}

	email := "john.doe@envi.org"

	return &Author{
		name:      name,
		email:     email,
		timestamp: timestamp,
	}
}

func NewAuthorFromByteArray(data string) (*Author, error) {
	chunks := strings.Split(data, " - ")
	// email := strings.ReplaceAll(chunks[1], "<", "")
	// email = strings.ReplaceAll(email, ">", "")
	timestamp, err := time.Parse(time.RFC3339, strings.Trim(chunks[2], ""))
	if err != nil {
		return nil, err
	}

	author := NewAuthor(chunks[0], timestamp)

	return author, nil
}

func (aut *Author) String() string {
	return fmt.Sprintf("%s - <%s> - %s", aut.name, aut.email, aut.timestamp.Format(time.RFC3339))
}

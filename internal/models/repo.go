package models

type Contributor struct {
	OwnerId string
	Role    string // admin, read, write
}

type Repo struct {
	Id           string
	Name         string // username:repo
	OwnerId      string
	Secret       string
	Contributors []Contributor
}

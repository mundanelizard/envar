package models

type Contributor struct {
	OwnerId string
	Role    string // admin, read, write
}

type Repo struct {
	Id           string        `bson:"_id"`
	Name         string        `bson:"name"` // username:repo
	OwnerId      string        `bson:"owner_id"`
	Secret       string        `bson:"secret"`
	CommitId     string        `bson:"commit_id"`
	TreeId       string        `json:"-" bson:"tree_id"`
	Contributors []Contributor `bson:"contibutors"`
}

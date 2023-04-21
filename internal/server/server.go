package server

import "fmt"

func CheckAuthentication() (bool, error) {
	return true, nil
}

func CreateNewRepo(name string) (string, error) {
	return fmt.Sprintf("https://localhost:8080/repos/%s", name), nil
}

func CheckAccess(repo string) (bool, error) {
	return true, nil
}

func PushCount(repo string) (int, error) {

	return 0, nil
}


type User struct {
	email string
	name  string
}

func GetUser() (*User, bool, error) {
	user := &User{
		email: "mundanelizard@gmail.com",
		name:  "Mundane Lizard",
	}

	return user, true, nil
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

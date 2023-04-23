package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var endpoint = "http://localhost:5000/"

func CreateAccount(username, password string) error {
	if len(username) == 0 {
		return errors.New("invalid username: expected username to be greater than zero")
	} else if len(password) < 12 {
		return errors.New("invalid password: expected password to have a min length of 12")
	}

	url, err := url.JoinPath(endpoint, "/repos/")
	if err != nil {
		return err
	}

	data := map[string]string {
		"Username": username,
		"Password": password,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		return errors.New(string(body))
	}

	return nil
}

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

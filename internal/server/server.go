package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mundanelizard/envi/internal/models"
)

var endpoint = "http://localhost:5000/"

func validUserDetail(username, password string) error {
	if len(username) == 0 {
		return errors.New("invalid username: expected username to be greater than zero")
	} else if len(password) < 12 {
		return errors.New("invalid password: expected password to have a min length of 12")
	}

	return nil
}

func CreateAccount(username, password string) error {
	err := validUserDetail(username, password)
	if err != nil {
		return err
	}

	url, err := url.JoinPath(endpoint, "/repos/")
	if err != nil {
		return err
	}

	data := map[string]string{
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
	defer res.Body.Close()

	if res.StatusCode != 201 {
		return errors.New(string(body))
	}

	return nil
}

func AuthenticationAccount(username, password string) error {
	err := validUserDetail(username, password)
	if err != nil {
		return err
	}

	url, err := url.JoinPath(endpoint, "/repos/")
	if err != nil {
		return err
	}

	data := map[string]string{
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
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(string(body))
	}

	token := string(body)

	// todo => save token to cache directory
	saveToken(token)

	return nil
}

func RetrieveUser() (*models.User, error) {
	token, err := retrieveToken()
	if err != nil {
		return nil, err
	}

	url, err := url.JoinPath(endpoint, "/users/me")
	if err != nil {
		return nil, err
	}

	// read token from cache directory and return the data back to the user.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Access-Token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	err = json.Unmarshal(body, user)

	return user, nil
}

func RetrieveRepo(username, name string) (*models.Repo, error) {
	token, err := retrieveToken()
	if err != nil {
		return nil, err
	}

	url, err := url.JoinPath(endpoint, "/repos/", username, name)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Access-Token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	repo := &models.Repo{}
	err = json.Unmarshal(body, repo)

	return repo, nil
}

func CreateNewRepo(username, name string) (string, error) {
	token, err := retrieveToken()
	if err != nil {
		return "", err
	}

	url, err := url.JoinPath(endpoint, "/repos/", username, name)
	if err != nil {
		return "", err
	}

	// gen secret

	data := map[string]string {
		"Name": name,
		"Secret": secret,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}


	req, err := http.NewRequest(http.MethodPost, url, body)

	return fmt.Sprintf("https://localhost:8080/repos/%s", name), nil
}

func PushRepo(username, name string) (error) {

	return nil
}

func PullRespoitory(username, name string) (io.ReadCloser, error) {
	return nil, nil
}
package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mundanelizard/envi/internal/crypto"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/mundanelizard/envi/internal/lockfile"
	"github.com/mundanelizard/envi/internal/models"
)

type Server struct {
	endpoint string
}

func New(endpoint string) *Server {
	return &Server{
		endpoint: endpoint,
	}
}

func validUserDetail(username, password string) error {
	if len(username) == 0 {
		return errors.New("invalid username: expected username to be greater than zero")
	} else if len(password) < 12 {
		return errors.New("invalid password: expected password to have a min length of 12")
	}

	return nil
}

func (srv *Server) CreateAccount(username, password string) error {
	err := validUserDetail(username, password)
	if err != nil {
		return err
	}

	url, err := url.JoinPath(srv.endpoint, "/users/")
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

func (srv *Server) AuthenticateAccount(username, password string) error {
	err := validUserDetail(username, password)
	if err != nil {
		return err
	}

	url, err := url.JoinPath(srv.endpoint, "/users/login")
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
	token = strings.ReplaceAll(token, "\"", "")

	err = srv.saveToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) RetrieveUser() (*models.User, error) {
	token, err := srv.retrieveToken()
	if err != nil {
		return nil, err
	}

	url, err := url.JoinPath(srv.endpoint, "/users/me")
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
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (srv *Server) RetrieveRepo(repoPath string) (*models.Repo, error) {
	token, err := srv.retrieveToken()
	if err != nil {
		return nil, err
	}

	url, err := url.JoinPath(srv.endpoint, repoPath)
	if err != nil {
		return nil, err
	}

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

	repo := &models.Repo{}
	err = json.Unmarshal(body, repo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (srv *Server) CreateNewRepo(repo string) (string, error) {
	token, err := srv.retrieveToken()
	if err != nil {
		return "", err
	}

	url, err := url.JoinPath(srv.endpoint, "/repos/")
	if err != nil {
		return "", err
	}

	secret := crypto.GenRandomString()
	hash, err := crypto.GenHash(secret)
	if err != nil {
		return "", err
	}

	data := map[string]string{
		"Name":   repo,
		"Secret": hash,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", nil
	}
	req.Header.Set("Access-Token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 201 {
		return "", errors.New(string(body))
	}

	fmt.Println("Repository secret [Save It]: ", secret)
	fmt.Println("Repository address: ", repo)

	return strings.ReplaceAll(string(body), "\"", ""), nil
}

func (srv *Server) PushRepo(repo, treeId, newTreeId, newCommitId, filepath, secret string) error {
	token, err := srv.retrieveToken()
	if err != nil {
		return err
	}

	url, err := url.JoinPath(srv.endpoint, repo, "push")
	if err != nil {
		return err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("repo", "envi.zip")
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Access-Token", token)
	req.Header.Set("Repo-Tree-Id", treeId)
	req.Header.Set("Next-Tree-Id", newTreeId)
	req.Header.Set("Next-Commit-Id", newCommitId)
	req.Header.Set("Repo-Secret", secret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(string(data))
	}

	fmt.Println("Repository secret [Save It]: ", secret)
	fmt.Println("Repository address: ", repo)

	return nil
}

func (srv *Server) PullRepo(repo string) (string, error) {
	token, err := srv.retrieveToken()
	if err != nil {
		return "", err
	}

	url, err := url.JoinPath(srv.endpoint, repo, "pull")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", nil
	}
	req.Header.Set("Access-Token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	tempDir := os.TempDir()
	dir := path.Join(tempDir, path.Base(repo)+".envi.download")

	err = lockfile.WriteWithLock(dir, body)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func (srv *Server) ShareRepo(repo, user, role string) error {
	token, err := srv.retrieveToken()
	if err != nil {
		return err
	}

	url, err := url.JoinPath(srv.endpoint, repo, "share")
	if err != nil {
		return err
	}

	data := map[string]string{
		"Username": user,
		"Role":     role,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil
	}
	req.Header.Set("Access-Token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return errors.New(string(body))
}

func (srv *Server) RevokeRepo(repo, user string) error {
	token, err := srv.retrieveToken()
	if err != nil {
		return err
	}

	url, err := url.JoinPath(srv.endpoint, repo, "revoke")
	if err != nil {
		return err
	}

	data := map[string]string{
		"Username": user,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil
	}
	req.Header.Set("Access-Token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return errors.New(string(body))
}

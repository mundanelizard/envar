package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/mundanelizard/envi/internal/models"
)

func (srv *server) extractUserFromBody(reader io.Reader) (*models.User, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var user models.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type LeanRepo struct {
	Name string
}

func (srv *server) extractRepoFromBody(reader io.Reader) (*LeanRepo, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var repo LeanRepo

	err = json.Unmarshal(body, &repo)
	if err != nil {
		return nil, err
	}

	if len(repo.Name) == 0 {
		return nil, errors.New("invalid body: expecting name field of type string")
	}

	return &repo, nil
}

var ErrUnauthorised = errors.New("unauthorised request")

func (srv *server) extractUserFromHeaderToken(header http.Header) (*models.User, error) {
	token := header.Get("Access-Token")
	if len(token) == 0 {
		return nil, ErrUnauthorised
	}

	var secret models.Secret
	query := map[string]string{"token": token}
	err := srv.db.Collection("secrets").FindOne(srv.ctx, query).Decode(&secret)
	if err != nil {
		srv.logger.Warn(err.Error())
		return nil, ErrUnauthorised
	}

	var user models.User
	query = map[string]string{"_id": secret.OwnerId}
	err = srv.db.Collection("user").FindOne(srv.ctx, query).Decode(&user)
	if err != nil {
		srv.logger.Warn(err.Error())
		return nil, ErrUnauthorised
	}

	return &user, nil
}

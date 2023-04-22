package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mundanelizard/envi/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func genErrorRes(msg string) map[string]string {
	return map[string]string{
		"message": msg,
	}
}

func (srv *server) handleSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := srv.extractUserFromBody(r.Body)
	if err != nil {
		srv.send(w, http.StatusBadRequest, genErrorRes(err.Error()))
		return
	}

	if err = models.IsValidUser(*user); err != nil {
		srv.send(w, http.StatusBadRequest, genErrorRes(err.Error()))
		return
	}

	user.Password = hashPassword(user.Password)

	_, err = srv.db.Collection("users").InsertOne(srv.ctx, user)
	if err != nil {
		srv.send(w, 500, genErrorRes(err.Error()))
		return
	}

	srv.send(w, 201, user)
}

func (srv *server) handleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := srv.extractUserFromBody(r.Body)
	if err != nil {
		srv.send(w, http.StatusBadRequest, genErrorRes(err.Error()))
		return
	}

	if err = models.IsValidUser(*req); err != nil {
		srv.send(w, http.StatusBadRequest, genErrorRes(err.Error()))
		return
	}

	var user models.User
	err = srv.db.Collection("users").FindOne(srv.ctx, map[string]string{"username": req.Username}).Decode(&user)
	if err != mongo.ErrNoDocuments {
		srv.send(w, http.StatusBadRequest, genErrorRes("user already exists in database"))
		return
	}

	if !verifyPassword(req.Password, user.Password) {
		srv.send(w, http.StatusBadRequest, genErrorRes("Invalid username or password"))
		return
	}

	token := genRandomString()
	secret := models.Secret{
		OwnerId: user.Id,
		Token:   token,
	}

	_, err = srv.db.Collection("secrets").InsertOne(srv.ctx, secret)
	if err != nil {
		srv.send(w, http.StatusBadRequest, genErrorRes(err.Error()))
		return
	}

	srv.send(w, 200, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (srv *server) handleCreateRepo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		srv.send(w, http.StatusUnauthorized, genErrorRes(err.Error()))
		return
	}

	body, err := srv.extractRepoFromBody(r.Body)
	if err != nil {
		srv.send(w, http.StatusBadRequest, genErrorRes(err.Error()))
		return
	}

	repoName := body.Name + ":" + user.Username

	var oldRepo models.Repo
	err = srv.db.Collection("users").FindOne(srv.ctx, map[string]string{"name": repoName}).Decode(oldRepo)
	if err != mongo.ErrNoDocuments {
		srv.send(w, http.StatusBadRequest, genErrorRes("user already exists in database"))
		return
	}

	secret := genRandomString()

	repo := &models.Repo{
		Name:         repoName,
		Secret:       hashPassword(secret),
		Contributors: []models.Contributor{},
	}

	_, err = srv.db.Collection("repos").InsertOne(srv.ctx, repo)
	if err != nil {
		srv.send(w, 500, genErrorRes(err.Error()))
		return
	}

	srv.send(w, 200, map[string]interface{}{
		"repo":   repo,
		"secret": secret,
	})
}

func (srv *server) handleGetRepos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		srv.send(w, http.StatusUnauthorized, genErrorRes(err.Error()))
		return
	}

	query := map[string]string{"owner_id": user.Id}
	cur, err := srv.db.Collection("repos").Find(srv.ctx, query)
	if err != nil {
		srv.send(w, http.StatusInternalServerError, genErrorRes(err.Error()))
		return
	}

	var results []models.Repo

	for cur.Next(srv.ctx) {
		var repo models.Repo
		err := cur.Decode(&repo)
		if err != nil {
			srv.send(w, http.StatusInternalServerError, genErrorRes(err.Error()))
			return
		}
		results = append(results, repo)
	}

	if err := cur.Err(); err != nil {
		srv.send(w, http.StatusInternalServerError, genErrorRes(err.Error()))
		return
	}

	cur.Close(srv.ctx)

	srv.send(w, http.StatusOK, results)
}

func (srv *server) handleGetRepo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		srv.send(w, http.StatusUnauthorized, genErrorRes(err.Error()))
		return
	}

	username := params.ByName("user")
	repoName := params.ByName("repo")

	key := username + ":" + repoName
	query := map[string]string{"owner_id": user.Id, "name": key}

	var repo models.Repo
	err = srv.db.Collection("repos").FindOne(srv.ctx, query).Decode(&repo)
	if err != mongo.ErrNoDocuments {
		srv.send(w, http.StatusBadRequest, genErrorRes(err.Error()))
		return
	}

	srv.send(w, http.StatusOK, repo)
}

func (srv *server) handlePull(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
}

func (srv *server) handlePush(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (srv *server) handleUpdateRepo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (srv *server) handleShareRepo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (srv *server) handleRemoveAccess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

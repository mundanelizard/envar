package main

import (
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mundanelizard/envi/internal/models"
)

func genErrorRes(msg string) map[string]string {
	return map[string]string{
		"message": msg,
	}
}

func (srv *server) handleSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		srv.send(w, 400, genErrorRes("Failed to read body"))
		return
	}

	user, err := srv.extractUserFromBody(body)
	if err != nil {
		srv.send(w, 400, genErrorRes(err.Error()))
		return
	}

	ok, err := srv.validUser(user)
	if err != nil {
		srv.send(w, 400, genErrorRes(err.Error()))
		return
	}

	user.Password = srv.hashPassword(user.Password)

	_, err = srv.db.Collection("users").InsertOne(srv.ctx, user)
	if err != nil {
		srv.send(w, 500, genErrorRes(err.Error()))
		return
	}

	srv.send(w, 201, user)
}

func (srv *server) handleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		srv.send(w, 400, genErrorRes("Failed to read body"))
		return
	}

	req, err := srv.extractUserFromBody(body)
	if err != nil {
		srv.send(w, 400, genErrorRes(err.Error()))
		return
	}

	var user models.User
	err = srv.db.Collection("users").FindOne(srv.ctx, map[string]string{"email": req.Email}).Decode(&user)
	if err != nil {
		srv.send(w, 400, genErrorRes(err.Error()))
		return
	}

	if !srv.verifyPassword(req.Password, user.Password) {
		srv.send(w, 400, genErrorRes("Invalid username or password"))
		return
	}

	token := genRandomString()
	code := srv.hashPassword(user.Password)

	secret := &models.Secret{
		OwnerId: user.Id,
		Secret:  code,
		Token:   token,
	}

	_, err = srv.db.Collection("secrets").InsertOne(srv.ctx, user)
	if err != nil {
		srv.send(w, 400, genErrorRes(err.Error()))
		return
	}

	srv.send(w, 200, map[string]interface{}{
		"token":  token,
		"secret": secret,
		"user":   user,
	})
}

func handleCreateRepo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func handleGetRepos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func handleUploadRepoHandshakeInit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func handleVerifyRepoHandshakeVerify(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func handleUpdateRepo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func handleShareRepo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func handleRemoveAccess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

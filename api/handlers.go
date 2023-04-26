package main

import (
	"fmt"
	"github.com/mundanelizard/envi/internal/crypto"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/mundanelizard/envi/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (srv *server) handleSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := srv.extractUserFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = models.IsValidUser(*user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = srv.db.Collection("users").FindOne(srv.ctx, map[string]string{"username": user.Username}).Decode(&user)
	if err != mongo.ErrNoDocuments {
		http.Error(w, "user already exists in database", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"username": user.Username,
		"password": crypto.HashPassword(user.Password),
	}

	_, err = srv.db.Collection("users").InsertOne(srv.ctx, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.send(w, http.StatusCreated, "user creation successful")
}

func (srv *server) handleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := srv.extractUserFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = models.IsValidUser(*req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	err = srv.db.Collection("users").FindOne(srv.ctx, map[string]string{"username": req.Username}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !crypto.VerifyPassword(req.Password, user.Password) {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}

	token := crypto.GenRandomString()
	secret := models.Secret{
		OwnerId: user.Id,
		Token:   token,
	}

	_, err = srv.db.Collection("secrets").InsertOne(srv.ctx, secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.send(w, http.StatusOK, token)
}

func (srv *server) handleGetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.send(w, http.StatusOK, user)
}

func (srv *server) handleCreateRepo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := srv.extractRepoFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	repoName := body.Name + "-" + user.Username

	var oldRepo models.Repo
	err = srv.db.Collection("repo").FindOne(srv.ctx, map[string]string{"name": repoName}).Decode(oldRepo)
	if err != mongo.ErrNoDocuments {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	repo := map[string]interface{}{
		"name":         repoName,
		"secret":       body.Secret,
		"owner_id":     user.Id,
		"contributors": []models.Contributor{},
	}

	_, err = srv.db.Collection("repos").InsertOne(srv.ctx, repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.send(w, http.StatusCreated, fmt.Sprintf("/repos/%s/%s", user.Username, body.Name))
}

func (srv *server) handleGetRepos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := map[string]string{"owner_id": user.Id}
	cur, err := srv.db.Collection("repos").Find(srv.ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var results []models.Repo

	for cur.Next(srv.ctx) {
		var repo models.Repo
		err := cur.Decode(&repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		results = append(results, repo)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cur.Close(srv.ctx)

	srv.send(w, http.StatusOK, results)
}

func (srv *server) handleGetRepo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := params.ByName("user")
	repoName := params.ByName("repo")

	key := repoName + "-" + username
	query := bson.M{
		"name": key,
		"$or": []interface{}{
			bson.M{"owner_id": user.Id},
			bson.M{"contributors.user_id": user.Id, "contributors.role": "W"},
		},
	}

	var repo models.Repo
	err = srv.db.Collection("repos").FindOne(srv.ctx, query).Decode(&repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.send(w, http.StatusOK, repo)
}

func (srv *server) handlePull(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := params.ByName("user")
	repoName := params.ByName("repo")

	key := repoName + "-" + username

	query := bson.M{
		"name": key,
		"$or": []interface{}{
			bson.M{"owner_id": user.Id},
			bson.M{"contributors.user_id": user.Id, "contributors.role": "W"},
		},
	}

	var repo models.Repo
	err = srv.db.Collection("repos").FindOne(srv.ctx, query).Decode(&repo)
	if err != mongo.ErrNoDocuments {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	repoPath := path.Join(srv.dir.uploads, key)

	srv.sendFile(w, repoPath)
}

func (srv *server) handlePush(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := params.ByName("user")
	repoName := params.ByName("repo")
	key := repoName + "-" + username

	query := bson.M{
		"name": key,
		"$or": []interface{}{
			bson.M{"owner_id": user.Id},
			bson.M{"contributors.user_id": user.Id, "contributors.role": "W"},
		},
	}

	commitId := r.Header.Get("Next-Commit-Id")
	treeId := r.Header.Get("Next-Tree-Id")

	if len(commitId) == 0 || len(treeId) == 0 {
		http.Error(w, "invalid commit-id or next-tree-id", http.StatusBadRequest)
		return
	}

	var repo models.Repo
	err = srv.db.Collection("repos").FindOne(srv.ctx, query).Decode(&repo)
	if err == mongo.ErrNoDocuments {
		http.Error(w, "repository not found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(repo.TreeId) != 0 && repo.TreeId != r.Header.Get("Repo-Tree-Id") {
		http.Error(w, "invalid header: repo tree id is invalid and doesn't match head", http.StatusBadRequest)
		return
	}

	secret := r.Header.Get("Repo-Secret")
	if err = crypto.VerifyHash(secret, repo.Secret); err != nil {
		http.Error(w, "invalid header: secret is invalid and doesn't match db content", http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("repo")
	if err != nil {
		http.Error(w, "error retrieving repository from form", http.StatusBadRequest)
		return
	}

	repoPath := path.Join(srv.dir.uploads, key)

	local, err := os.OpenFile(repoPath, os.O_CREATE|os.O_RDWR, 0655)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := io.Copy(local, file); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query = bson.M{
		"name": key,
	}

	update := bson.M{
		"$set": bson.M{
			"commit_id": commitId,
			"tree_id":   treeId,
		},
	}

	_, err = srv.db.Collection("repos").UpdateOne(srv.ctx, query, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	srv.send(w, http.StatusOK, "")
}

func (srv *server) handleShareRepo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := params.ByName("user")
	repoName := params.ByName("repo")
	key := username + "-" + repoName

	share, err := srv.extractShareRepoFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := map[string]string{"owner_id": user.Id, "name": key}

	var repo models.Repo
	err = srv.db.Collection("repos").FindOne(srv.ctx, query).Decode(&repo)
	if err != nil {
		http.Error(w, "invalid header: can not find repository", http.StatusBadRequest)
		return
	}

	secret := r.Header.Get("Repo-Secret")
	if !crypto.VerifyPassword(secret, repo.Secret) {
		http.Error(w, "invalid header: secret is invalid and doesn't match db content", http.StatusBadRequest)
		return
	}

	contributor := models.Contributor{
		UserId: share.Id,
		Role:   share.Role,
	}

	filter := bson.M{"name": key, "owner_id": user.Id}
	update := bson.M{"$push": bson.M{"contributors": contributor}}

	_, err = srv.db.Collection("repos").UpdateOne(srv.ctx, filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.send(w, http.StatusOK, map[string]bool{"success": true})
}

func (srv *server) handleRemoveAccess(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	share, err := srv.extractShareRepoFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contributor := models.Contributor{
		UserId: share.Id,
		Role:   share.Role,
	}

	username := params.ByName("user")
	repoName := params.ByName("repo")
	key := username + "-" + repoName

	query := map[string]string{"owner_id": user.Id, "name": key}

	var repo models.Repo
	err = srv.db.Collection("repos").FindOne(srv.ctx, query).Decode(&repo)
	if err != nil {
		http.Error(w, "invalid header: can not find repository", http.StatusBadRequest)
		return
	}

	secret := r.Header.Get("Repo-Secret")
	if !crypto.VerifyPassword(secret, repo.Secret) {
		http.Error(w, "invalid header: secret is invalid and doesn't match db content", http.StatusBadRequest)
		return
	}

	filter := bson.M{"name": key, "owner_id": user.Id}
	update := bson.M{"$pull": bson.M{"contributors": contributor}}

	_, err = srv.db.Collection("repos").UpdateOne(srv.ctx, filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.send(w, http.StatusOK, map[string]bool{"success": true})
}

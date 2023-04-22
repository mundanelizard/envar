package main

import (
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

	user.Password = hashPassword(user.Password)

	_, err = srv.db.Collection("users").InsertOne(srv.ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.send(w, 201, user)
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
	if err != mongo.ErrNoDocuments {
		http.Error(w, "user already exists in database", http.StatusBadRequest)
		return
	}

	if !verifyPassword(req.Password, user.Password) {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}

	token := genRandomString()
	secret := models.Secret{
		OwnerId: user.Id,
		Token:   token,
	}

	_, err = srv.db.Collection("secrets").InsertOne(srv.ctx, secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	repo := &models.Repo{
		Name:         repoName,
		Secret:       body.Secret,
		Contributors: []models.Contributor{},
	}

	_, err = srv.db.Collection("repos").InsertOne(srv.ctx, repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.send(w, 200, map[string]interface{}{
		"repo": repo,
	})
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

	key := username + "-" + repoName
	contributorQuery := map[string]string{"contributors": user.Id}
	ownerQuery := map[string]string{"owner_id": user.Id}
	subQueries := []map[string]string{contributorQuery, ownerQuery}
	query := map[string]interface{}{"$or": subQueries, "name": key}

	var repo models.Repo
	err = srv.db.Collection("repos").FindOne(srv.ctx, query).Decode(&repo)
	if err != mongo.ErrNoDocuments {
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

	key := username + "-" + repoName

	contributorQuery := map[string]string{"contributors": user.Id}
	ownerQuery := map[string]string{"owner_id": user.Id}
	subQueries := []map[string]string{contributorQuery, ownerQuery}
	query := map[string]interface{}{"$or": subQueries, "name": key}

	var repo models.Repo
	err = srv.db.Collection("repos").FindOne(srv.ctx, query).Decode(&repo)
	if err != mongo.ErrNoDocuments {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := path.Join(dir, ".envi-server", "uploads", username)

	srv.sendFile(w, path)
}

func (srv *server) handlePush(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := srv.extractUserFromHeaderToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := params.ByName("user")
	repoName := params.ByName("repo")
	key := username + "-" + repoName

	contributorQuery := map[string]string{"contributors.owner_id": user.Id, "contributors.role": "W"}
	ownerQuery := map[string]string{"owner_id": user.Id}
	subQueries := []map[string]string{contributorQuery, ownerQuery}
	query := map[string]interface{}{"$or": subQueries, "name": key}

	var repo models.Repo
	err = srv.db.Collection("repos").FindOne(srv.ctx, query).Decode(&repo)
	if err != mongo.ErrNoDocuments {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secret := r.Header.Get("Repo-Secret")
	if !verifyPassword(secret, repo.Secret) {
		http.Error(w, "invalid header: secret is invalid and doesn't match db content", http.StatusBadRequest)
		return
	}

	content := r.MultipartForm.File["repo"][0]
	if content == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := content.Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	local, err := os.OpenFile(repoName, os.O_CREATE|os.O_RDWR, 0655)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := io.Copy(local, file); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.send(w, http.StatusOK, map[string]bool{"success": true})
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
	if !verifyPassword(secret, repo.Secret) {
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
	if !verifyPassword(secret, repo.Secret) {
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

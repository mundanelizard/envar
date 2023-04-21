package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (srv *server) routes() http.Handler {
	router := httprouter.New()

	router.POST("/users/", handleSignup)
	router.POST("/users/login", handleLogin)

	router.POST("/repos/:user", handleCreateRepo)
	router.POST("/repos/:user", handleGetRepos)

	router.POST("/repos/:user/:repo/handshakes/init", handleUploadRepoHandshakeInit)
	router.POST("/repos/:user/:repo/handshakes/verify", handleVerifyRepoHandshakeVerify)

	router.POST("/repos/:user/:repo/upload", handleUpdateRepo)
	router.POST("/repos/:user/:repo/share", handleShareRepo)
	router.POST("repos/:user/:repo/revoke", handleRemoveAccess)

	return router
}

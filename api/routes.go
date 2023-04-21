package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (srv *server) routes() http.Handler {
	router := httprouter.New()

	router.POST("/users/", srv.handleSignup)
	router.POST("/users/login", srv.handleLogin)

	router.POST("/repos/:user", handleCreateRepo)
	router.POST("/repos/:user", handleGetRepos)

	router.POST("/repos/:user/:repo/recent", handleUploadRepoHandshakeInit)
	router.POST("/repos/:user/:repo/verify", handleVerifyRepoHandshakeVerify)

	router.GET("/repos/:user/:repo/pull", srv.handlePullRequest)

	router.POST("/repos/:user/:repo/upload", handleUpdateRepo)
	router.POST("/repos/:user/:repo/share", handleShareRepo)
	router.POST("repos/:user/:repo/revoke", handleRemoveAccess)

	return router
}

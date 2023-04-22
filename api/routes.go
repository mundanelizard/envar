package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (srv *server) routes() http.Handler {
	router := httprouter.New()

	router.POST("/users/", srv.handleSignup)
	router.POST("/users/login", srv.handleLogin)

	router.POST("/repos/", srv.handleCreateRepo)

	router.GET("/repos/:user/", srv.handleGetRepos)
	router.GET("/repos/:user/:repo/", handleGetRepo)

	router.GET("/repos/:user/:repo/pull", srv.handlePull)
	router.POST("/repos/:user/:repo/push", srv.handlePush)

	router.POST("/repos/:user/:repo/upload", handleUpdateRepo)
	router.POST("/repos/:user/:repo/share", handleShareRepo)
	router.POST("repos/:user/:repo/revoke", handleRemoveAccess)

	return router
}

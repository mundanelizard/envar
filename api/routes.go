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
	router.GET("/repos/:user/:repo/", srv.handleGetRepo)

	router.GET("/repos/:user/:repo/pull", srv.handlePull)
	router.POST("/repos/:user/:repo/push", srv.handlePush)

	router.POST("/repos/:user/:repo/share", srv.handleShareRepo)
	router.POST("repos/:user/:repo/revoke", srv.handleRemoveAccess)

	return router
}

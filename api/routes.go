package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (srv *server) routes() http.Handler {
	router := httprouter.New()

	return router
}

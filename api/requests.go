package main

import (
	"encoding/json"
	"net/http"
)

func (srv *server) send(w http.ResponseWriter, status int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	
	w.WriteHeader(status)
	w.Write(bytes)
}
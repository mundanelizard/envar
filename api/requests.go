package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (srv *server) send(w http.ResponseWriter, status int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(status)
	w.Write(bytes)
}

func (srv *server) sendFile(w http.ResponseWriter, path string) {

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer file.Close()

	contentType := "application/octet-stream"
	w.Header().Set("Content-Type", contentType)

	filename := filepath.Base(file.Name())
	disposition := fmt.Sprintf("attachment; filename=%s", filename)

	w.Header().Set("Content-Disposition", disposition)

	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

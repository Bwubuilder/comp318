package database

import (
	"log/slog"
	"net/http"
)

type MainHandler struct{}

func New() MainHandler {
	slog.Info("Created MainHandler")
	return MainHandler{}
}

func (m MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Serve Activated")
	db := newDatabase()
	http.HandleFunc("/", db.dbMethods)
}

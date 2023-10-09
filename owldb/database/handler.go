package database

import (
	"log/slog"
	"net/http"
)

type MainHandler struct {
}

func New() http.Handler {
	mux := http.NewServeMux()
	auth := newAuth()
	db := newDatabase()
	mux.HandleFunc("/auth", auth.authFunction)
	mux.HandleFunc("/", db.dbMethods)

	slog.Info("Created MainHandler")
	return mux
}

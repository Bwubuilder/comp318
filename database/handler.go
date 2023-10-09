package database

import (
	"log/slog"
	"net/http"
)

func setMux() *http.ServeMux {
	mux := http.NewServeMux()
	auth := newAuth()
	db := newDatabase()
	mux.HandleFunc("/auth", auth.authFunction)
	mux.HandleFunc("/", db.dbMethods)

	slog.Info("Created MainHandler")
	return mux
}

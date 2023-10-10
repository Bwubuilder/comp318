package database

import (
	"log/slog"
	"net/http"
)

func SetMux() *http.ServeMux {
	mux := http.NewServeMux()
	db := newDatabase()
	mux.Handle("/auth", newAuth())
	mux.HandleFunc("/", db.dbMethods)

	slog.Info("Created MainHandler")
	return mux
}

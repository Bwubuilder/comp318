package database

import (
	"log/slog"
	"net/http"

	"github.com/Bwubuilder/owldb/authorization"
)

func SetMux() *http.ServeMux {
	mux := http.NewServeMux()
	db := newDatabase()
	mux.Handle("/auth", authorization.NewAuth())
	mux.HandleFunc("/", db.dbMethods)

	slog.Info("Created MainHandler")
	return mux
}

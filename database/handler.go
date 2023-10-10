package database

import (
	"log/slog"
	"net/http"

	"github.com/Bwubuilder/owldb/authorization"
)

func SetMux() *http.ServeMux {
	mux := http.NewServeMux()
	db := newDatabase()
	auth := authorization.NewAuth()
	mux.HandleFunc("/auth", auth.HandleAuthFunctions)
	mux.HandleFunc("/", db.dbMethods)

	slog.Info("Created MainHandler")
	return mux
}

package database

import (
	"log/slog"
	"net/http"
)

type databaseHandler struct {
	db DatabaseService
}

func newDatabase() databaseHandler {
	return databaseHandler{db: *NewDatabaseService()}
}

func (d *databaseHandler) ServeHTTP() {
	slog.Info("Database Created")
}

func (dh *databaseHandler) dbMethods(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		dh.db.HandleGet(w, r)
	case http.MethodPut:
		dh.db.HandlePut(w, r)
	case http.MethodPost:
		dh.db.HandlePost(w, r)
	case http.MethodPatch:
		dh.db.HandlePatch(w, r)
	case http.MethodDelete:
		dh.db.HandleDelete(w, r)
	case http.MethodOptions:
		dh.db.HandleOptions(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

package handlers

import (
	"net/http"
)

// PingDB is a handler that returns a 200 OK response with a message.
func (us *URLShortener) PingDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	if us.DBStorage == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	if err := us.DBStorage.Ping(); err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database connection is alive"))
}

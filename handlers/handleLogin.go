package handlers

import (
	"net/http"
)

// HandleLogin resolve login request
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login"))
}
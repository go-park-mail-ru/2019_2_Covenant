package handlers

import "net/http"

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout"))
}